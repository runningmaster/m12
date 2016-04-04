package core

/*
	//ip key tag
	// timestamp time
	//IHashcheck   string sha hash
*/

type headGeoaptV3 struct {
	//ID   string `json:"id"`
	//Name string `json:"name"`
	//Head string `json:"head"`
	//Addr string `json:"addr"`
	//Code string `json:"code"`// egrpou (okpo)

}

type bodyGeoaptV3 struct {
	Code  string  `json:"code,omitempty"`
	Name  string  `json:"name,omitempty"`
	Quant float64 `json:"quant,omitempty"`
	Price float64 `json:"price,omitempty"`
	URL   string  `json:"url,omitempty"` // formerly link -> addr, home, url (?)
}

type headSaleInV3 struct {
	//ID   string `json:"id"`
	//Name string `json:"name"`
	//Head string `json:"head"`
	//Addr string `json:"addr"`
	//Code string `json:"code"`// egrpou (okpo)
	// span []
}

type bodySaleInV3 struct {
	Code  string  `json:"code,omitempty"`
	Name  string  `json:"name,omitempty"`
	Quant float64 `json:"quant,omitempty"`
	Price float64 `json:"price,omitempty"`
	Stock float64 `json:"stock,omitempty"`

	//"company": "string",
	//"egrpou": "string",

	//	"rel_id": "string",    // ID_DRUG
	//	"rel_id_s": "string",  // ID_ADDRESS
	//	"rel_sha": "string",   // SHA1(name)
	//	"rel_sha_s": "string", // SHA1(supp)
}

type headSaleOutV3 struct {
	//ID   string `json:"id"`
	//Name string `json:"name"`
	//Head string `json:"head"`
	//Addr string `json:"addr"`
	//Code string `json:"code"`// egrpou (okpo)
	// span []
}

type bodySaleOutV3 struct {
	Code  string  `json:"code,omitempty"`
	Name  string  `json:"name,omitempty"`
	Quant float64 `json:"quant,omitempty"`
	Price float64 `json:"price,omitempty"`
	Stock float64 `json:"stock,omitempty"`
	Rmbrs int64   `json:"rmbrs,omitempty"` // Reimburse 0 / 1
	//
	QuantIn float64 `json:"quant_in,omitempty"`
	PriceIn float64 `json:"price_in,omitempty"`
}

type headSaleOutByV3 struct {
	// span []
	//TRangeLower string `json:",omitempty"`
	//TRangeUpper string `json:",omitempty"`
	//"Head": {
	//		"Source": "",  // "APTEKA","S",155
	//		"FileName": "" // Имя файла
	//	},
	/*
		Source    string `json:",omitempty"`
		Drugstore string `json:",omitempty"`
	*/

}

type bodySaleOutByV3 struct {
	Code     string  `json:"code,omitempty"`
	Name     string  `json:"name,omitempty"`     // formerly Drug
	QuantIn  float64 `json:"quant_in,omitempty"` // formerly QuantInp
	PriceIn  float64 `json:"price_in,omitempty"` // formerly PriceInp
	QuantOut float64 `json:"quant_out,omitempty"`
	PriceOut float64 `json:"price_out,omitempty"`
	PriceRoc float64 `json:"price_roc,omitempty"`
	Stock    float64 `json:"stock,omitempty"`   // formerly Balance
	StockT   float64 `json:"stock_t,omitempty"` // formerly BalanceT
}
