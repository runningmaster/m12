package core

import (
	"encoding/json"
	"net/http"

	"internal/redis"
)

type HTTPHeadReader interface {
	ReadHeader(http.Header)
}

type HTTPHeadWriter interface {
	WriteHeader(http.Header)
}

type Master interface {
	NewWorker() Worker
}

type Worker interface {
	Work([]byte) (interface{}, error)
}

type WorkFunc func([]byte) (interface{}, error)

func (f WorkFunc) Work(b []byte) (interface{}, error) {
	return f(b)
}

// Ping calls Redis PING
func Ping(_ []byte) (interface{}, error) {
	return redis.Ping()
}

// Info calls Redis INFO
func Info(_ []byte) (interface{}, error) {
	return redis.Info()
}

type jsonMeta struct {
	UUID string `json:"uuid,omitempty"`
	Auth string `json:"auth,omitempty"`
	Host string `json:"host,omitempty"`
	Time int64  `json:"time,omitempty"`
	User string `json:"user,omitempty"`

	HTag string `json:"htag,omitempty"` // *
	Spn1 int64  `json:"spn1,omitempty"` // *
	Spn2 int64  `json:"spn2,omitempty"` // *
	Nick string `json:"nick,omitempty"` // * BR_NICK:id_addr | Source:MDS_LICENSE | file:FileName (?) deprecated

	Name string `json:"name,omitempty"` // *
	Head string `json:"head,omitempty"` // *
	Addr string `json:"addr,omitempty"` // *
	Code string `json:"code,omitempty"` // egrpou (okpo)

	ETag string `json:"etag,omitempty"`
	Size int64  `json:"size,omitempty"`

	CTag string `json:"ctag,omitempty"`
	Test bool   `json:"test,omitempty"`

	Link linkAddr `json:"link,omitempty"` // ?
}

func unmarshalMeta(b []byte) (jsonMeta, error) {
	m := jsonMeta{}
	err := json.Unmarshal(b, &m)
	return m, err
}

func (m *jsonMeta) marshal() []byte {
	b, _ := json.Marshal(m)
	return b
}

type pair struct {
	Backet string `json:"backet,omitempty"`
	Object string `json:"object,omitempty"`
}

func (p pair) marshal() []byte {
	b, _ := json.Marshal(p)
	return b
}

func unmarshaPair(data []byte) (pair, error) {
	p := pair{}
	err := json.Unmarshal(data, &p)
	return p, err
}

func unmarshaPairExt(data []byte) (string, string, error) {
	p, err := unmarshaPair(data)
	if err != nil {
		return "", "", err
	}
	return p.Backet, p.Object, nil
}

// Redis scheme:
// SET => key="auth"
// SADD key v [v...]
// SREM key v [v...]
// SISMEMBER key v
// type auth string

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

type linkAddrer interface {
	ruler
	getSupp(int) string
	setLinkAddr(int, linkAddr)
}

type linkDruger interface {
	ruler
	getName(int) string
	setLinkDrug(int, linkDrug)
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

func (j jsonV3Geoa) setLinkDrug(i int, l linkDrug) {
	j[i].Link = l
}

func (j jsonV3Sale) len() int {
	return len(j)
}

func (j jsonV3Sale) getName(i int) string {
	return j[i].Name
}

func (j jsonV3Sale) setLinkDrug(i int, l linkDrug) {
	j[i].LinkDrug = l
}

func (j jsonV3Sale) getSupp(i int) string {
	return j[i].SuppName
}

func (j jsonV3Sale) setLinkAddr(i int, l linkAddr) {
	j[i].LinkAddr = l
}

func (j jsonV3SaleBy) len() int {
	return len(j)
}

func (j jsonV3SaleBy) getName(i int) string {
	return j[i].Name
}

func (j jsonV3SaleBy) setLinkDrug(i int, l linkDrug) {
	j[i].Link = l
}
