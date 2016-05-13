package core

import (
	"archive/tar"
	"bytes"
	"io"
	"net/http"
	"time"

	"golang.org/x/net/context"
)

// Upld puts data to s3 interface
func Upld(ctx context.Context, _ http.ResponseWriter, r *http.Request) (interface{}, error) {
	v, _ := ctx.Value("ctxMeta").(string)

	m, err := makeMetaFromBase64String("FIXME ctxutil.MetaFromContext(ctx)")
	if err != nil {
		return nil, err
	}

	//m.ID = ctxutil.IDFromContext(ctx)
	//m.IP = ctxutil.IPFromContext(ctx)
	//m.Auth = ctxutil.AuthFromContext(ctx)
	m.Time = time.Now().Unix()
	m.SrcCE = r.Header.Get("Content-Encoding")
	m.SrcCT = r.Header.Get("Content-Type")

	t, err := tarMetaData(makeReadCloser(m.packToJSON()), r.Body) // FIXME base64?
	if err != nil {
		return nil, err
	}

	goToStreamIn(m.ID, t)

	return m.ID, nil
}

func tarMetaData(m, d io.ReadCloser) (io.Reader, error) {
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

func writeToTar(name string, rc io.ReadCloser, w *tar.Writer) error {
	b, err := readClose(rc)
	if err != nil {
		return err
	}

	h := &tar.Header{
		Name: name,
		Size: int64(len(b)),
	}

	err = w.WriteHeader(h)
	if err != nil {
		return err
	}

	_, err = w.Write(b)
	if err != nil {
		return err
	}

	return nil
}

func untarMetaData(rc io.ReadCloser) ([]byte, []byte, error) {
	defer func() { _ = rc.Close() }()

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
