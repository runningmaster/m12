package s3

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/url"

	minio "github.com/minio/minio-go"
)

var (
	cli    *minio.Client
	logger = log.New(ioutil.Discard, "", log.LstdFlags)
)

// Run FIXME
func Run(addr, akey, skey string, l *log.Logger) error {
	if l != nil {
		logger = l
	}

	u, err := url.Parse(addr)
	if err != nil {
		return err
	}

	cli, err = minio.New(u.Host, akey, skey, true)
	if err != nil {
		return fmt.Errorf("s3: %s", err)
	}

	return nil
}

// InitBacketList FIXME
func InitBacketList(list ...string) error {
	for i := range list {
		b := list[i]
		err := cli.BucketExists(b)
		if err != nil {
			err = cli.MakeBucket(b, "")
			if err != nil {
				return fmt.Errorf("s3: %s", err)
			}
		}
	}

	return nil
}

func listObjectsN(bucket, prefix string, recursive bool, n int) ([]minio.ObjectInfo, error) {
	doneCh := make(chan struct{}, 1)
	defer func() { close(doneCh) }()

	i := 0
	objs := make([]minio.ObjectInfo, 0, n)
	for object := range cli.ListObjects(bucket, prefix, recursive, doneCh) {
		if object.Err != nil {
			return nil, object.Err
		}
		objs = append(objs, object)
		i++
		if i == n {
			doneCh <- struct{}{}
		}
	}

	return objs, nil
}

// PutObject FIXME
func PutObject(bucket, object string, r io.Reader) error {
	_, err := cli.PutObject(bucket, object, r, "")
	return err
}

// PopObject FIXME
func PopObject(bucket, object string) (io.Reader, error) {
	obj, err := cli.GetObject(bucket, object)
	if err != nil {
		return nil, err
	}

	err = cli.RemoveObject(bucket, object)
	if err != nil {
		return nil, err
	}

	return obj, nil
}

// ListObjects FIXME
func ListObjects(backet string, n int) ([]string, error) {
	objs, err := listObjectsN(backet, "", false, n)
	if err != nil {
		return nil, err
	}

	list := make([]string, n)
	for i := range list {
		list[i] = objs[i].Key
	}

	return list, nil
}
