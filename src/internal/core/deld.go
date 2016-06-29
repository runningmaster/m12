package core

var Deld = &deld{}

type deld struct {
}

func (w *deld) New() interface{} {
	return &deld{}
}

func (w *deld) Work(data []byte) (interface{}, error) {
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
