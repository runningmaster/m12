package core

import (
	"encoding/base64"
	"encoding/json"
	"net/http"

	"internal/redis"
	"internal/s3"
)

const (
	backetStreamIn  = "stream-in"
	backetStreamOut = "stream-out"
	backetStreamErr = "stream-err"

	subjectSteamIn  = backetStreamIn + ".67a7ea16"
	subjectSteamOut = backetStreamOut + ".0566ce58"
)

type HTTPHeadReader interface {
	ReadHeader(http.Header)
}

type HTTPHeadWriter interface {
	WriteHeader(http.Header)
}

type Worker interface {
	Work([]byte) (interface{}, error)
}

type WorkFunc func([]byte) (interface{}, error)

func (f WorkFunc) Work(b []byte) (interface{}, error) {
	return f(b)
}

func Init() error {
	return s3.InitBacketList(backetStreamIn, backetStreamOut, backetStreamErr)
}

/*



func goNotifyStream(n int) {
	go func() {
		c := time.Tick(1 * time.Second)
		var err error
		for _ = range c {
			err = notifyStream(backetStreamIn, subjectSteamIn, n)
			if err != nil {
				log.Println(err)
			}
			err = notifyStream(backetStreamOut, subjectSteamOut, n)
			if err != nil {
				log.Println(err)
			}
		}
	}()
}

func notifyStream(backet, subject string, n int) error {
	objs, err := s3.listObjectsN(backet, "", false, n)
	if err != nil {
		return err
	}
	return nil
}

*/

// Ping calls Redis PING
func Ping(_ []byte) (interface{}, error) {
	return redis.Ping()
}

// Info calls Redis INFO
func Info(_ []byte) (interface{}, error) {
	return redis.Info()
}

type meta struct {
	UUID string `json:"uuid,omitempty"`
	Host string `json:"host,omitempty"`
	Auth string `json:"auth,omitempty"`
	Time int64  `json:"time,omitempty"`
	Test bool   `json:"test,omitempty"`

	HTag string `json:"htag,omitempty"` // *
	Spn1 int64  `json:"spn1,omitempty"` // *
	Spn2 int64  `json:"spn2,omitempty"` // *
	Nick string `json:"nick,omitempty"` // * BR_NICK:id_addr | MDS_LICENSE / file:FileName (?) depecated

	Name string `json:"name,omitempty"` // *
	Head string `json:"head,omitempty"` // *
	Addr string `json:"addr,omitempty"` // *
	Code string `json:"code,omitempty"` // egrpou (okpo)

	Link linkAddr `json:"link,omitempty"` // ?

	ETag string `json:"etag,omitempty"`
	Size int64  `json:"size,omitempty"`
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
