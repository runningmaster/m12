package core

import (
	"archive/tar"
	"bytes"
	"io"
	"time"
)

const (
	tarMeta = "meta.json"
	tarData = "data.json.gz"
)

func tarMetaData(m, d []byte) (io.Reader, error) {
	b := new(bytes.Buffer)
	t := tar.NewWriter(b)

	err := writeToTar(tarMeta, m, t)
	if err != nil {
		return nil, err
	}

	err = writeToTar(tarData, d, t)
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
		Name:    name,
		Mode:    0644,
		ModTime: time.Now(),
		Size:    int64(len(data)),
	}

	err := w.WriteHeader(h)
	if err != nil {
		return err
	}

	_, err = w.Write(data)
	return err
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
		case h.Name == tarMeta:
			if _, err := io.Copy(meta, tr); err != nil {
				return nil, nil, err
			}
		case h.Name == tarData:
			if _, err := io.Copy(data, tr); err != nil {
				return nil, nil, err
			}
		}
	}

	return meta.Bytes(), data.Bytes(), nil
}
