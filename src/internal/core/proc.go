package core

import (
	"bytes"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"

	"internal/strutil"
)

const (
	fromBY = "by"
	fromKZ = "kz"
	fromRU = "ru"
	fromUA = "ua"

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

func isHTag(t string) bool {
	if _, ok := htags[t]; ok {
		return ok
	}

	return false
}

func proc(key string) error {
	s, err := s3cli.StatObject(backetStreamIn, key)
	if err != nil {
		return err
	}

	m, err := makeMeta(s.ContentType)
	if err != nil {
		return err
	}

	if !isHTag(m.HTag) {
		return fmt.Errorf("core: proc: invalid htag %s", m.HTag)
	}

	m.Time = s.LastModified.Format("02.01.2006 15:04:05.999999999")
	m.ETag = s.ETag // MD5
	m.Path = ""     // ?
	m.Size = s.Size
	if m.Name != "" {
		l, err := findLinkAddr(makeSHA1String(makeMagicAddr(m.Name, m.Head, m.Addr)))
		if err != nil {
			return err
		}
		m.Link = l[0]
	}

	o, err := s3cli.GetObject(backetStreamIn, key)
	if err != nil {
		return err
	}

	_, err = readMendClose(o)
	if err != nil {
		return err
	}

	// process object
	// put object backetStreamOut
	// put object backetStreamErr

	return nil
}

func Test(t string, b []byte) ([]byte, error) {
	var (
		src interface{}
		err error
	)

	if strings.Contains(t, "geo") {
		src = listDataGeoV3{}
	} else if strings.Contains(t, "daily.by") {
		src = listDataSaleBYV3{}
	} else if strings.Contains(t, "sale") {
		src = listDataSaleV3{}
	}

	if src == nil {
		return nil, fmt.Errorf("core: proc: invalid interface")
	}

	if err = json.Unmarshal(b, &src); err != nil {
		return nil, err
	}

	if err = mineLinkName(t, src.(nameLinker)); err != nil {
		return nil, err
	}

	if strings.Contains(t, "sale-in") {
		if err = mineLinkSupp(src.(suppLinker)); err != nil {
			return nil, err
		}
	}

	//
	//ext  =
	return json.Marshal(src)
}

func mineLinkName(t string, l nameLinker) error {
	var (
		keys = make([]string, l.len())
		from = filepath.Ext(t)
		name []byte
	)
	for i := 0; i < l.len(); i++ {
		switch {
		case from == fromUA:
			name = makeMagicDrugUA(l.getName(i))
		case from == fromRU:
			name = makeMagicDrugRU(l.getName(i))
		case from == fromKZ:
			name = makeMagicDrugKZ(l.getName(i))
		case from == fromBY:
			name = makeMagicDrugBY(l.getName(i))
		default:
			name = makeMagicDrug(l.getName(i))
		}
		keys[i] = makeSHA1String(name)
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

func mineLinkSupp(l suppLinker) error {
	var keys = make([]string, l.len())
	for i := 0; i < l.len(); i++ {
		keys[i] = makeSHA1String(makeMagicSupp(l.getSupp(i)))
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

func makeMagicAddr(name, head, addr string) []byte {
	return bytes.TrimSpace(
		[]byte(strutil.First(
			fmt.Sprintf("%s/%s: %s", name, head, addr),
			magicAddrLength,
		)),
	)
}

func makeMagicSupp(name string) []byte {
	return bytes.TrimSpace(
		[]byte(strutil.First(
			name,
			magicAddrLength,
		)),
	)
}

func makeMagicDrug(name string) []byte {
	return bytes.TrimSpace(
		[]byte(strutil.First(
			name,
			magicDrugLength,
		)),
	)
}

func makeMagicDrugBY(name string) []byte {
	return append(makeMagicDrug(name), magicSuffixBY...)
}

func makeMagicDrugKZ(name string) []byte {
	return append(makeMagicDrug(name), magicSuffixKZ...)
}

func makeMagicDrugRU(name string) []byte {
	return append(makeMagicDrug(name), magicSuffixRU...)
}

func makeMagicDrugUA(name string) []byte {
	return append(makeMagicDrug(name), magicSuffixUA...)
}

func makeSHA1String(b []byte) string {
	return fmt.Sprintf("%x", sha1.Sum(b))
}
