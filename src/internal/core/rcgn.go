package core

func RcgnAddr(meta, data []byte) (interface{}, error) {
	m, err := unmarshalMeta(meta)
	if err != nil {
		return nil, err
	}

	err = testHTag(m.HTag)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func RcgnDrug(meta, data []byte) (interface{}, error) {
	m, err := unmarshalMeta(meta)
	if err != nil {
		return nil, err
	}

	err = testHTag(m.HTag)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
