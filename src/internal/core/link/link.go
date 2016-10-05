package link

// Redis scheme:
// HASH => key="stat"
// HMSET key i->n [i->n...]
// HMGET key i [i..]
type Auth struct {
	ID   string `json:"id,omitempty"   redis:"i"`
	Name string `json:"name,omitempty" redis:"n"`
}

// Redis scheme:
// HASH => key=ID (SHA1)
// HMSET key l/v a/v s/v e/v (if exists in json)
// HMGET key l a s e
// JSON array: [{"id":"key1","id_link":1,"id_addr":2,"id_stat":0,"egrpou":"egrpou1"}]
type Addr struct {
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
type Drug struct {
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
type Stat struct {
	ID   int64  `json:"id,omitempty"   redis:"i"`
	Name string `json:"name,omitempty" redis:"n"`
}

type ItemV3Geoa struct {
	ID    string  `json:"id,omitempty"`
	Name  string  `json:"name,omitempty"`
	Home  string  `json:"home,omitempty"` // formerly link
	Quant float64 `json:"quant,omitempty"`
	Price float64 `json:"price,omitempty"`
	Link  Drug    `json:"link,omitempty"`
}

type ItemV3Sale struct {
	ID        string  `json:"id,omitempty"`
	Name      string  `json:"name,omitempty"`
	QuantIn   float64 `json:"quant_in,omitempty"`
	PriceIn   float64 `json:"price_in,omitempty"`
	QuantOut  float64 `json:"quant_out,omitempty"`
	PriceOut  float64 `json:"price_out,omitempty"`
	Stock     float64 `json:"stock,omitempty"`
	Reimburse bool    `json:"reimburse,omitempty"`
	SuppName  string  `json:"supp_name,omitempty"`
	SuppCode  string  `json:"supp_code,omitempty"`
	LinkAddr  Addr    `json:"link_addr,omitempty"`
	LinkDrug  Drug    `json:"link_drug,omitempty"`
}

type ItemV3SaleBy struct {
	ID       string  `json:"id,omitempty"`
	Name     string  `json:"name,omitempty"`
	QuantIn  float64 `json:"quant_in,omitempty"` // formerly QuantInp
	PriceIn  float64 `json:"price_in,omitempty"` // formerly PriceInp
	QuantOut float64 `json:"quant_out,omitempty"`
	PriceOut float64 `json:"price_out,omitempty"`
	PriceRoc float64 `json:"price_roc,omitempty"`
	Stock    float64 `json:"stock,omitempty"`     // formerly Balance
	StockTab float64 `json:"stock_tab,omitempty"` // formerly BalanceT
	Link     Drug    `json:"link,omitempty"`
}

type Ruler interface {
	Len() int
}

type Addrer interface {
	Ruler
	GetSupp(int) string
	SetAddr(int, Addr) bool
}

type Druger interface {
	Ruler
	GetName(int) string
	SetDrug(int, Drug) bool
}

type DataV3Geoa []ItemV3Geoa
type DataV3Sale []ItemV3Sale
type DataV3SaleBy []ItemV3SaleBy

func (d DataV3Geoa) Len() int {
	return len(d)
}

func (d DataV3Geoa) GetName(i int) string {
	return d[i].Name
}

func (d DataV3Geoa) SetDrug(i int, l Drug) bool {
	d[i].Link = l
	return l.IDLink != 0
}

func (d DataV3Sale) len() int {
	return len(d)
}

func (d DataV3Sale) GetName(i int) string {
	return d[i].Name
}

func (d DataV3Sale) SetDrug(i int, l Drug) bool {
	d[i].LinkDrug = l
	return l.IDLink != 0
}

func (d DataV3Sale) GetSupp(i int) string {
	return d[i].SuppName
}

func (d DataV3Sale) SetAddr(i int, l Addr) bool {
	d[i].LinkAddr = l
	return l.IDLink != 0
}

func (d DataV3SaleBy) Len() int {
	return len(d)
}

func (d DataV3SaleBy) GetName(i int) string {
	return d[i].Name
}

func (d DataV3SaleBy) SetDrug(i int, l Drug) bool {
	d[i].Link = l
	return l.IDLink != 0
}

type DataV1Sale struct {
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

type DataV1SaleBy struct {
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

type DataV1Geoa struct {
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
