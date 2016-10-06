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

	"internal/core/link"
	"internal/database/minio"
	"internal/strings/strutil"

	"github.com/spkg/bom"
)

func proc(data []byte) {
	s := time.Now()
	b, o, err := minio.Unmarshal(data)
	if err != nil {
		log.Println("proc: err: pair:", err)
		return
	}

	f1, err := minio.Get(b, o)
	if err != nil {
		log.Println("proc: err: load:", o, err)
		return
	}
	defer minio.Free(f1)

	var f string
	m, d, err := procObject(f1)
	if err != nil {
		log.Println("proc: err:", o, err)
		m.Fail = err.Error()

		f = o + ".txt"
		err = minio.Put(BucketStreamErr, f, bytes.NewReader(m.MarshalIndent()))
		if err != nil {
			log.Println("proc: err: save:", f, err)
		}

		err = minio.Copy(BucketStreamErr, o, b, o)
		if err != nil {
			log.Println("proc: err: copy:", o, err)
		}
	} else {
		f = MakeFileName(m.Auth.ID, m.UUID, m.HTag)
		err = minio.Put(BucketStreamOut, f, d)
		if err != nil {
			log.Println("proc: err: save:", f, err)
		}

		log.Println("proc:", f, m.Proc, time.Since(s).String())
	}

	err = minio.Del(b, o)
	if err != nil {
		log.Println("proc: err: kill:", o, err)
	}

	err = SetZlog(m)
	if err != nil {
		log.Println("proc: err: zlog:", o, err)
	}
}

func procObject(r io.Reader) (Meta, io.Reader, error) {
	m := Meta{}

	meta, data, err := UnpackMetaData(r)
	if err != nil {
		return m, nil, err
	}

	m, err = UnmarshalMeta(meta)
	if err != nil {
		return m, nil, err
	}

	v, err := unmarshalData(data, &m)
	if err != nil {
		return m, nil, err
	}

	d, err := mineLinks(v, &m)
	if err != nil {
		return m, nil, err
	}

	t, err := PackMetaData(m.Marshal(), d)
	if err != nil {
		return m, nil, err
	}

	return m, t, nil
}

func killUTF8BOM(data []byte) []byte {
	if strings.Contains(http.DetectContentType(data), "text/plain; charset=utf-8") {
		return bom.Clean(data)
	}
	return data
}

const magicConvString = "conv"

func unmarshalData(data []byte, m *Meta) (interface{}, error) {
	d := killUTF8BOM(data)
	m.ETag = btsToMD5(d)
	m.Size = int64(len(d))

	if strings.HasPrefix(m.CTag, magicConvString) {
		m.CTag = fmt.Sprintf("converted from %s format", m.HTag)
		m.HTag = ConvHTag(m.HTag)
		return unmarshalDataOLD(d, m)
	}

	return unmarshalDataNEW(d, m)
}

func unmarshalDataOLD(data []byte, m *Meta) (interface{}, error) {
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

func unmarshalDataNEW(data []byte, m *Meta) (interface{}, error) {
	t := m.HTag

	switch {
	case isGeo(t):
		v := link.DataV3Geoa{}
		err := json.Unmarshal(data, &v)
		return v, err
	case isSaleBY(t):
		v := link.DataV3SaleBy{}
		err := json.Unmarshal(data, &v)
		return v, err
	default:
		v := link.DataV3Sale{}
		err := json.Unmarshal(data, &v)
		return v, err
	}
}

func mineLinks(v interface{}, m *Meta) ([]byte, error) {
	t := m.HTag
	s := time.Now()

	a, err := GetLinkAuth([]string{m.Auth.ID})
	if err != nil {
		return nil, err
	}
	m.Auth = a[0]

	l, err := GetLinkAddr([]string{strToSHA1(makeMagicHead(m.Name, m.Head, m.Addr))})
	if err != nil {
		return nil, err
	}
	m.Link = l[0]

	n := 0
	if r, ok := v.(link.Ruler); ok {
		n = r.Len()
		m.Proc = fmt.Sprintf("%d", n)
		if n == 0 {
			return nil, fmt.Errorf("no data")
		}
	}

	if d, ok := v.(link.Druger); ok {
		n, err = mineDrugs(d, t)
	}
	if err != nil {
		return nil, err
	}
	m.Proc = fmt.Sprintf("%s:%d", m.Proc, n)

	if isSaleIn(t) {
		if a, ok := v.(link.Addrer); ok {
			n, err = mineAddrs(a)
		}
		if err != nil {
			return nil, err
		}
		m.Proc = fmt.Sprintf("%s:%d", m.Proc, n)
	}
	m.Proc = fmt.Sprintf("%s:%s", m.Proc, time.Since(s).String())

	return json.Marshal(v)
}

func mineDrugs(v link.Druger, t string) (int, error) {
	var (
		ext  = filepath.Ext(t)
		keys = make([]string, v.Len())
		name string
	)
	for i := 0; i < v.Len(); i++ {
		name = v.GetName(i)
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

	lds, err := GetLinkDrug(keys)
	if err != nil {
		return 0, err
	}

	n := 0
	for i := 0; i < v.Len(); i++ {
		if v.SetDrug(i, lds[i]) {
			n++
		}
	}

	return n, nil
}

func mineAddrs(v link.Addrer) (int, error) {
	var keys = make([]string, v.Len())
	for i := 0; i < v.Len(); i++ {
		keys[i] = strToSHA1(makeMagicAddr(v.GetSupp(i)))
	}

	lds, err := GetLinkAddr(keys)
	if err != nil {
		return 0, err
	}

	n := 0
	for i := 0; i < v.Len(); i++ {
		if v.SetAddr(i, lds[i]) {
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

func makeMagicHead(name, head, addr string) string {
	return strings.TrimSpace(
		strutil.TrimRightN(
			fmt.Sprintf("%s/%s: %s", name, head, addr),
			magicLength,
		),
	)
}

func makeMagicAddr(name string) string {
	return strings.TrimSpace(
		strutil.TrimRightN(
			name,
			magicLength,
		),
	)
}

func makeMagicDrug(name string) string {
	return strings.TrimSpace(
		strutil.TrimRightN(
			name,
			magicLength,
		),
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
