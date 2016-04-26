package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

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
	"data.geoapt.ru":           {},
	"data.geoapt.ua":           {},
	"data.sale-in.monthly.kz":  {},
	"data.sale-in.monthly.ua":  {},
	"data.sale-in.weekly.ua":   {},
	"data.sale-in.daily.kz":    {},
	"data.sale-in.daily.ua":    {},
	"data.sale-out.monthly.kz": {},
	"data.sale-out.monthly.ua": {},
	"data.sale-out.weekly.ua":  {},
	"data.sale-out.daily.by":   {},
	"data.sale-out.daily.kz":   {},
	"data.sale-out.daily.ua":   {},
}

func proc(backet, object string) error {
	meta, data, err := mineMetaData(backet, object)
	if err != nil {
		return err
	}

	err = checkGzip(data)
	if err != nil {
		return err
	}

	m, err := makeMetaFromJSON(meta)
	if err != nil {
		return err
	}

	m.HTag = strings.ToLower(m.HTag)
	err = checkHTag(m.HTag)
	if err != nil {
		return err
	}

	m.Link, err = checkMeta(m)
	if err != nil {
		return err
	}

	m.ETag = btsToMD5(data)
	m.Size = int64(len(data))

	b, err := gunzip(data)
	if err != nil {
		return err
	}

	b, err = mendIfUTF8(data)
	if err != nil {
		return err
	}

	b, err = mineLinks(m.HTag, b)
	if err != nil {
		return err
	}

	t, err := tarMetaData(m.makeReadCloser(), ioutil.NopCloser(bytes.NewReader(data)))
	if err != nil {
		return err
	}

	goToStreamOut(m.ID, t)

	//goToStreamErr(m.ID, ?)

	return nil
}

func checkHTag(t string) error {
	if _, ok := htags[t]; ok {
		return nil
	}
	return fmt.Errorf("core: proc: invalid htag %s", t)
}

func checkGzip(b []byte) error {
	if isTypeGzip(b) {
		return nil
	}
	return fmt.Errorf("core: content must contain gzip")
}

func checkMeta(m meta) (linkAddr, error) {
	if m.Name == "" {
		return linkAddr{}, nil
	}
	l, err := findLinkAddr(strToSHA1(makeMagicHead(m.Name, m.Head, m.Addr)))
	if err != nil {
		return linkAddr{}, err
	}
	return l[0], nil
}

func mineLinks(t string, b []byte) ([]byte, error) {
	var src interface{}

	switch {
	case isGeo(t):
		src = listGeoV3{}
	case isSaleBY(t):
		src = listSaleBYV3{}
	default:
		src = listSaleV3{}
	}

	err := json.Unmarshal(b, &src)
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
		return fmt.Errorf("core: proc: invalid len (name): got %d, want %d ", len(lds), l.len())
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
		return fmt.Errorf("core: proc: invalid len (supp): got %d, want %d ", len(lds), l.len())
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
