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

	dataGeoV3 struct {
		ID    string   `json:"id,omitempty"`
		Name  string   `json:"name,omitempty"`
		Quant float64  `json:"quant,omitempty"`
		Price float64  `json:"price,omitempty"`
		URL   string   `json:"url,omitempty"` // formerly link -> addr, home, url (?)
		Link  linkDrug `json:"link,omitempty"`
	}

	dataSaleV3 struct {
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

	dataSaleBYV3 struct {
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

	suppLinker interface {
		len() int
		getSupp(int) string
		setLinkAddr(int, linkAddr)
	}

	nameLinker interface {
		len() int
		getName(int) string
		setLinkDrug(int, linkDrug)
	}

	listDataGeoV3    []dataGeoV3
	listDataSaleV3   []dataSaleV3
	listDataSaleBYV3 []dataSaleBYV3
)

func (l listDataGeoV3) len() int {
	return len(l)
}

func (l listDataGeoV3) getName(i int) string {
	return l[i].Name
}

func (l listDataGeoV3) setLinkDrug(i int, link linkDrug) {
	l[i].Link = link
}

func (l listDataSaleV3) len() int {
	return len(l)
}

func (l listDataSaleV3) getName(i int) string {
	return l[i].Name
}

func (l listDataSaleV3) setLinkDrug(i int, link linkDrug) {
	l[i].LinkDrug = link
}

func (l listDataSaleV3) getSupp(i int) string {
	return l[i].SuppName
}

func (l listDataSaleV3) setLinkAddr(i int, link linkAddr) {
	l[i].LinkAddr = link
}

func (l listDataSaleBYV3) len() int {
	return len(l)
}

func (l listDataSaleBYV3) getName(i int) string {
	return l[i].Name
}

func (l listDataSaleBYV3) setLinkDrug(i int, link linkDrug) {
	l[i].Link = link
}

// Init is caled from other package for manually initialization
func Init() error {
	var err error
	if err = initRedis(); err != nil {
		return err
	}

	if err = initS3Cli(); err != nil {
		return err
	}

	return nil
}

//$ curl --verbose --insecure --request 'POST' --header 'Content-Encoding: application/x-gzip' --header 'Content-Type: application/json; charset=utf-8' --header 'Content-Meta-JSON-Base64: eyJuYW1lIjoi0JDQv9GC0LXQutCwIDMiLCAiaGVhZCI6ItCR0IbQm9CQINCg0J7QnNCQ0KjQmtCQIiwiYWRkciI6ItCR0L7RgNC40YHQv9C+0LvRjCDRg9C7LiDQmtC40LXQstGB0LrQuNC5INCo0LvRj9GFLCA5OCIsImNvZGUiOiIxMjM0NTYifQ==' --upload-file 'data.json.gz' --user 'api:key-masterkey' --url http://localhost:8080/upload
