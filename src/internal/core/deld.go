package core

var Deld = newDeldWorker()

type deldWorker struct {
}

func newDeldWorker() Worker {
	return &deldWorker{}
}

func (w *deldWorker) NewWorker() Worker {
	return newDeldWorker()
}

func (w *deldWorker) Work(data []byte) (interface{}, error) {
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
