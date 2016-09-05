package minio

import (
	"net/url"

	minio "github.com/minio/minio-go"
)

// Init return active connection to MINIO Server
func Init(addr string) (*minio.Client, error) {
	return makeConn(addr)
	// return makeBuckets(bucketStreamIn, bucketStreamOut, bucketStreamErr)
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

func makeBuckets(c *minio.Client, list ...string) error {
	for i := range list {
		b := list[i]
		ok, err := c.BucketExists(b)
		if err != nil {
			return err
		}
		if !ok {
			err = c.MakeBucket(b, "")
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func listObjectsN(c *minio.Client, bucket string, n int) ([]string, error) {
	doneCh := make(chan struct{}, 1)
	defer func() { close(doneCh) }()

	i := 0
	out := make([]string, 0, n)
	for o := range c.ListObjects(bucket, "", false, doneCh) {
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
