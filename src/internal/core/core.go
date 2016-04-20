package core

import (
	"net/http"

	"golang.org/x/net/context"
)

// Handler is func for processing data from api.
type (
	Handler func(context.Context, *http.Request) (interface{}, error)

	meta struct {
		ID string `json:"id,omitempty"` // ?
		IP string `json:"ip,omitempty"` // ?

		Auth string   `json:"auth,omitempty"` // *
		HTag string   `json:"htag,omitempty"` // *
		Nick string   `json:"nick,omitempty"` // * BR_NICK:id_addr | MDS_LICENSE / file:FileName (?) depecated
		Name string   `json:"name,omitempty"` // *
		Head string   `json:"head,omitempty"` // *
		Addr string   `json:"addr,omitempty"` // *
		Code string   `json:"code,omitempty"` // egrpou (okpo)
		Span []string `json:"span,omitempty"` // *

		Link linkAddr `json:"link,omitempty"` // ?

		Time string `json:"time,omitempty"` // ?
		ETag string `json:"etag,omitempty"` // ?
		Path string `json:"path,omitempty"` // ?
		Size int64  `json:"size,omitempty"` // ?

		SrcCE string `json:"src_ce,omitempty"` // Source ContentEncoding
		SrcCT string `json:"src_ct,omitempty"` // Source ContentType
	}

	itemGeoV3 struct {
		ID    string   `json:"id,omitempty"`
		Name  string   `json:"name,omitempty"`
		Quant float64  `json:"quant,omitempty"`
		Price float64  `json:"price,omitempty"`
		URL   string   `json:"url,omitempty"` // formerly link -> addr, home, url (?)
		Link  linkDrug `json:"link,omitempty"`
	}

	itemSaleV3 struct {
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

	itemSaleBYV3 struct {
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

	linkAddrer interface {
		len() int
		getSupp(int) string
		setLinkAddr(int, linkAddr)
	}

	linkDruger interface {
		len() int
		getName(int) string
		setLinkDrug(int, linkDrug)
	}

	listGeoV3    []itemGeoV3
	listSaleV3   []itemSaleV3
	listSaleBYV3 []itemSaleBYV3
)

func (l listGeoV3) len() int {
	return len(l)
}

func (l listGeoV3) getName(i int) string {
	return l[i].Name
}

func (l listGeoV3) setLinkDrug(i int, link linkDrug) {
	l[i].Link = link
}

func (l listSaleV3) len() int {
	return len(l)
}

func (l listSaleV3) getName(i int) string {
	return l[i].Name
}

func (l listSaleV3) setLinkDrug(i int, link linkDrug) {
	l[i].LinkDrug = link
}

func (l listSaleV3) getSupp(i int) string {
	return l[i].SuppName
}

func (l listSaleV3) setLinkAddr(i int, link linkAddr) {
	l[i].LinkAddr = link
}

func (l listSaleBYV3) len() int {
	return len(l)
}

func (l listSaleBYV3) getName(i int) string {
	return l[i].Name
}

func (l listSaleBYV3) setLinkDrug(i int, link linkDrug) {
	l[i].Link = link
}

// Init is caled from other package for manually initialization
func Init() error {
	err := initRedis()
	if err != nil {
		return err
	}

	return initS3Cli()
}

//$ curl --verbose --insecure --request 'POST' --header 'Content-Encoding: application/x-gzip' --header 'Content-Type: application/json; charset=utf-8' --header 'Content-Meta: eyJuYW1lIjoi0JDQv9GC0LXQutCwIDMiLCAiaGVhZCI6ItCR0IbQm9CQINCg0J7QnNCQ0KjQmtCQIiwiYWRkciI6ItCR0L7RgNC40YHQv9C+0LvRjCDRg9C7LiDQmtC40LXQstGB0LrQuNC5INCo0LvRj9GFLCA5OCIsImNvZGUiOiIxMjM0NTYifQ==' --upload-file 'data.json.gz' --user 'api:key-masterkey' --url http://localhost:8080/upload
