package core

import (
	"encoding/base64"
	"encoding/json"
	"net/http"

	"golang.org/x/net/context"
)

// Handler is func for processing data from api.

type Handler func(context.Context, http.ResponseWriter, *http.Request) (interface{}, error)

type meta struct {
	ID string `json:"id,omitempty"` // ?
	IP string `json:"ip,omitempty"` // ?

	SKey string `json:"skey,omitempty"` // *
	HTag string `json:"htag,omitempty"` // *
	Nick string `json:"nick,omitempty"` // * BR_NICK:id_addr | MDS_LICENSE / file:FileName (?) depecated
	Name string `json:"name,omitempty"` // *
	Head string `json:"head,omitempty"` // *
	Addr string `json:"addr,omitempty"` // *
	Code string `json:"code,omitempty"` // egrpou (okpo)
	Spn1 int64  `json:"spn1,omitempty"` // *
	Spn2 int64  `json:"spn2,omitempty"` // *

	Link linkAddr `json:"link,omitempty"` // ?

	ETag string `json:"etag,omitempty"` // ?
	Size int64  `json:"size,omitempty"` // ?
	Time int64  `json:"time,omitempty"` // ?

	SrcCE string `json:"src_ce,omitempty"` // Source ContentEncoding
	SrcCT string `json:"src_ct,omitempty"` // Source ContentType
	Debug bool   `json:"debug,omitempty"`  // for debug purpose only
}

func makeMetaFromJSON(b []byte) (meta, error) {
	m := meta{}
	err := json.Unmarshal(b, &m)
	return m, err
}

func makeMetaFromBase64String(s string) (meta, error) {
	b, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return meta{}, err
	}
	return makeMetaFromJSON(b)
}

func (m meta) packToJSON() []byte {
	b, _ := json.Marshal(m)
	return b
}

func (m meta) packToBase64String() string {
	return base64.StdEncoding.EncodeToString(m.packToJSON())
}

type pair struct {
	Backet string `json:"backet,omitempty"`
	Object string `json:"object,omitempty"`
}

func makePairFromJSON(b []byte) (pair, error) {
	p := pair{}
	err := json.Unmarshal(b, &p)
	return p, err
}

func (p pair) packToJSON() []byte {
	b, _ := json.Marshal(p)
	return b
}

type itemV3Geoa struct {
	ID    string   `json:"id,omitempty"`
	Name  string   `json:"name,omitempty"`
	Quant float64  `json:"quant,omitempty"`
	Price float64  `json:"price,omitempty"`
	URL   string   `json:"url,omitempty"` // formerly link -> addr, home, url (?)
	Link  linkDrug `json:"link,omitempty"`
}

type itemV3Sale struct {
	ID        string   `json:"id,omitempty"`
	Name      string   `json:"name,omitempty"`
	QuantIn   float64  `json:"quant_in,omitempty"`
	PriceIn   float64  `json:"price_in,omitempty"`
	QuantOut  float64  `json:"quant_out,omitempty"`
	PriceOut  float64  `json:"price_out,omitempty"`
	Stock     float64  `json:"stock,omitempty"`
	Reimburse bool     `json:"reimburse,omitempty"`
	SuppName  string   `json:"supp_name,omitempty"`
	SuppCode  string   `json:"supp_code,omitempty"`
	LinkAddr  linkAddr `json:"link_addr,omitempty"`
	LinkDrug  linkDrug `json:"link_drug,omitempty"`
}

type itemV3Soby struct {
	ID       string   `json:"id,omitempty"`
	Name     string   `json:"name,omitempty"`
	QuantIn  float64  `json:"quant_in,omitempty"` // formerly QuantInp
	PriceIn  float64  `json:"price_in,omitempty"` // formerly PriceInp
	QuantOut float64  `json:"quant_out,omitempty"`
	PriceOut float64  `json:"price_out,omitempty"`
	PriceRoc float64  `json:"price_roc,omitempty"`
	Stock    float64  `json:"stock,omitempty"`     // formerly Balance
	StockTab float64  `json:"stock_tab,omitempty"` // formerly BalanceT
	Link     linkDrug `json:"link,omitempty"`
}

type linkAddrer interface {
	len() int
	getSupp(int) string
	setLinkAddr(int, linkAddr)
}

type linkDruger interface {
	len() int
	getName(int) string
	setLinkDrug(int, linkDrug)
}

type listV3Geoa []itemV3Geoa
type listV3Sale []itemV3Sale
type listV3Soby []itemV3Soby

func (l listV3Geoa) len() int {
	return len(l)
}

func (l listV3Geoa) getName(i int) string {
	return l[i].Name
}

func (l listV3Geoa) setLinkDrug(i int, link linkDrug) {
	l[i].Link = link
}

func (l listV3Sale) len() int {
	return len(l)
}

func (l listV3Sale) getName(i int) string {
	return l[i].Name
}

func (l listV3Sale) setLinkDrug(i int, link linkDrug) {
	l[i].LinkDrug = link
}

func (l listV3Sale) getSupp(i int) string {
	return l[i].SuppName
}

func (l listV3Sale) setLinkAddr(i int, link linkAddr) {
	l[i].LinkAddr = link
}

func (l listV3Soby) len() int {
	return len(l)
}

func (l listV3Soby) getName(i int) string {
	return l[i].Name
}

func (l listV3Soby) setLinkDrug(i int, link linkDrug) {
	l[i].Link = link
}

// Init is caled from other package for manually initialization
func Init() error {
	err := initRedis()
	if err != nil {
		return err
	}

	err = initS3Cli()
	if err != nil {
		return err
	}

	err = initNATSCli()
	if err != nil {
		return err
	}

	return nil
}
