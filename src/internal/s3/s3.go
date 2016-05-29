package s3

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"

	minio "github.com/minio/minio-go"
)

var (
	cli    *minio.Client
	logger = log.New(ioutil.Discard, "", log.LstdFlags)
)

// Run FIXME
func Run(addr, akey, skey string, log *log.Logger) error {
	if log != nil {
		logger = log
	}

	var err error
	cli, err = minio.New(addr, akey, skey, true)
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

type pair struct {
	Backet string `json:"backet,omitempty"`
	Object string `json:"object,omitempty"`
}

func (p pair) marshalJSON() []byte {
	b, _ := json.Marshal(p)
	return b
}

func unmarshaJSON(data []byte) (pair, error) {
	p := pair{}
	err := json.Unmarshal(data, &p)
	return p, err
}

// GetPathJSONList FIXME
func ListObjectsMarshal(backet string, n int) ([][]byte, error) {
	objs, err := listObjectsN(backet, "", false, n)
	if err != nil {
		return nil, err
	}

	list := make([][]byte, n)
	for i := range list {
		list[i] = pair{backet, objs[i].Key}.marshalJSON()
	}

	return list, nil
}

// PopObject FIXME
func PopObjectUnmarshal(data []byte) (io.Reader, error) {
	p, err := unmarshaJSON(data)
	if err != nil {
		return nil, err
	}

	return PopObject(p.Backet, p.Object)
}
