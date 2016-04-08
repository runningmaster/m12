package s3

import (
	"fmt"
	"io"
	"sync"

	"internal/flag"

	s3 "github.com/minio/minio-go"
)

var (
	pool = sync.Pool{
		New: func() interface{} {
			c, err := s3.New(flag.S3Address, flag.S3AccessKey, flag.S3SecretKey, true)
			if err != nil {
				return err
			}
			return c
		},
	}
)

func getCli() (*s3.Client, error) {
	switch c := pool.Get().(type) {
	case *s3.Client:
		return c, nil
	case error:
		return nil, c
	}
	return nil, fmt.Errorf("s3: unreachable")
}

func putCli(x interface{}) {
	pool.Put(x)
}

func MkB(bucketName string) error {
	c, err := getCli()
	if err != nil {
		return err
	}
	defer putCli(c)

	return c.MakeBucket(bucketName, "")
}

func RmB(bucketName string) error {
	c, err := getCli()
	if err != nil {
		return err
	}
	defer putCli(c)

	return c.RemoveBucket(bucketName)
}

func Put(bucketName, objectName string, r io.Reader, contentType string) error {
	c, err := getCli()
	if err != nil {
		return err
	}
	defer putCli(c)

	if _, err := c.PutObject(bucketName, objectName, r, contentType); err != nil {
		return err
	}

	return nil
}

func Get(bucketName, objectName string) (io.Reader, error) {
	c, err := getCli()
	if err != nil {
		return nil, err
	}
	defer putCli(c)

	return c.GetObject(bucketName, objectName)
}
