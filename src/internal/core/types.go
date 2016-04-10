package core

type (
	// Redis scheme:
	// SET => key="auth"
	// SADD key v [v...]
	// SREM key v [v...]
	// SISMEMBER key v

	// Redis scheme:
	// HASH => key=ID (SHA1)
	// HMSET key l/v a/v s/v e/v (if exists in json)
	// HMGET key l a s e
	// JSON array: [{"id":"key1","id_link":1,"id_addr":2,"id_stat":0,"egrpou":"egrpou1"}]
	linkAddr struct {
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
	linkDrug struct {
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
	linkStat struct {
		ID   int64  `json:"id,omitempty"   redis:"i"`
		Name string `json:"name,omitempty" redis:"n"`
	}
)

func (l linkAddr) keyflds(k string) []interface{} {
	return []interface{}{k, "l", "a", "s", "e"}
}

func (l linkAddr) keyvals() []interface{} {
	v := make([]interface{}, 0, 1+4*2)
	v = append(v, l.ID)
	if l.IDLink != 0 {
		v = append(v, "l", l.IDLink)
	}
	if l.IDAddr != 0 {
		v = append(v, "a", l.IDAddr)
	}
	if l.IDStat != 0 {
		v = append(v, "s", l.IDStat)
	}
	if l.EGRPOU != "" {
		v = append(v, "e", l.EGRPOU)
	}
	return v
}

func isEmpty(v []interface{}) bool {
	for i := range v {
		if v[i] != nil {
			return false
		}
	}
	return true
}

func (l linkAddr) makeFrom(k string, v []interface{}) interface{} {
	if isEmpty(v) {
		return nil
	}
	return linkAddr{
		ID:     k,
		IDLink: toInt64(v[0]),  // "l"
		IDAddr: toInt64(v[1]),  // "a"
		IDStat: toInt64(v[2]),  // "s"
		EGRPOU: toString(v[3]), // "e"
	}
}

func (l linkDrug) keyflds(k string) []interface{} {
	return []interface{}{k, "l", "d", "b", "c", "s"}
}

func (l linkDrug) keyvals() []interface{} {
	v := make([]interface{}, 0, 5*2)
	v = append(v, l.ID)
	if l.IDLink != 0 {
		v = append(v, "l", l.IDLink)
	}
	if l.IDDrug != 0 {
		v = append(v, "d", l.IDDrug)
	}
	if l.IDBrnd != 0 {
		v = append(v, "b", l.IDBrnd)
	}
	if l.IDCatg != 0 {
		v = append(v, "c", l.IDCatg)
	}
	if l.IDStat != 0 {
		v = append(v, "s", l.IDStat)
	}
	return v
}

func (l linkDrug) makeFrom(k string, v []interface{}) interface{} {
	if isEmpty(v) {
		return nil
	}
	return linkDrug{
		ID:     k,
		IDLink: toInt64(v[0]), // "l"
		IDDrug: toInt64(v[1]), // "d"
		IDBrnd: toInt64(v[2]), // "b"
		IDCatg: toInt64(v[3]), // "c"
		IDStat: toInt64(v[4]), // "s"
	}
}

func (l linkStat) makeFrom(k int64, v interface{}) interface{} {
	if v == nil {
		return v
	}
	return linkStat{
		ID:   k,
		Name: toString(v), // "n"
	}
}

// geo	key, tag, n,h,a,c
// sle	key, tag, src, span
type нead struct {
	ID string `json:"id,omitempty"` // ?
	IP string `json:"ip,omitempty"` // ?

	Key  string   `json:"key,omitempty"`  // *
	Tag  string   `json:"tag,omitempty"`  // *
	Src  string   `json:"src,omitempty"`  // * BR_NICK:id_addr | MDS_LICENSE / file:FileName (?) depecated
	Name string   `json:"name,omitempty"` // *
	Head string   `json:"head,omitempty"` // *
	Addr string   `json:"addr,omitempty"` // *
	Code string   `json:"code,omitempty"` // egrpou (okpo)
	Span []string `json:"span,omitempty"` // *

	From int64  `json:"from,omitempty"` // ?
	Time string `json:"time,omitempty"` // ?
	Hash string `json:"hash,omitempty"` // ?
	Path string `json:"path,omitempty"` // ?
}

/*

type bodyGeoV3 struct {
	ID    string   `json:"id,omitempty"`
	Name  string   `json:"name,omitempty"`
	Quant float64  `json:"quant,omitempty"`
	Price float64  `json:"price,omitempty"`
	URL   string   `json:"url,omitempty"` // formerly link -> addr, home, url (?)
	Link  linkDrug `json:"link,omitempty"`
}

type bodySaleV3 struct {
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

type bodySaleBYV3 struct {
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
*/
