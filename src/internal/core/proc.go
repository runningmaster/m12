package core

import (
	"crypto/sha1"
	"fmt"
	"strings"

	"internal/strutil"
)

const (
	htagExtBY = "by"
	htagExtKZ = "kz"
	htagExtRU = "ru"

	magicAddrLength = 1024
	magicDrugLength = 512

	magicSuffixBY = "{\"COUNTRY_ID\":\"1010\"}"
	magicSuffixKZ = "{\"COUNTRY_ID\":\"106\"}"
	magicSuffixRU = "{\"COUNTRY_ID\":\"1027\"}"
)

/*
CheckGzip (if not then fail)
GetHeader [if not exist then find query param (Form)]
Calc body hash md5
sha1:	af064923bbf2301596aac4c273ba32178ebc4a96
md5:	b0804ec967f48520697662a204f5fe72
Put "H.Tag" "H.UUID"+.gz
SendToChannel
*/

func proc(key string) error {
	s, err := s3cli.StatObject(backetStreamIn, key)
	if err != nil {
		return err
	}

	m, err := makeMeta(s.ContentType)
	if err != nil {
		return err
	}

	m.Time = s.LastModified.Format("02.01.2006 15:04:05.999999999")
	m.Path = backetStreamIn + "/" + key
	m.ETag = s.ETag
	m.Size = s.Size
	//if m.Ownr, err = makeOwner(m.Name, m.Head, m.Addr); err != nil {
	//
	//}

	//

	_, err = s3cli.GetObject(backetStreamIn, key)
	if err != nil {
		return err
	}
	return nil
}

func makeOwner(name, head, addr string) (int64, error) {
	if name == "" {
		return 0, fmt.Errorf("proc: name not found")
	}
	_ = fmt.Sprintf("%x", sha1.Sum([]byte(makeMagicAddr(name, head, addr))))
	return 0, nil
}

func makeMagicAddr(name, head, addr string) string {
	return strings.TrimSpace(
		strutil.First(
			fmt.Sprintf("%s/%s: %s", name, head, addr),
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
