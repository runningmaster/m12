package minio

import (
	"io"
	"net/url"

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

	_, err = c.BucketExists("test")
	if err != nil {
		return nil, err
	}

	return c, nil
}

func Put(b, o string, r io.Reader) error {
	//_, err := cli.PutObject(b, o, r, -1, minio.PutObjectOptions{})
	_, err := cli.PutObject(b, o, r, "")
	return err
}

func Get(b, o string) (io.ReadCloser, error) {
	//return cli.GetObject(b, o, minio.GetObjectOptions{})
	return cli.GetObject(b, o)
}

func Del(b, o string) error {
	return cli.RemoveObject(b, o)
}

func Make(b string) error {
	ok, err := cli.BucketExists(b)
	if err != nil {
		return err
	}

	if !ok {
		return cli.MakeBucket(b, "")
	}

	return nil
}

func Copy(bDst, oDst, bSrc, oSrc string) error {
	dst, err := minio.NewDestinationInfo(bDst, oDst, nil, nil)
	if err != nil {
		return err
	}
	src := minio.NewSourceInfo(bSrc, oSrc, nil)
	return cli.CopyObject(dst, src)
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
