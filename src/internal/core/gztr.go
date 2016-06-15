package core

import (
	"archive/tar"
	"bytes"
	"io"
	"time"

	"internal/gzpool"
)

const (
	tarMeta = "meta.json.gz"
	tarData = "data.json.gz"
)

func gztarMetaData(m, d []byte) (io.Reader, error) {
	b := new(bytes.Buffer)
	t := tar.NewWriter(b)

	err := writeGzTar(tarMeta, m, t)
	if err != nil {
		return nil, err
	}

	err = writeGzTar(tarData, d, t)
	if err != nil {
		return nil, err
	}

	err = t.Close()
	if err != nil {
		return nil, err
	}

	return b, nil
}

func writeGzTar(name string, data []byte, w *tar.Writer) error {
	d, err := gzpool.MustGzip(data)
	if err != nil {
		return err
	}

	h := &tar.Header{
		Name:    name,
		Mode:    0666,
		ModTime: time.Now(),
		Size:    int64(len(d)),
	}

	err = w.WriteHeader(h)
	if err != nil {
		return err
	}

	_, err = w.Write(d)
	return err
}

func ungztarMetaData(r io.Reader, gz ...bool) ([]byte, []byte, error) {
	tr := tar.NewReader(r)
	var (
		m = new(bytes.Buffer)
		d = new(bytes.Buffer)
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
		case h.Name == tarMeta:
			if len(gz) > 1 && gz[0] == true {
				_, err = io.Copy(m, tr)
			} else {
				err = gzpool.Copy(m, tr)
			}
		case h.Name == tarData:
			if len(gz) == 2 && gz[1] == true {
				_, err = io.Copy(d, tr)
			} else {
				err = gzpool.Copy(d, tr)
			}
		}
		if err != nil {
			return nil, nil, err
		}
	}

	return m.Bytes(), d.Bytes(), nil
}
