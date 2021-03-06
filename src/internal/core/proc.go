package core

import (
	"bytes"
	"crypto/md5"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"internal/database/minio"
	"internal/strings/strutil"

	"github.com/spkg/bom"
)

func proc(data []byte) {
	t := time.Now()

	p, err := decodePath(data)
	if err != nil {
		log.Println(err)
		return
	}

	f, err := minio.Get(p.Bucket, p.Object)
	if err != nil {
		log.Println(err)
		return
	}
	defer minio.Free(f)
	defer func() {
		err = minio.Del(p.Bucket, p.Object)
		if err != nil {
			log.Println(err)
		}
	}()

	m, d, err := procObject(f)
	if err != nil {
		// Fail
		err = copyToErrs(p, bytes.NewReader(m.marshalIndent()))
		if err != nil {
			log.Println(err)
		}
	} else {
		// Success
		err = copyToOuts(p, m, d)
		if err != nil {
			log.Println(err)
		}
	}

	m.Proc = fmt.Sprintf("%s %s", m.Proc, time.Since(t).String())
	err = setZlog(m)
	if err != nil {
		log.Println(err)
	}

	if m.Fail == "" {
		log.Println("-->", p.Object, m.Proc)
	} else {
		log.Println("-x-", p.Object, "err:", m.Fail)
	}

}

func copyToErrs(p path, r io.Reader) error {
	err := minio.Copy(bucketStreamErr, p.Object, p.Bucket, p.Object)
	if err != nil {
		return err
	}
	return minio.Put(bucketStreamErr, p.Object+".txt", r)
}

func copyToOuts(p path, m *meta, r io.Reader) error {
	err := minio.Put(bucketStreamOut, p.Object, r)
	if err != nil {
		return err
	}
	if isGeo(m.HTag) {
		p.Bucket = bucketStreamOut
		err = minio.Copy(bucketStreamOutGeo, p.Object, p.Bucket, p.Object)
		if err != nil {
			return err
		}
		//if strings.HasSuffix(m.HTag, ".ua") {
		//	p.Bucket = bucketStreamOutGeo
		//	return minio.Copy(bucketStreamOutGeoTest, p.Object, p.Bucket, p.Object)
		//}
		return nil
	}
	if m.Frwd != "" {
		p.Bucket = bucketStreamOut
		return minio.Copy(bucketStreamOutFrwd, p.Object, p.Bucket, p.Object)
	}
	return nil
}

func procObject(r io.Reader) (*meta, io.Reader, error) {
	m := &meta{}

	fail := func(m *meta, err error) (*meta, io.Reader, error) {
		m.Fail = err.Error()
		return m, nil, err
	}

	meta, data, err := unpackMetaData(r)
	if err != nil {
		return fail(m, err)
	}

	m, err = unmarshalMeta(meta)
	if err != nil {
		return fail(m, err)
	}

	v, err := unmarshalData(data, m)
	if err != nil {
		return fail(m, err)
	}

	l, err := mineLinks(v, m)
	if err != nil {
		return fail(m, err)
	}

	d, err := json.Marshal(l)
	if err != nil {
		return fail(m, err)
	}

	p, err := packMetaData(m.marshal(), d)
	if err != nil {
		return fail(m, err)
	}

	return m, p, nil
}

func killUTF8BOM(data []byte) []byte {
	if strings.Contains(http.DetectContentType(data), "text/plain; charset=utf-8") {
		return bom.Clean(data)
	}
	return data
}

const magicConvString = "conv"

func unmarshalData(data []byte, m *meta) (interface{}, error) {
	d := killUTF8BOM(data)
	m.ETag = btsToMD5(d)
	m.Size = int64(len(d))

	if strings.HasPrefix(m.CTag, magicConvString) {
		m.CTag = fmt.Sprintf("converted from %s format", m.HTag)
		m.HTag = normHTag(m.HTag)
		return unmarshalDataOLD(d, m)
	}

	return unmarshalDataNEW(d, m)
}

func unmarshalDataOLD(data []byte, m *meta) (interface{}, error) {
	t := m.HTag

	switch {
	case isGeo(t):
		return convGeoa(data, m)
	case isSaleBY(t):
		return convSaleBy(data, m)
	default:
		return convSale(data, m)
	}
}

func unmarshalDataNEW(data []byte, m *meta) (interface{}, error) {
	t := m.HTag

	switch {
	case isGeo(t):
		v := jsonV3Geoa{}
		err := json.Unmarshal(data, &v)
		return v, err
	case isSaleBY(t):
		v := jsonV3SaleBy{}
		err := json.Unmarshal(data, &v)
		return v, err
	default:
		v := jsonV3Sale{}
		err := json.Unmarshal(data, &v)
		return v, err
	}
}

func mineLinks(v interface{}, m *meta) (interface{}, error) {
	t := m.HTag
	s := time.Now()

	a, err := getLinkAuth([]string{m.Auth.ID})
	if err != nil {
		return nil, err
	}
	m.Auth = a[0]

	l, err := getLinkAddr([]string{strToSHA1(makeMagicHead(m.Name, m.Head, m.Addr))})
	if err != nil {
		return nil, err
	}
	m.Link = l[0]

	n := 0
	if r, ok := v.(ruler); ok {
		n = r.len()
		m.Proc = fmt.Sprintf("%d", n)
		if n == 0 {
			return nil, fmt.Errorf("no data")
		}
	}

	if d, ok := v.(druger); ok {
		n, err = mineDrugs(d, t)
	}
	if err != nil {
		return nil, err
	}
	m.Proc = fmt.Sprintf("%s:%d", m.Proc, n)

	if isSaleIn(t) || isRcgnAddr(t) {
		if a, ok := v.(addrer); ok {
			n, err = mineAddrs(a)
		}
		if err != nil {
			return nil, err
		}
		m.Proc = fmt.Sprintf("%s:%d", m.Proc, n)
	}
	m.Proc = fmt.Sprintf("%s:%s", m.Proc, time.Since(s).String())

	return v, nil
}

func mineDrugs(v druger, t string) (int, error) {
	var (
		ext  = filepath.Ext(t)
		keys = make([]string, v.len())
		name string
	)
	for i := 0; i < v.len(); i++ {
		name = v.getName(i)
		switch {
		case isUA(ext):
			name = makeMagicDrugUA(name)
		case isRU(ext):
			name = makeMagicDrugRU(name)
		case isKZ(ext):
			name = makeMagicDrugKZ(name)
		case isBY(ext):
			name = makeMagicDrugBY(name)
		default:
			name = makeMagicDrug(name)
		}
		keys[i] = strToSHA1(name)
	}

	lds, err := getLinkDrug(keys)
	if err != nil {
		return 0, err
	}

	n := 0
	for i := 0; i < v.len(); i++ {
		if v.setDrug(i, lds[i]) {
			n++
		}
	}

	return n, nil
}

func mineAddrs(v addrer) (int, error) {
	var keys = make([]string, v.len())
	for i := 0; i < v.len(); i++ {
		keys[i] = strToSHA1(makeMagicAddr(v.getSupp(i)))
	}

	lds, err := getLinkAddr(keys)
	if err != nil {
		return 0, err
	}

	n := 0
	for i := 0; i < v.len(); i++ {
		if v.setAddr(i, lds[i]) {
			n++
		}
	}

	return n, nil
}

const (
	magicLength   = 1024
	magicSuffixBY = "{\"COUNTRY_ID\":\"1010\"}"
	magicSuffixKZ = "{\"COUNTRY_ID\":\"106\"}"
	magicSuffixRU = "{\"COUNTRY_ID\":\"1027\"}"
	magicSuffixUA = ""
)

func makeMagicName(name, addr string) string {
	return strings.Trim(
		strutil.TrimRightN(
			fmt.Sprintf("%s %s", name, addr),
			magicLength,
		),
		" ",
	)
}

func makeMagicHead(name, head, addr string) string {
	return strings.Trim(
		strutil.TrimRightN(
			fmt.Sprintf("%s/%s: %s", name, head, addr),
			magicLength,
		),
		" ",
	)
}

func makeMagicAddr(name string) string {
	return strings.Trim(
		strutil.TrimRightN(
			name,
			magicLength,
		),
		" ",
	)
}

func makeMagicDrug(name string) string {
	return strings.Trim(
		strutil.TrimRightN(
			name,
			magicLength,
		),
		" ",
	)
}

func makeMagicDrugBY(name string) string {
	return makeMagicDrug(name) + magicSuffixBY
}

func makeMagicDrugKZ(name string) string {
	return makeMagicDrug(name) + magicSuffixKZ
}

func makeMagicDrugRU(name string) string {
	return makeMagicDrug(name) + magicSuffixRU
}

func makeMagicDrugUA(name string) string {
	return makeMagicDrug(name) + magicSuffixUA
}

func isGeo(s string) bool {
	return strings.Contains(s, "geo")
}

func isSaleBY(s string) bool {
	return strings.Contains(s, "daily.by")
}

func isSaleIn(s string) bool {
	return strings.Contains(s, "sale-in")
}

func isRcgnAddr(s string) bool {
	return strings.Contains(s, "rcgn.addr")
}

const (
	extBY = ".by"
	extKZ = ".kz"
	extRU = ".ru"
	extUA = ".ua"
)

func isBY(s string) bool {
	return strings.EqualFold(s, extBY)
}

func isKZ(s string) bool {
	return strings.EqualFold(s, extKZ)
}

func isRU(s string) bool {
	return strings.EqualFold(s, extRU)
}

func isUA(s string) bool {
	return strings.EqualFold(s, extUA)
}

func btsToMD5(b []byte) string {
	return fmt.Sprintf("%x", md5.Sum(b))
}

func btsToSHA1(b []byte) string {
	return fmt.Sprintf("%x", sha1.Sum(b))
}

func strToSHA1(s string) string {
	return btsToSHA1([]byte(s))
}
