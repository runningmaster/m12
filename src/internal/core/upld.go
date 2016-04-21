package core

import (
	"archive/tar"
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"internal/context/ctxutil"

	"golang.org/x/net/context"
)

// Upld sends data to s3 interface
func Upld(ctx context.Context, r *http.Request) (interface{}, error) {
	d, err := readClose(r.Body)
	if err != nil {
		return nil, err
	}

	if !isTypeGzip(d) {
		return nil, fmt.Errorf("core: content must contain gzip")
	}

	m := meta{}
	err = m.initFromJSON([]byte(ctxutil.MetaFromContext(ctx)))
	if err != nil {
		return nil, err
	}

	m.ID = ctxutil.IDFromContext(ctx)
	m.IP = ctxutil.IPFromContext(ctx)
	m.Auth = ctxutil.AuthFromContext(ctx)
	m.HTag = strings.ToLower(m.HTag)
	m.Time = time.Now().Format("02.01.2006 15:04:05.999999999")
	m.ETag = btsToMD5(d)
	m.Path = "" // ?
	m.Size = int64(len(d))
	m.SrcCE = r.Header.Get("Content-Encoding")
	m.SrcCT = r.Header.Get("Content-Type")

	p, err := m.packToJSON()
	if err != nil {
		return nil, err
	}

	t, err := tarMetaData(p, d)
	if err != nil {
		return nil, err
	}

	err = putObject(backetStreamIn, m.ID+".tar", t)
	if err != nil {
		return nil, err
	}

	o, err := popObject(backetStreamIn, m.ID+".tar")
	if err != nil {
		return nil, err
	}

	meta, data, err := untarMetaData(o)
	if err != nil {
		return nil, err
	}

	fmt.Println("meta:", string(meta), isTypeGzip(data))

	// send ?

	return m.ID, nil
}

func tarMetaData(meta, data []byte) (io.Reader, error) {
	bf := new(bytes.Buffer)
	tw := tar.NewWriter(bf)

	hm := &tar.Header{
		Name: "meta",
		Size: int64(len(meta)),
	}
	err := tw.WriteHeader(hm)
	if err != nil {
		return nil, err
	}
	_, err = tw.Write(meta)
	if err != nil {
		return nil, err
	}

	hd := &tar.Header{
		Name: "data",
		Size: int64(len(data)),
	}
	err = tw.WriteHeader(hd)
	if err != nil {
		return nil, err
	}
	_, err = tw.Write(data)
	if err != nil {
		return nil, err
	}

	err = tw.Close()
	if err != nil {
		log.Fatalln(err)
	}

	return bf, nil
}

func untarMetaData(rc io.ReadCloser) ([]byte, []byte, error) {
	tr := tar.NewReader(rc)
	defer rc.Close()

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
