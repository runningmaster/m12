package api

func deld(data []byte) (interface{}, error) {
	b, o, err := minio.Unmarshal(data)
	if err != nil {
		return nil, err
	}

	err = minio.Del(b, o)
	if err != nil {
		return nil, err
	}

	return "OK", nil
}
