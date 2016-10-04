package api

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

	"internal/conns/minio"
	"internal/strutil"

	"github.com/spkg/bom"
)

var listHTag = map[string]struct{}{
	"geoapt.ru":           {},
	"geoapt.ua":           {},
	"sale-in.monthly.by":  {},
	"sale-in.monthly.kz":  {},
	"sale-in.monthly.ua":  {},
	"sale-in.weekly.ua":   {},
	"sale-in.daily.kz":    {},
	"sale-in.daily.ua":    {},
	"sale-out.monthly.by": {},
	"sale-out.monthly.kz": {},
	"sale-out.monthly.ru": {},
	"sale-out.monthly.ua": {},
	"sale-out.weekly.ru":  {},
	"sale-out.weekly.ua":  {},
	"sale-out.daily.by":   {},
	"sale-out.daily.kz":   {},
	"sale-out.daily.ua":   {},
}

var convHTag = map[string]string{
	// version 1 -> version 3
	"data.geostore":         "geoapt.ua",
	"data.sale-inp.monthly": "sale-in.monthly.ua",
	"data.sale-inp.weekly":  "sale-in.weekly.ua",
	"data.sale-inp.daily":   "sale-in.daily.ua",
	"data.sale-out.monthly": "sale-out.monthly.ua",
	"data.sale-out.weekly":  "sale-out.weekly.ua",
	"data.sale-out.daily":   "sale-out.daily.ua",
	// version 2 -> version 3
	"data.geoapt.ru":           "geoapt.ru",
	"data.geoapt.ua":           "geoapt.ua",
	"data.sale-inp.monthly.kz": "sale-out.monthly.kz",
	"data.sale-inp.monthly.ua": "sale-in.monthly.ua",
	"data.sale-inp.weekly.ua":  "sale-in.weekly.ua",
	"data.sale-inp.daily.kz":   "sale-in.daily.kz",
	"data.sale-inp.daily.ua":   "sale-in.daily.ua",
	"data.sale-out.monthly.kz": "sale-out.monthly.kz",
	"data.sale-out.monthly.ua": "sale-out.monthly.ua",
	"data.sale-out.weekly.ua":  "sale-out.weekly.ua",
	"data.sale-out.daily.by":   "sale-out.daily.by",
	"data.sale-out.daily.kz":   "sale-out.daily.kz",
	"data.sale-out.daily.ua":   "sale-out.daily.ua",
}

func checkHTag(t string) error {
	t = strings.ToLower(t)
	_, ok1 := convHTag[t]
	_, ok2 := listHTag[t]

	if ok1 || ok2 {
		return nil
	}

	return fmt.Errorf("api: invalid htag %s", t)
}

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
		err = minio.Put(bucketStreamErr, f, bytes.NewReader(m.marshalIndent()))
		if err != nil {
			log.Println("proc: err: save:", f, err)
		}

		err = minio.Copy(bucketStreamErr, o, b, o)
		if err != nil {
			log.Println("proc: err: copy:", o, err)
		}
	} else {
		f = makeFileName(m.Auth.ID, m.UUID, m.HTag)
		err = minio.Put(bucketStreamOut, f, d)
		if err != nil {
			log.Println("proc: err: save:", f, err)
		}

		log.Println("proc:", f, m.Proc, time.Since(s).String())
	}

	err = minio.Del(b, o)
	if err != nil {
		log.Println("proc: err: kill:", o, err)
	}

	err = zlog(m)
	if err != nil {
		log.Println("proc: err: zlog:", o, err)
	}
}

func procObject(r io.Reader) (jsonMeta, io.Reader, error) {
	m := jsonMeta{}

	meta, data, err := ungztarMetaData(r)
	if err != nil {
		return m, nil, err
	}

	m, err = unmarshalMeta(meta)
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

	t, err := gztarMetaData(m.marshal(), d)
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

func unmarshalData(data []byte, m *jsonMeta) (interface{}, error) {
	d := killUTF8BOM(data)
	m.ETag = btsToMD5(d)
	m.Size = int64(len(d))

	if strings.HasPrefix(m.CTag, magicConvString) {
		m.CTag = fmt.Sprintf("converted from %s format", m.HTag)
		m.HTag = convHTag[m.HTag]
		return unmarshalDataOLD(d, m)
	}

	return unmarshalDataNEW(d, m)
}

func unmarshalDataOLD(data []byte, m *jsonMeta) (interface{}, error) {
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

func unmarshalDataNEW(data []byte, m *jsonMeta) (interface{}, error) {
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

func mineLinks(v interface{}, m *jsonMeta) ([]byte, error) {
	t := m.HTag
	s := time.Now()

	a, err := getAuthREDIS([]string{m.Auth.ID})
	if err != nil {
		return nil, err
	}
	m.Auth = a[0]

	l, err := getAddrREDIS([]string{strToSHA1(makeMagicHead(m.Name, m.Head, m.Addr))})
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

	if isSaleIn(t) {
		if a, ok := v.(addrer); ok {
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

	lds, err := getDrugREDIS(keys)
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

	lds, err := getAddrREDIS(keys)
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
