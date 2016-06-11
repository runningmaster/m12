package core

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"

	"internal/gzpool"
	"internal/minio"
	"internal/strutil"
)

const (
	extBY = ".by"
	extKZ = ".kz"
	extRU = ".ru"
	extUA = ".ua"

	magicSuffixBY = "{\"COUNTRY_ID\":\"1010\"}"
	magicSuffixKZ = "{\"COUNTRY_ID\":\"106\"}"
	magicSuffixRU = "{\"COUNTRY_ID\":\"1027\"}"
	magicSuffixUA = ""

	magicAddrLength = 1024
	magicDrugLength = 512
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
	"__ version 1 __":          "",
	"data.geostore":            "geoapt.ua",
	"data.sale-inp.monthly":    "sale-in.monthly.ua",
	"data.sale-inp.weekly":     "sale-in.weekly.ua",
	"data.sale-inp.daily":      "sale-in.daily.ua",
	"data.sale-out.monthly":    "sale-out.monthly.ua",
	"data.sale-out.weekly":     "sale-out.weekly.ua",
	"data.sale-out.daily":      "sale-out.daily.ua",
	"__ version 2 __":          "",
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
	if _, ok := listHTag[strings.ToLower(t)]; ok {
		return nil
	}
	return fmt.Errorf("core: invalid htag %s", t)
}

func proc(data []byte) error {
	meta, data, err := popMetaData(data)
	if err != nil {
		return err
	}

	m, err := procMeta(meta, btsToMD5(data), int64(len(data)))
	if err != nil {
		return err
	}

	d, err := procData(data, &m)
	if err != nil {
		return err
	}

	t, err := tarMetaData(m.marshalJSON(), d)
	if err != nil {
		return err
	}

	return minio.PutObject(backetStreamOut, m.UUID, t)
}

func procMeta(data []byte, etag string, size int64) (jsonMeta, error) {
	m, err := unmarshalJSONmeta(data)
	if err != nil {
		return m, err
	}
	/*
		t := m.HTag
		if s, ok := convHTag[t]; ok {
			m.HTag = s
		}

	*/
	err = checkHTag(m.HTag)
	if err != nil {
		return m, err
	}

	l, err := findLinkMeta(m)
	if err != nil {
		return m, err
	}

	m.Link = l
	m.ETag = etag
	m.Size = size

	return m, nil
}

func findLinkMeta(m jsonMeta) (linkAddr, error) {
	l, err := getLinkAddr(strToSHA1(makeMagicHead(m.Name, m.Head, m.Addr)))
	if err != nil {
		return linkAddr{}, err
	}
	return l[0], nil
}

func procData(data []byte, m *jsonMeta) ([]byte, error) {
	data, err := gzpool.Gunzip(data)
	if err != nil {
		return nil, err
	}

	data, err = mendIfUTF8(data)
	if err != nil {
		return nil, err
	}

	v, err := convData(data, m)
	if err != nil {
		return nil, err
	}

	data, err = mineLinks(m.HTag, v)
	if err != nil {
		return nil, err
	}

	return gzpool.Gzip(data)
}

func convData(data []byte, m *jsonMeta) (interface{}, error) {
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
	if err != nil {
		return nil, err
	}

	switch {
	case isGeo(t):
		v = jsonV3Geoa{}
	case isSaleBY(t):
		v = jsonV3SaleBy{}
	default:
		v = jsonV3Sale{}
	}

	err = json.Unmarshal(data, &v)
	if err != nil {
		return nil, err
	}

	return v, nil
}

func mineLinks(t string, v interface{}) ([]byte, error) {
	err := mineLinkDrug(t, v.(linkDruger))
	if err != nil {
		return nil, err
	}

	if isSaleIn(t) {
		err = mineLinkAddr(v.(linkAddrer))
		if err != nil {
			return nil, err
		}
	}

	return json.Marshal(v)
}

func mineLinkDrug(t string, l linkDruger) error {
	var (
		ext  = filepath.Ext(t)
		keys = make([]string, l.len())
		name string
	)
	for i := 0; i < l.len(); i++ {
		name = l.getName(i)
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

	lds, err := getLinkDrug(keys...)
	if err != nil {
		return err
	}

	for i := 0; i < l.len(); i++ {
		l.setLinkDrug(i, lds[i])
	}

	return nil
}

func mineLinkAddr(l linkAddrer) error {
	var keys = make([]string, l.len())
	for i := 0; i < l.len(); i++ {
		keys[i] = strToSHA1(makeMagicAddr(l.getSupp(i)))
	}

	lds, err := getLinkAddr(keys...)
	if err != nil {
		return err
	}

	for i := 0; i < l.len(); i++ {
		l.setLinkAddr(i, lds[i])
	}

	return nil
}

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
