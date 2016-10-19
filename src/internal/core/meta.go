package core

import (
	"encoding/json"
	"fmt"
)

type meta struct {
	UUID string   `json:"uuid,omitempty"`
	Auth linkAuth `json:"auth,omitempty"`
	Host string   `json:"host,omitempty"`
	User string   `json:"user,omitempty"`
	Time string   `json:"time,omitempty"`
	Unix int64    `json:"unix,omitempty"`

	HTag string   `json:"htag,omitempty"` // *
	Span []string `json:"span,omitempty"` // *
	Nick string   `json:"nick,omitempty"` // * Source | Source:MDSLns | Source:Drugstore -> conv.go

	Name string `json:"name,omitempty"` // *
	Head string `json:"head,omitempty"` // *
	Addr string `json:"addr,omitempty"` // *
	Code string `json:"code,omitempty"` // egrpou (okpo)

	Link linkAddr `json:"link,omitempty"`

	CTag string `json:"ctag,omitempty"`
	ETag string `json:"etag,omitempty"`
	Size int64  `json:"size,omitempty"`
	Proc string `json:"proc,omitempty"`
	Fail string `json:"fail,omitempty"`
	Test bool   `json:"test,omitempty"`
}

func unmarshalMeta(b []byte) (*meta, error) {
	m := &meta{}
	err := json.Unmarshal(b, &m)
	return m, err
}

func (m *meta) marshal() []byte {
	b, _ := json.Marshal(m)
	return b
}

func (m *meta) marshalIndent() []byte {
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

func makeFileName(auth, uuid, htag string) string {
	return fmt.Sprintf("%s_%s_%s.tar", trimPart(auth), trimPart(uuid), htag)
}
