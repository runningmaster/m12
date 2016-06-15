package core

import (
	"encoding/json"
	"fmt"
	"time"
)

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
		return nil, fmt.Errorf("core: conv: no data")
	}
	if len(v.Data[0].Item) == 0 {
		return nil, fmt.Errorf("core: conv: no data items")
	}

	d := make(jsonV3Sale, len(v.Data[0].Item))
	for i, v := range v.Data[0].Item {
		d[i].ID = v.Code
		d[i].Name = v.Drug
		d[i].QuantIn = v.QuantInp
		d[i].PriceIn = v.PriceInp
		d[i].QuantOut = v.QuantOut
		d[i].PriceOut = v.PriceOut
		d[i].Stock = v.Balance
		d[i].Reimburse = v.Reimburs != 0
		d[i].SuppName = v.Supp
		d[i].SuppCode = v.SuppOKPO
	}

	m.Spn1, err = convDateTimeToUnix(v.Meta.TRangeLower)
	if err != nil {
		return nil, err
	}
	m.Spn2, err = convDateTimeToUnix(v.Meta.TRangeUpper)
	if err != nil {
		return nil, err
	}
	m.Nick = v.Data[0].Head.Source
	if v.Data[0].Head.MDSLns != "" {
		m.Nick = m.Nick + ":" + v.Data[0].Head.MDSLns
	}

	return d, nil
}

func convDateTimeToUnix(s string) (int64, error) {
	t, err := time.Parse("02.01.2006 15:04:05", s)
	if err != nil {
		t, err = time.Parse("02.01.2006", s)
	}
	return t.Unix(), err
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
		return nil, fmt.Errorf("core: conv: no data")
	}
	if len(v.Data[0].Item) == 0 {
		return nil, fmt.Errorf("core: conv: no data items")
	}

	d := make(jsonV3SaleBy, len(v.Data[0].Item))
	for i, v := range v.Data[0].Item {
		d[i].ID = v.Code
		d[i].Name = v.Drug
		d[i].QuantIn = v.QuantInp
		d[i].PriceIn = v.PriceInp
		d[i].QuantOut = v.QuantOut
		d[i].PriceOut = v.PriceOut
		d[i].PriceRoc = v.PriceRoc
		d[i].Stock = v.Balance
		d[i].StockTab = v.BalanceT

	}

	m.Spn1, err = convDateTimeToUnix(v.Meta.TRangeLower)
	if err != nil {
		return nil, err
	}
	m.Spn2, err = convDateTimeToUnix(v.Meta.TRangeUpper)
	if err != nil {
		return nil, err
	}
	m.Nick = v.Data[0].Head.Source + ":" + v.Data[0].Head.Drugstore

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
		return nil, fmt.Errorf("core: conv: no data")
	}

	d := make(jsonV3Geoa, len(v.Data))
	for i, v := range v.Data {
		d[i].ID = v.Code
		d[i].Name = v.Name
		d[i].Home = v.Addr
		d[i].Quant = v.Quant
		d[i].Price = v.Price
		// workaround
		if v.Link != "" {
			d[i].Home = v.Link
		}
	}

	m.Name = v.Meta.Name
	m.Head = v.Meta.Head
	m.Addr = v.Meta.Addr
	m.Code = v.Meta.EGRPOU

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
		return nil, fmt.Errorf("core: conv: no data")
	}

	d := make(jsonV3Geoa, len(v.Data))
	for i, v := range v.Data {
		d[i].ID = v.ID
		d[i].Name = v.Name
		d[i].Home = v.Link
		d[i].Quant = v.Quant
		d[i].Price = v.Price
	}

	m.Name = v.Meta.Name
	m.Head = v.Meta.Head
	m.Addr = v.Meta.Addr
	m.Code = v.Meta.Code

	return d, nil
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
		Head   string
		Name   string
		Addr   string
		EGRPOU string
	}
	Data []struct {
		Code  string
		Name  string
		Desc  string
		Addr  string `json:",omitempty"`
		Link  string
		Price float64
		Quant float64
	}
}

type jsonV2Geoa struct {
	Meta struct {
		Name string `json:"name"`
		Head string `json:"head"`
		Addr string `json:"addr"`
		Code string `json:"code"`
	} `json:"meta"`
	Data []struct {
		ID    string  `json:"id"`
		Name  string  `json:"name"`
		Link  string  `json:"link"`
		Quant float64 `json:"quant"`
		Price float64 `json:"price"`
	} `json:"data"`
}
