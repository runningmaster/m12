package core

// Redis scheme:
// HASH => key="stat"
// HMSET key i->n [i->n...]
// HMGET key i [i..]
type LinkAuth struct {
	ID   string `json:"id,omitempty"   redis:"i"`
	Name string `json:"name,omitempty" redis:"n"`
}

// Redis scheme:
// HASH => key=ID (SHA1)
// HMSET key l/v a/v s/v e/v (if exists in json)
// HMGET key l a s e
// JSON array: [{"id":"key1","id_link":1,"id_addr":2,"id_stat":0,"egrpou":"egrpou1"}]
type LinkAddr struct {
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
type LinkDrug struct {
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
type LinkStat struct {
	ID   int64  `json:"id,omitempty"   redis:"i"`
	Name string `json:"name,omitempty" redis:"n"`
}

type itemV3Geoa struct {
	ID    string   `json:"id,omitempty"`
	Name  string   `json:"name,omitempty"`
	Home  string   `json:"home,omitempty"` // formerly link
	Quant float64  `json:"quant,omitempty"`
	Price float64  `json:"price,omitempty"`
	Link  LinkDrug `json:"link,omitempty"`
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
	LinkAddr  LinkAddr `json:"link_addr,omitempty"`
	LinkDrug  LinkDrug `json:"link_drug,omitempty"`
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
	Link     LinkDrug `json:"link,omitempty"`
}

type ruler interface {
	len() int
}

type addrer interface {
	ruler
	getSupp(int) string
	setAddr(int, LinkAddr) bool
}

type druger interface {
	ruler
	getName(int) string
	setDrug(int, LinkDrug) bool
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

func (j jsonV3Geoa) setDrug(i int, l LinkDrug) bool {
	j[i].Link = l
	return l.IDLink != 0
}

func (j jsonV3Sale) len() int {
	return len(j)
}

func (j jsonV3Sale) getName(i int) string {
	return j[i].Name
}

func (j jsonV3Sale) setDrug(i int, l LinkDrug) bool {
	j[i].LinkDrug = l
	return l.IDLink != 0
}

func (j jsonV3Sale) getSupp(i int) string {
	return j[i].SuppName
}

func (j jsonV3Sale) setAddr(i int, l LinkAddr) bool {
	j[i].LinkAddr = l
	return l.IDLink != 0
}

func (j jsonV3SaleBy) len() int {
	return len(j)
}

func (j jsonV3SaleBy) getName(i int) string {
	return j[i].Name
}

func (j jsonV3SaleBy) setDrug(i int, l LinkDrug) bool {
	j[i].Link = l
	return l.IDLink != 0
}

type jsonV1Sale struct {
	Meta struct {
		TRangeLower string
		TRangeUpper string
	}
	Data []struct {
		Head struct {
			Source string
			MDSLns string `json:",omitempty"`
		}
		Item []struct {
			Code     string
			Drug     string
			Supp     string  `json:",omitempty"`
			SuppOKPO string  `json:",omitempty"`
			QuantInp float64 `json:",omitempty"`
			PriceInp float64 `json:",omitempty"`
			QuantOut float64 `json:",omitempty"`
			PriceOut float64 `json:",omitempty"`
			Balance  float64 `json:",omitempty"`
			Reimburs int     `json:",omitempty"`
		}
	}
}

type jsonV1SaleBy struct {
	Meta struct {
		TRangeLower string
		TRangeUpper string
	}
	Data []struct {
		Head struct {
			Source    string
			Drugstore string
		}
		Item []struct {
			Code     string
			Drug     string
			QuantInp float64
			QuantOut float64
			PriceInp float64
			PriceOut float64
			PriceRoc float64
			Balance  float64
			BalanceT float64
		}
	}
}

type jsonV1Geoa struct {
	Meta struct {
		Name   string `json:"name,omitempty"`
		Head   string `json:"head,omitempty"`
		Addr   string `json:"addr,omitempty"`
		Code   string `json:"code,omitempty"`
		EGRPOU string `json:"EGRPOU,omitempty"` // deprecated from 1.0
	} `json:"meta"`
	Data []struct {
		ID    string  `json:"id,omitempty"`
		Code  string  `json:"Code,omitempty"` // deprecated from 1.0
		Name  string  `json:"name"`
		Desc  string  `json:"Desc,omitempty"` // deprecated from 1.0
		Addr  string  `json:"Addr,omitempty"` // deprecated from 1.0
		Link  string  `json:"link"`
		Quant float64 `json:"quant"`
		Price float64 `json:"price"`
	} `json:"data"`
}
