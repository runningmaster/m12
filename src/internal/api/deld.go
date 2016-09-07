package api

func deld(data []byte) (interface{}, error) {
	b, o, err := cMINIO.Unmarshal(data)
	if err != nil {
		return nil, err
	}

	err = cMINIO.Del(b, o)
	if err != nil {
		return nil, err
	}

	return "OK", nil
}
