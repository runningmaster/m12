package core

import (
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

	"internal/strutil"

	minio "github.com/minio/minio-go"
	"github.com/spkg/bom"
)

var listHTag = map[string]struct{}{
	"geoapt.ru":           {},
	"geoapt.ua":           {},
	"sale-in.monthly.kz":  {},
	"sale-in.monthly.ua":  {},
	"sale-in.weekly.ua":   {},
	"sale-in.daily.kz":    {},
	"sale-in.daily.ua":    {},
	"sale-out.monthly.kz": {},
	"sale-out.monthly.ua": {},
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

	return fmt.Errorf("core: invalid htag %s", t)
}

func proc(data []byte) {
	bucket, object, err := unmarshaPairExt(data)
	if err != nil {
		log.Println("core: proc: err: pair:", object, err)
		return
	}
	defer func(t time.Time) {
		log.Println("core: proc:", object, time.Since(t).String())
	}(time.Now())

	o, err := cMINIO.GetObject(bucket, object)
	if err != nil {
		log.Println("core: proc: err: load:", object, err)
		return
	}
	defer func(c io.Closer) {
		if c != nil {
			_ = c.Close()
		}
	}(o)

	f, r, err := procObject(o)
	if err != nil {
		log.Println("core: proc: err:", object, err)
		err = cMINIO.CopyObject(bucketStreamErr, object, bucket+"/"+object, minio.NewCopyConditions())
		if err != nil {
			log.Println("core: proc: err: copy:", object, err)
		}
	} else {
		_, err = cMINIO.PutObject(bucketStreamOut, f, r, "")
		if err != nil {
			log.Println("core: proc: err: save:", f, err)
		}
	}

	err = cMINIO.RemoveObject(bucket, object)
	if err != nil {
		log.Println("core: proc: err: kill:", f, err)
	}
}

func procObject(r io.Reader) (string, io.Reader, error) {
	meta, data, err := ungztarMetaData(r)
	if err != nil {
		return "", nil, err
	}

	m, err := unmarshalMeta(meta)
	if err != nil {
		return "", nil, err
	}

	v, err := unmarshalData(data, &m)
	if err != nil {
		return "", nil, err
	}

	if r, ok := v.(ruler); ok {
		if r.len() == 0 {
			return "", nil, fmt.Errorf("no data")
		}
	}

	d, err := mineLinks(v, &m)
	if err != nil {
		return "", nil, err
	}

	t, err := gztarMetaData(m.marshal(), d)
	if err != nil {
		return "", nil, err
	}

	return makeFileName(m.UUID, m.Auth, m.HTag), t, nil
}

func mendIfUTF8(data []byte) ([]byte, error) {
	if strings.Contains(http.DetectContentType(data), "text/plain; charset=utf-8") {
		return bom.Clean(data), nil
	}
	return data, nil
}

const magicConvString = "conv"

func unmarshalData(data []byte, m *jsonMeta) (interface{}, error) {
	d, err := mendIfUTF8(data)
	if err != nil {
		return nil, err
	}

	m.ETag = btsToMD5(d)
	m.Size = int64(len(d))

	if strings.HasPrefix(m.CTag, magicConvString) {
		return unmarshalDataOLD(d, m)
	}

	return unmarshalDataNEW(d, m)
}

func unmarshalDataOLD(data []byte, m *jsonMeta) (interface{}, error) {
	t := m.HTag

	var v interface{}
	var err error
	switch {
	case isGeoV2(t):
		v, err = convGeo2(data, m)
	case isGeoV1(t):
		v, err = convGeo1(data, m)
	case isSaleBY(t):
		v, err = convSaleBy(data, m)
	default:
		v, err = convSale(data, m)
	}

	// cleanup
	m.HTag = convHTag[t]
	m.CTag = ""

	return v, err
}

func unmarshalDataNEW(data []byte, m *jsonMeta) (interface{}, error) {
	t := m.HTag

	var v interface{}
	switch {
	case isGeo(t):
		v = jsonV3Geoa{}
	case isSaleBY(t):
		v = jsonV3SaleBy{}
	default:
		v = jsonV3Sale{}
	}

	err := json.Unmarshal(data, &v)
	if err != nil {
		return nil, err
	}

	return v, err
}

func mineLinks(v interface{}, m *jsonMeta) ([]byte, error) {
	t := m.HTag

	var err error
	if d, ok := v.(druger); ok {
		err = mineDrugs(d, t)
	}
	if err != nil {
		return nil, err
	}

	if isSaleIn(t) {
		if a, ok := v.(addrer); ok {
			err = mineAddrs(a)
		}
		if err != nil {
			return nil, err
		}
	}

	if isGeo(t) {
		l, err := getAddr(strToSHA1(makeMagicHead(m.Name, m.Head, m.Addr)))
		if err != nil {
			return nil, err
		}
		m.Link = l[0]
	}

	return json.Marshal(v)
}

func mineDrugs(v druger, t string) error {
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

	lds, err := getDrug(keys...)
	if err != nil {
		return err
	}

	for i := 0; i < v.len(); i++ {
		v.setDrug(i, lds[i])
	}

	return nil
}

func mineAddrs(v addrer) error {
	var keys = make([]string, v.len())
	for i := 0; i < v.len(); i++ {
		keys[i] = strToSHA1(makeMagicAddr(v.getSupp(i)))
	}

	lds, err := getAddr(keys...)
	if err != nil {
		return err
	}

	for i := 0; i < v.len(); i++ {
		v.setAddr(i, lds[i])
	}

	return nil
}

const (
	magicSuffixBY = "{\"COUNTRY_ID\":\"1010\"}"
	magicSuffixKZ = "{\"COUNTRY_ID\":\"106\"}"
	magicSuffixRU = "{\"COUNTRY_ID\":\"1027\"}"
	magicSuffixUA = ""

	magicAddrLength = 1024
	magicDrugLength = 512
)

func makeMagicHead(name, head, addr string) string {
	return strings.TrimSpace(
		strutil.First(
			fmt.Sprintf("%s/%s: %s", name, head, addr),
			magicAddrLength,
		),
	)
}

func makeMagicAddr(name string) string {
	return strings.TrimSpace(
		strutil.First(
			name,
			magicAddrLength,
		),
	)
}

func makeMagicDrug(name string) string {
	return strings.TrimSpace(
		strutil.First(
			name,
			magicDrugLength,
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

func isGeoV1(s string) bool {
	return strings.Contains(s, "geostore")
}

func isGeoV2(s string) bool {
	return strings.Contains(s, "geoapt")
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
