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
