package core

import (
	"encoding/json"
	"fmt"
	"net/http"

	"internal/gzutil"
)

var Conv = &convWorker{}

type convWorker struct {
	meta []byte
	uuid string
}

func (w *convWorker) ReadHeader(h http.Header) {
	w.meta = []byte(h.Get("Content-Meta"))
	w.uuid = h.Get("Content-UUID")
}

func (w *convWorker) Work(data []byte) (interface{}, error) {
	m, err := unmarshalJSONmeta(w.meta)
	if err != nil {
		return nil, err
	}

	t := m.HTag
	var v interface{}
	switch {
	case isGeoV2(t):
		v, err = convGeo2(data, &m)
	case isGeoV1(t):
		v, err = convGeo1(data, &m)
	case isSaleBY(t):
		v, err = convSaleBy(data, &m)
	default:
		v, err = convSale(data, &m)
	}
	if err != nil {
		return nil, err
	}

	data, err = json.Marshal(v)
	if err != nil {
		return nil, err
	}

	data, err = gzutil.Gzip(data)
	if err != nil {
		return nil, err
	}

	m.HTag = convHTag[t]

	putd := &putdWorker{m.marshalJSON(), w.uuid}
	return putd.Work(data)
}

func unmarshalSale(data []byte) (*jsonV1Sale, error) {
	v := &jsonV1Sale{}

	err := json.Unmarshal(data, &v)
	if err != nil {
		return nil, err
	}

	return v, nil
}

func convSale(data []byte, m *jsonMeta) (jsonV3Sale, error) {
	v, err := unmarshalSale(data)
	if err != nil {
		return nil, err
	}

	if len(v.Data) == 0 {
		return nil, fmt.Errorf("core: conv data: no data")
	}
	if len(v.Data[0].Item) == 0 {
		return nil, fmt.Errorf("core: conv data: no data items")
	}

	d := make(jsonV3Sale, len(v.Data[0].Item))
	for i := range v.Data[0].Item {
		_ = v.Data[0].Item[i]
	}

	return d, nil
}

func unmarshalSaleBy(data []byte) (*jsonV1SaleBy, error) {
	v := &jsonV1SaleBy{}

	err := json.Unmarshal(data, &v)
	if err != nil {
		return nil, err
	}

	return v, nil
}

func convSaleBy(data []byte, m *jsonMeta) (jsonV3SaleBy, error) {
	v, err := unmarshalSaleBy(data)
	if err != nil {
		return nil, err
	}

	if len(v.Data) == 0 {
		return nil, fmt.Errorf("core: conv data: no data")
	}
	if len(v.Data[0].Item) == 0 {
		return nil, fmt.Errorf("core: conv data: no data items")
	}

	d := make(jsonV3SaleBy, len(v.Data[0].Item))
	for i := range v.Data[0].Item {
		_ = v.Data[0].Item[i]
	}

	return d, nil
}

func unmarshalGeo1(data []byte) (*jsonV1Geoa, error) {
	v := &jsonV1Geoa{}

	err := json.Unmarshal(data, &v)
	if err != nil {
		return nil, err
	}

	return v, nil
}

func convGeo1(data []byte, m *jsonMeta) (jsonV3Geoa, error) {
	v, err := unmarshalGeo1(data)
	if err != nil {
		return nil, err
	}

	if len(v.Data) == 0 {
		return nil, fmt.Errorf("core: conv data: no data")
	}

	d := make(jsonV3Geoa, len(v.Data))
	for i := range v.Data {
		_ = v.Data[i]
	}

	return d, nil
}

func unmarshalGeo2(data []byte) (*jsonV2Geoa, error) {
	v := &jsonV2Geoa{}

	err := json.Unmarshal(data, &v)
	if err != nil {
		return nil, err
	}

	return v, nil
}

func convGeo2(data []byte, m *jsonMeta) (jsonV3Geoa, error) {
	v, err := unmarshalGeo2(data)
	if err != nil {
		return nil, err
	}

	if len(v.Data) == 0 {
		return nil, fmt.Errorf("core: conv data: no data")
	}

	d := make(jsonV3Geoa, len(v.Data))
	for i := range v.Data {
		_ = v.Data[i]
	}

	return d, nil
}

type jsonV1Sale struct {
	Meta struct {
		Version     int
		Agent       string `json:",omitempty"`
		Timestamp   string
		TRangeLower string
		TRangeUpper string
		AccountType string `json:",omitempty"`
		ISender     string `json:",omitempty"`
		IHashtag    string `json:",omitempty"`
		ITimestamp  string `json:",omitempty"`
		IHashstamp  string `json:",omitempty"`
		IHashcheck  string `json:",omitempty"`
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
			IDrugSHA string  `json:",omitempty"`
			IDrugLNK string  `json:",omitempty"`
			ISuppSHA string  `json:",omitempty"`
			ISuppLNK string  `json:",omitempty"`
		}
	}
}

type jsonV1SaleBy struct {
	Meta struct {
		Version     int
		Agent       string
		Timestamp   string
		TRangeLower string
		TRangeUpper string
		ISender     string
		IHashtag    string
		ITimestamp  string
		IHashstamp  string
		IHashcheck  string
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
			IDrugSHA string
			IDrugLNK string
		}
	}
}

type jsonV1Geoa struct {
	Meta struct {
		Code         string
		Head         string
		Name         string
		Addr         string
		EGRPOU       string
		Axioma       string // Name/Head: Addr
		AxiomaSHA    string
		AxiomaCode   string
		DateTimeZone string
		Agent        string
		Version      string
		ISender      string
		IHashtag     string
		ITimestamp   string
		IHashstamp   string
		IHashcheck   string
	}
	Data []struct {
		Code     string
		Name     string
		Desc     string
		Addr     string `json:",omitempty"`
		Link     string
		Price    float64
		Quant    float64
		IDrugSHA string
		IDrugLNK string
	}
}

type jsonV2Geoa struct {
	Meta struct {
		ID      string `json:"id"`
		Name    string `json:"name"`
		Head    string `json:"head"`
		Addr    string `json:"addr"`
		Code    string `json:"code"`
		Time    string `json:"time"`
		RelID   string `json:"rel_id,omitempty"`
		RelSHA  string `json:"rel_sha,omitempty"`
		SkyIP   string `json:"sky_ip"`
		SkyKey  string `json:"sky_key"`
		SkyTag  string `json:"sky_tag"`
		SkySHA  string `json:"sky_sha"`
		SkyTime string `json:"sky_time"`
	} `json:"meta"`
	Data []struct {
		ID     string  `json:"id"`
		Name   string  `json:"name"`
		Link   string  `json:"link"`
		Quant  float64 `json:"quant"`
		Price  float64 `json:"price"`
		RelID  string  `json:"rel_id,omitempty"`
		RelSHA string  `json:"rel_sha,omitempty"`
	} `json:"data"`
}
