package minio

import (
	"encoding/json"
	"io"
	"net/url"
	"path/filepath"

	minio "github.com/minio/minio-go"
)

var cli *minio.Client

// Init inits client for MINIO Server
func Init(addr string) error {
	c, err := makeConn(addr)
	if err != nil {
		return err
	}
	cli = c
	return nil
}

func makeConn(addr string) (*minio.Client, error) {
	u, err := url.Parse(addr)
	if err != nil {
		return nil, err
	}

	var aKey, sKey string
	if u.User != nil {
		aKey = u.User.Username()
		sKey, _ = u.User.Password()
	}

	c, err := minio.New(u.Host, aKey, sKey, u.Scheme == "https")
	if err != nil {
		return nil, err
	}

	return c, nil
}

func mustBucket(b string) error {
	ok, err := cli.BucketExists(b)
	if err != nil {
		return err
	}

	if !ok {
		return cli.MakeBucket(b, "")
	}

	return nil
}

func Put(b string, o string, r io.Reader) error {
	err := mustBucket(b)
	if err != nil {
		return err
	}

	_, err = cli.PutObject(b, o, r, "")
	return err
}

func Get(b string, o string) (io.ReadCloser, error) {
	return cli.GetObject(b, o)
}

func Del(b string, o string) error {
	return cli.RemoveObject(b, o)
}

func Copy(bDst string, oDst string, bSrc string, oSrc string) error {
	err := mustBucket(bDst)
	if err != nil {
		return err
	}

	return cli.CopyObject(bSrc, oSrc, filepath.Join(bDst, oDst), minio.NewCopyConditions())
}

func List(b string, n int) ([]string, error) {
	doneCh := make(chan struct{}, 1)
	defer func() { close(doneCh) }()

	i := 0
	out := make([]string, 0, n)
	for o := range cli.ListObjects(b, "", false, doneCh) {
		if o.Err != nil {
			return nil, o.Err
		}
		if len(out) < n { // workaround
			out = append(out, o.Key)
		}
		i++
		if i == n {
			doneCh <- struct{}{}
		}
	}

	return out, nil
}

func Free(o io.Closer) {
	if o != nil {
		_ = o.Close()
	}
}

type pair struct {
	Bucket string `json:"bucket,omitempty"`
	Object string `json:"object,omitempty"`
}

func Pair(b string, o string) ([]byte, error) {
	p := pair{b, o}
	return json.Marshal(p)
}

func Unmarshal(data []byte) (string, string, error) {
	p := pair{}
	err := json.Unmarshal(data, &p)
	return p.Bucket, p.Object, err
}
