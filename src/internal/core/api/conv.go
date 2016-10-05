package api

import (
	"encoding/json"
	"fmt"
	"time"

	"internal/core/link"
	"internal/core/structs"
)

func unmarshalSale(data []byte) (*link.DataV1Sale, error) {
	v := &link.DataV1Sale{}

	err := json.Unmarshal(data, &v)
	if err != nil {
		return nil, err
	}

	return v, nil
}

func convSale(data []byte, m *structs.Meta) (link.DataV3Sale, error) {
	v, err := unmarshalSale(data)
	if err != nil {
		return nil, err
	}

	m.Span = []string{v.Meta.TRangeLower, v.Meta.TRangeUpper}
	err = testDateTimeSpan(m.Span)
	if err != nil {
		return nil, err
	}

	if len(v.Data) == 0 {
		return link.DataV3Sale{}, nil
	}

	m.Nick = v.Data[0].Head.Source
	if v.Data[0].Head.MDSLns != "" {
		m.Nick = m.Nick + ":" + v.Data[0].Head.MDSLns
	}

	d := make(link.DataV3Sale, len(v.Data[0].Item))
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

	return d, nil
}

func unmarshalSaleBy(data []byte) (*link.DataV1SaleBy, error) {
	v := &link.DataV1SaleBy{}

	err := json.Unmarshal(data, &v)
	if err != nil {
		return nil, err
	}

	return v, nil
}

func convSaleBy(data []byte, m *structs.Meta) (link.DataV3SaleBy, error) {
	v, err := unmarshalSaleBy(data)
	if err != nil {
		return nil, err
	}

	m.Span = []string{v.Meta.TRangeLower, v.Meta.TRangeUpper}
	err = testDateTimeSpan(m.Span)
	if err != nil {
		return nil, err
	}

	if len(v.Data) == 0 {
		return link.DataV3SaleBy{}, nil
	}

	m.Nick = v.Data[0].Head.Source + ":" + v.Data[0].Head.Drugstore

	d := make(link.DataV3SaleBy, len(v.Data[0].Item))
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

	return d, nil
}

func unmarshalGeoa(data []byte) (*link.DataV1Geoa, error) {
	v := &link.DataV1Geoa{}

	err := json.Unmarshal(data, &v)
	if err != nil {
		return nil, err
	}

	return v, nil
}

func convGeoa(data []byte, m *structs.Meta) (link.DataV3Geoa, error) {
	v, err := unmarshalGeoa(data)
	if err != nil {
		return nil, err
	}

	m.Name = v.Meta.Name
	m.Head = v.Meta.Head
	m.Addr = v.Meta.Addr
	m.Code = v.Meta.Code
	if v.Meta.EGRPOU != "" {
		m.Code = v.Meta.EGRPOU
	}

	if len(v.Data) == 0 {
		return link.DataV3Geoa{}, nil
	}

	d := make(link.DataV3Geoa, len(v.Data))
	for i, v := range v.Data {
		d[i].ID = v.ID
		if v.Code != "" {
			d[i].ID = v.Code
		}
		d[i].Name = v.Name
		d[i].Home = v.Link
		if v.Addr != "" {
			d[i].Home = v.Addr
		}
		d[i].Quant = v.Quant
		d[i].Price = v.Price
	}

	return d, nil
}

func testDateTimeSpan(s []string) error {
	if len(s) != 2 {
		return fmt.Errorf("api: conv: not enough values inside time span")
	}

	var err error
	for i := range s {
		_, err = time.Parse("02.01.2006 15:04:05", s[i])
		if err != nil {
			_, err = time.Parse("02.01.2006", s[i])
			if err != nil {
				return err
			}
		}
	}

	return nil
}
