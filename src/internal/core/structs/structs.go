package structs

import (
	"encoding/json"
	"fmt"
	"strings"
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

func unmarshalMeta(b []byte) (Meta, error) {
	m := Meta{}
	err := json.Unmarshal(b, &m)
	return m, err
}

func (m *Meta) marshal() []byte {
	b, _ := json.Marshal(m)
	return b
}

func (m *Meta) marshalIndent() []byte {
	b, _ := json.MarshalIndent(m, "", "\t")
	return b
}

// Redis scheme:
// HASH => key="stat"
// HMSET key i->n [i->n...]
// HMGET key i [i..]
type linkAuth struct {
	ID   string `json:"id,omitempty"   redis:"i"`
	Name string `json:"name,omitempty" redis:"n"`
}

// Redis scheme:
// HASH => key=ID (SHA1)
// HMSET key l/v a/v s/v e/v (if exists in json)
// HMGET key l a s e
// JSON array: [{"id":"key1","id_link":1,"id_addr":2,"id_stat":0,"egrpou":"egrpou1"}]
type linkAddr struct {
	ID     string `json:"id,omitempty"      redis:"key"`
	IDLink int64  `json:"id_link,omitempty" redis:"l"`
	IDAddr int64  `json:"id_addr,omitempty" redis:"a"`
	IDStat int64  `json:"id_stat,omitempty" redis:"s"`
	EGRPOU string `json:"egrpou,omitempty"  redis:"e"`
}

// Redis scheme:
// HASH => key=ID (SHA1)
// HMSET key l/v d/v b/v c/v s/v (if exists in json)
// HMGET key l d b c s
type linkDrug struct {
	ID     string `json:"id,omitempty"      redis:"key"`
	IDLink int64  `json:"id_link,omitempty" redis:"l"`
	IDDrug int64  `json:"id_drug,omitempty" redis:"d"`
	IDBrnd int64  `json:"id_brnd,omitempty" redis:"b"`
	IDCatg int64  `json:"id_catg,omitempty" redis:"c"`
	IDStat int64  `json:"id_stat,omitempty" redis:"s"`
}

// Redis scheme:
// HASH => key="stat"
// HMSET key i->n [i->n...]
// HMGET key i [i..]
type linkStat struct {
	ID   int64  `json:"id,omitempty"   redis:"i"`
	Name string `json:"name,omitempty" redis:"n"`
}

type itemV3Geoa struct {
	ID    string   `json:"id,omitempty"`
	Name  string   `json:"name,omitempty"`
	Home  string   `json:"home,omitempty"` // formerly link
	Quant float64  `json:"quant,omitempty"`
	Price float64  `json:"price,omitempty"`
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

type itemV3SaleBy struct {
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

type ruler interface {
	len() int
}

type addrer interface {
	ruler
	getSupp(int) string
	setAddr(int, linkAddr) bool
}

type druger interface {
	ruler
	getName(int) string
	setDrug(int, linkDrug) bool
}

type jsonV3Geoa []itemV3Geoa
type jsonV3Sale []itemV3Sale
type jsonV3SaleBy []itemV3SaleBy

func (j jsonV3Geoa) len() int {
	return len(j)
}

func (j jsonV3Geoa) getName(i int) string {
	return j[i].Name
}

func (j jsonV3Geoa) setDrug(i int, l linkDrug) bool {
	j[i].Link = l
	return l.IDLink != 0
}

func (j jsonV3Sale) len() int {
	return len(j)
}

func (j jsonV3Sale) getName(i int) string {
	return j[i].Name
}

func (j jsonV3Sale) setDrug(i int, l linkDrug) bool {
	j[i].LinkDrug = l
	return l.IDLink != 0
}

func (j jsonV3Sale) getSupp(i int) string {
	return j[i].SuppName
}

func (j jsonV3Sale) setAddr(i int, l linkAddr) bool {
	j[i].LinkAddr = l
	return l.IDLink != 0
}

func (j jsonV3SaleBy) len() int {
	return len(j)
}

func (j jsonV3SaleBy) getName(i int) string {
	return j[i].Name
}

func (j jsonV3SaleBy) setDrug(i int, l linkDrug) bool {
	j[i].Link = l
	return l.IDLink != 0
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
