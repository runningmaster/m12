package core

import (
	"encoding/json"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"internal/gzutil"
	"internal/s3"
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

func findLinkAddr(keys ...string) ([]linkAddr, error) {
	return nil, nil
}

func findLinkDrug(keys ...string) ([]linkDrug, error) {
	return nil, nil
}

func popObjectByJSON(data []byte) (io.Reader, error) {
	p, err := unmarshaJSONpair(data)
	if err != nil {
		return nil, err
	}

	return s3.PopObject(p.Backet, p.Object)
}

func proc(data []byte) error {
	o, err := popObjectByJSON(data)
	if err != nil {
		return err
	}

	meta, data, err := untarMetaData(o)
	if err != nil {
		return err
	}

	m, err := procMeta(meta, btsToMD5(data), int64(len(data)))
	if err != nil {
		return err
	}

	d, err := procData(data)
	if err != nil {
		return err
	}

	t, err := tarMetaData(m.marshalBase64(), d)
	if err != nil {
		return err
	}

	go func() { // ?
		err := s3.PutObject(backetStreamOut, m.UUID, t)
		if err != nil {
			// log.
		}
	}()
	//goToStreamOut(m.ID, t)

	//goToStreamErr(m.ID, ?)

	return nil
}

func procMeta(data []byte, etag string, size int64) (meta, error) {
	m, err := unmarshalBase64meta(data)
	if err != nil {
		return meta{}, err
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

func findLinkMeta(m meta) (linkAddr, error) {
	if m.Name == "" {
		return linkAddr{}, nil
	}
	l, err := findLinkAddr(strToSHA1(makeMagicHead(m.Name, m.Head, m.Addr)))
	if err != nil {
		return linkAddr{}, err
	}
	return l[0], nil
}

func procData(data []byte) ([]byte, error) {
	var err error

	data, err = gunzip(data)
	if err != nil {
		return nil, err
	}

	data, err = mendIfUTF8(data)
	if err != nil {
		return nil, err
	}

	data, err = mineLinks("m.HTag", data)
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
		src = listV3Geoa{}
	case isSaleBY(t):
		src = listV3Soby{}
	default:
		src = listV3Sale{}
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

	lds, err := findLinkDrug(keys...)
	if err != nil {
		return err
	}

	if len(lds) != l.len() {
		return fmt.Errorf("core: proc: invalid len (name): got %d, want %d", len(lds), l.len())
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

	lds, err := findLinkAddr(keys...)
	if err != nil {
		return err
	}

	if len(lds) != l.len() {
		return fmt.Errorf("core: proc: invalid len (supp): got %d, want %d", len(lds), l.len())
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
