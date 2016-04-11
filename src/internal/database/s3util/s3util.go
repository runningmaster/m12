package s3util

import (
	"fmt"
	"io"
	"log"
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

func init() {
	_, err := LsB()
	if err != nil {
		log.Fatal(err)
	}
}

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

// LsB is wrapper for ListBuckets()
func LsB() ([]string, error) {
	c, err := getCli()
	if err != nil {
		return nil, err
	}
	defer putCli(c)

	l, err := c.ListBuckets()
	if err != nil {
		return nil, err
	}

	s := make([]string, len(l))
	for i := range l {
		s[i] = l[i].Name
	}

	return s, nil
}

// MkB is wrapper for MakeBucket()
func MkB(bucketName string) error {
	c, err := getCli()
	if err != nil {
		return err
	}
	defer putCli(c)

	return c.MakeBucket(bucketName, "")
}

// RmB is wrapper for RemoveBucket()
func RmB(bucketName string) error {
	c, err := getCli()
	if err != nil {
		return err
	}
	defer putCli(c)

	return c.RemoveBucket(bucketName)
}

// Put is wrapper for PutObject()
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

// Get is wrapper for GetObject()
func Get(bucketName, objectName string) (io.Reader, error) {
	c, err := getCli()
	if err != nil {
		return nil, err
	}
	defer putCli(c)

	return c.GetObject(bucketName, objectName)
}
