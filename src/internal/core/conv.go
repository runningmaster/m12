package core

import (
	"fmt"
	"net/http"
)

var Conv = &convWorker{}

type convWorker struct {
	meta []byte
}

func (w *convWorker) ReadHeader(h http.Header) {
	w.meta = []byte(h.Get("Content-Meta"))
}

func (w *convWorker) Work(data []byte) (interface{}, error) {

	//	go func() { // ?
	//		err := minio.PutObject(backetStreamIn, uuid, t)
	//		if err != nil {
	//			// log.
	//		}
	//	}()

	return "uuid", nil
}

func convSale(src *jsonV1Sale) (*jsonMeta, jsonV3Sale, error) {
	if len(src.Data) == 0 {
		return nil, nil, fmt.Errorf("core: conv data: no data")
	}
	if len(src.Data[0].Item) == 0 {
		return nil, nil, fmt.Errorf("core: conv data: no data items")
	}

	data := make(jsonV3Sale, len(src.Data[0].Item))
	for i := range src.Data[0].Item {
		_ = src.Data[0].Item[i]
	}
	return nil, data, nil
}

func convSaleBy(src *jsonV1SaleBy) (*jsonMeta, jsonV3SaleBy, error) {
	if len(src.Data) == 0 {
		return nil, nil, fmt.Errorf("core: conv data: no data")
	}
	if len(src.Data[0].Item) == 0 {
		return nil, nil, fmt.Errorf("core: conv data: no data items")
	}

	data := make(jsonV3SaleBy, len(src.Data[0].Item))
	for i := range src.Data[0].Item {
		_ = src.Data[0].Item[i]
	}
	return nil, data, nil
}

func convGeo1(src *jsonV1Geoa) (*jsonMeta, jsonV3Geoa, error) {
	if len(src.Data) == 0 {
		return nil, nil, fmt.Errorf("core: conv data: no data")
	}

	data := make(jsonV3Geoa, len(src.Data))
	for i := range src.Data {
		_ = src.Data[i]
	}
	return nil, data, nil
}

func convGeo2(src *jsonV2Geoa) (*jsonMeta, jsonV3Geoa, error) {
	if len(src.Data) == 0 {
		return nil, nil, fmt.Errorf("core: conv data: no data")
	}

	data := make(jsonV3Geoa, len(src.Data))
	for i := range src.Data {
		_ = src.Data[i]
	}
	return nil, data, nil
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
