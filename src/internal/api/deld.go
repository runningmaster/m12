package api

func deld(data []byte) (interface{}, error) {
	bucket, object, err := unmarshaPairExt(data)
	if err != nil {
		return nil, err
	}

	err = cMINIO.RemoveObject(bucket, object)
	if err != nil {
		return nil, err
	}

	return "OK", nil
}
