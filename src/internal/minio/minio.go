package minio

import (
	"encoding/json"
	"io"
	"net/url"
	"path/filepath"

	minio "github.com/minio/minio-go"
)

type Clienter interface {
	Put(string, string, io.Reader) error
	Get(string, string) (io.ReadCloser, error)
	Del(string, string) error
	Copy(string, string, string, string) error
	List(string, int) ([]string, error)
	Free(io.Closer)
	Marshal(string, string) ([]byte, error)
	Unmarshal([]byte) (string, string, error)
}

type client struct {
	cli *minio.Client
}

// Init inits client for MINIO Server
func NewClient(addr string) (Clienter, error) {
	c, err := makeConn(addr)
	if err != nil {
		return nil, err
	}

	return &client{c}, nil
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

func (c *client) mustBucket(b string) error {
	ok, err := c.cli.BucketExists(b)
	if err != nil {
		return err
	}

	if !ok {
		return c.cli.MakeBucket(b, "")
	}

	return nil
}

func (c *client) Put(b string, o string, r io.Reader) error {
	err := c.mustBucket(b)
	if err != nil {
		return err
	}

	_, err = c.cli.PutObject(b, o, r, "")
	return err
}

func (c *client) Get(b string, o string) (io.ReadCloser, error) {
	return c.cli.GetObject(b, o)
}

func (c *client) Del(b string, o string) error {
	return c.cli.RemoveObject(b, o)
}

func (c *client) Copy(bDst string, oDst string, bSrc string, oSrc string) error {
	err := c.mustBucket(bDst)
	if err != nil {
		return err
	}

	return c.cli.CopyObject(bSrc, oSrc, filepath.Join(bDst, oDst), minio.NewCopyConditions())
}

func (c *client) List(b string, n int) ([]string, error) {
	doneCh := make(chan struct{}, 1)
	defer func() { close(doneCh) }()

	i := 0
	out := make([]string, 0, n)
	for o := range c.cli.ListObjects(b, "", false, doneCh) {
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

func (c *client) Free(o io.Closer) {
	if o != nil {
		_ = o.Close()
	}
}

type pair struct {
	Bucket string `json:"bucket,omitempty"`
	Object string `json:"object,omitempty"`
}

func (c *client) Marshal(b string, o string) ([]byte, error) {
	p := pair{b, o}
	return json.Marshal(p)
}

func (c *client) Unmarshal(data []byte) (string, string, error) {
	p := pair{}
	err := json.Unmarshal(data, &p)
	return p.Bucket, p.Object, err
}
