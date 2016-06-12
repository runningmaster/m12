package minio

import (
	"fmt"
	"io"
	"net/url"

	minio "github.com/minio/minio-go"
)

var cli *minio.Client

// Run FIXME
func Run(addr string) error {
	u, err := url.Parse(addr)
	if err != nil {
		return err
	}

	var aKey, sKey string
	if u.User != nil {
		aKey = u.User.Username()
		sKey, _ = u.User.Password()
	}

	cli, err = minio.New(u.Host, aKey, sKey, u.Scheme == "https")
	if err != nil {
		return fmt.Errorf("minio: %s", err)
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
				return fmt.Errorf("minio: %s", err)
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

func GetObject(bucket, object string) (io.ReadCloser, error) {
	return cli.GetObject(bucket, object)
}

func DelObject(bucket, object string) error {
	return cli.RemoveObject(bucket, object)
}

// ListObjects FIXME
func ListObjects(backet string, n int) ([]string, error) {
	objs, err := listObjectsN(backet, "", false, n)
	if err != nil {
		return nil, err
	}

	list := make([]string, len(objs))
	for i := range list {
		list[i] = objs[i].Key
	}

	return list, nil
}
