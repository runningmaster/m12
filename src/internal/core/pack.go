package core

import (
	"archive/tar"
	"bytes"
	"io"
	"time"

	"internal/compress/gziputil"
)

const (
	tarMeta = "meta.json.gz"
	tarData = "data.json.gz"
)

func packMetaData(m, d []byte) (io.Reader, error) {
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
	d, err := gziputil.MustCompress(data)
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

func UnpackMetaData(r io.Reader, gz ...bool) ([]byte, []byte, error) {
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
			err = copyMetaData(m, tr, len(gz) > 1 && gz[0])
		case h.Name == tarData:
			err = copyMetaData(d, tr, len(gz) == 2 && gz[1])
		}
		if err != nil {
			return nil, nil, err
		}
	}

	return m.Bytes(), d.Bytes(), nil
}

func copyMetaData(dst io.Writer, src io.Reader, gz bool) error {
	if gz {
		_, err := io.Copy(dst, src)
		return err

	}
	return gziputil.Copy(dst, src)
}
