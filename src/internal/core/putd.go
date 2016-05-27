package core

import (
	"archive/tar"
	"bytes"
	"io"
	"net/http"
)

var Putd = &putd{}

type putd struct {
	test string
}

func (u *putd) ReadHeader(h http.Header) {
	_ = h.Get("Content-Meta")
	u.test = h.Get("Content-Test")
}

func (u *putd) WriteHeader(h http.Header) {
	h.Set("Content-Test", u.test+" DEBUG")
}

func (u *putd) Work([]byte) (interface{}, error) {
	return nil, nil
}

/*
// Upld puts data to s3 interface
func Upld(ctx context.Context, _ http.ResponseWriter, r *http.Request) (interface{}, error) {
	m, err := makeMetaFromBase64String("FIXME ctxutil.MetaFromContext(ctx)")
	if err != nil {
		return nil, err
	}


	_, err = tarMetaData(makeReadCloser(m.packToJSON()), r.Body) // FIXME base64?
	if err != nil {
		return nil, err
	}

	//goToStreamIn(m.ID, t)

	return m.ID, nil
}
*/
func tarMetaData(m, d []byte) (io.Reader, error) {
	b := new(bytes.Buffer)
	t := tar.NewWriter(b)

	err := writeToTar("meta", m, t)
	if err != nil {
		return nil, err
	}

	err = writeToTar("data", d, t)
	if err != nil {
		return nil, err
	}

	err = t.Close()
	if err != nil {
		return nil, err
	}

	return b, nil
}

func writeToTar(name string, data []byte, w *tar.Writer) error {
	h := &tar.Header{
		Name: name,
		Size: int64(len(data)),
	}

	err := w.WriteHeader(h)
	if err != nil {
		return err
	}

	_, err = w.Write(data)
	if err != nil {
		return err
	}

	return nil
}

func untarMetaData(rc io.Reader) ([]byte, []byte, error) {
	//defer func() { _ = rc.Close() }()

	tr := tar.NewReader(rc)
	var (
		meta = new(bytes.Buffer)
		data = new(bytes.Buffer)
	)
	for {
		h, err := tr.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, nil, err
		}
		switch {
		case h.Name == "meta":
			if _, err := io.Copy(meta, tr); err != nil {
				return nil, nil, err
			}
		case h.Name == "data":
			if _, err := io.Copy(data, tr); err != nil {
				return nil, nil, err
			}
		}
	}

	return meta.Bytes(), data.Bytes(), nil
}
