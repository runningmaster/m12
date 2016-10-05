package structs

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"internal/core/link"
)

const (
	BucketStreamIn  = "stream-in"
	BucketStreamOut = "stream-out"
	BucketStreamErr = "stream-err"

	SubjectSteamIn  = "m12." + BucketStreamIn
	SubjectSteamOut = "m12." + BucketStreamOut

	// should be move to pref
	ListN = 100
	TickD = 10 * time.Second
)

var (
	listHTag = map[string]struct{}{
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

	convHTag = map[string]string{
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
)

func FindHTag(t string) string {
	return convHTag[t]
}

func CheckHTag(t string) error {
	t = strings.ToLower(t)
	_, ok1 := convHTag[t]
	_, ok2 := listHTag[t]

	if ok1 || ok2 {
		return nil
	}

	return fmt.Errorf("meta: invalid htag %s", t)
}

type Meta struct {
	UUID string    `json:"uuid,omitempty"`
	Auth link.Auth `json:"auth,omitempty"`
	Host string    `json:"host,omitempty"`
	User string    `json:"user,omitempty"`
	Time string    `json:"time,omitempty"`
	Unix int64     `json:"unix,omitempty"`

	HTag string   `json:"htag,omitempty"` // *
	Span []string `json:"span,omitempty"` // *
	Nick string   `json:"nick,omitempty"` // * Source | Source:MDSLns | Source:Drugstore -> conv.go

	Name string `json:"name,omitempty"` // *
	Head string `json:"head,omitempty"` // *
	Addr string `json:"addr,omitempty"` // *
	Code string `json:"code,omitempty"` // egrpou (okpo)

	Link link.Addr `json:"link,omitempty"`

	CTag string `json:"ctag,omitempty"`
	ETag string `json:"etag,omitempty"`
	Size int64  `json:"size,omitempty"`
	Proc string `json:"proc,omitempty"`
	Fail string `json:"fail,omitempty"`
	Test bool   `json:"test,omitempty"`
}

func UnmarshalMeta(b []byte) (Meta, error) {
	m := Meta{}
	err := json.Unmarshal(b, &m)
	return m, err
}

func (m *Meta) Marshal() []byte {
	b, _ := json.Marshal(m)
	return b
}

func (m *Meta) MarshalIndent() []byte {
	b, _ := json.MarshalIndent(m, "", "\t")
	return b
}

const magicLen = 8

func trimPart(s string) string {
	if len(s) > magicLen {
		return s[:magicLen]
	}
	return s
}

func MakeFileName(auth, uuid, htag string) string {
	return fmt.Sprintf("%s_%s_%s.tar", trimPart(auth), trimPart(uuid), htag)
}
