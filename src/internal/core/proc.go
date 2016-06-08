package core

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"

	"internal/gzutil"
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

var htags = map[string]struct{}{
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

func proc(data []byte) error {
	meta, data, err := popMetaData(data)
	if err != nil {
		return err
	}

	m, err := procMeta(meta, btsToMD5(data), int64(len(data)))
	if err != nil {
		return err
	}

	d, err := procData(m.HTag, data)
	if err != nil {
		return err
	}

	t, err := tarMetaData(m.marshalBase64(), d)
	if err != nil {
		return err
	}

	go func() { // ?
		err := minio.PutObject(backetStreamOut, m.UUID, t)
		if err != nil {
			// log.
		}
	}()
	//goToStreamErr(m.ID, ?) // FIXME
	return nil
}

func procMeta(data []byte, etag string, size int64) (jsonMeta, error) {
	m, err := unmarshalBase64meta(data)
	if err != nil {
		return m, err
	}

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

func checkHTag(t string) error {
	if _, ok := htags[strings.ToLower(t)]; ok {
		return nil
	}
	return fmt.Errorf("core: invalid htag %s", t)
}

func findLinkMeta(m jsonMeta) (*linkAddr, error) {
	l, err := getLinkAddr(strToSHA1(makeMagicHead(m.Name, m.Head, m.Addr)))
	if err != nil {
		return nil, err
	}
	return l[0], nil
}

func procData(htag string, data []byte) ([]byte, error) {
	var err error

	data, err = gzutil.Gunzip(data)
	if err != nil {
		return nil, err
	}

	data, err = mendIfUTF8(data)
	if err != nil {
		return nil, err
	}

	data, err = mineLinks(htag, data)
	if err != nil {
		return nil, err
	}

	return gzutil.Gzip(data)
}

func checkGzip(b []byte) error {
	if isTypeGzip(b) {
		return nil
	}
	return fmt.Errorf("core: content must contain gzip")
}

func mineLinks(t string, data []byte) ([]byte, error) {
	var src interface{}

	switch {
	case isGeo(t):
		src = jsonV3Geoa{}
	case isSaleBY(t):
		src = jsonV3SaleBy{}
	default:
		src = jsonV3Sale{}
	}

	err := json.Unmarshal(data, &src)
	if err != nil {
		return nil, err
	}

	err = mineLinkDrug(t, src.(linkDruger))
	if err != nil {
		return nil, err
	}

	if isSaleIn(t) {
		err = mineLinkAddr(src.(linkAddrer))
		if err != nil {
			return nil, err
		}
	}

	return json.Marshal(src)
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
