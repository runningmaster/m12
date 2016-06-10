package core

import (
	"archive/tar"
	"bytes"
	"crypto/md5"
	"crypto/sha1"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"internal/gzutil"
	"internal/minio"

	"github.com/spkg/bom"
)

func btsToMD5(b []byte) string {
	return fmt.Sprintf("%x", md5.Sum(b))
}

func strToMD5(s string) string {
	return btsToMD5([]byte(s))
}

func btsToSHA1(b []byte) string {
	return fmt.Sprintf("%x", sha1.Sum(b))
}

func strToSHA1(s string) string {
	return btsToSHA1([]byte(s))
}

func readClose(r io.ReadCloser) ([]byte, error) {
	defer func() { _ = r.Close() }()
	return ioutil.ReadAll(r)
}

func isTypeGzip(b []byte) bool {
	return gzutil.IsGzipInString(http.DetectContentType(b))
}

func isTypeUTF8(b []byte) bool {
	return strings.Contains(http.DetectContentType(b), "text/plain; charset=utf-8")
}

func mendIfGzip(b []byte) ([]byte, error) {
	if isTypeGzip(b) {
		return gzutil.Gunzip(b)
	}

	return b, nil
}

func mendIfUTF8(b []byte) ([]byte, error) {
	if isTypeUTF8(b) {
		return bom.Clean(b), nil
	}

	return b, nil
}

func isEmpty(v []interface{}) bool {
	for i := range v {
		if v[i] != nil {
			return false
		}
	}
	return true
}

func makeReadCloser(b []byte) io.ReadCloser {
	return ioutil.NopCloser(bytes.NewReader(b))
}

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
	return err
}

func untarMetaData(r io.Reader) ([]byte, []byte, error) {
	tr := tar.NewReader(r)
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

func popMetaData(data []byte) ([]byte, []byte, error) {
	p, err := unmarshaJSONpair(data)
	if err != nil {
		return nil, nil, err
	}

	o, err := minio.GetObject(p.Backet, p.Object)
	if err != nil {
		return nil, nil, err
	}
	defer func() {
		if o == nil {
			return
		}
		err := o.Close()
		if err != nil {
			log.Println("popMetaData", err)
			return
		}
		//err = minio.DelObject(p.Backet, p.Object)
		//if err != nil {
		//	log.Println("popMetaData", err)
		//}
	}()

	return untarMetaData(o)
}
