package core

import (
	"fmt"
	"internal/flag"
	"io"
	"log"

	minio "github.com/minio/minio-go"
)

const (
	backetStreamIn  = "stream-in"
	backetStreamOut = "stream-out"
	backetStreamErr = "stream-err"
)

var (
	s3cli   *minio.Client
	backets = [...]string{backetStreamIn, backetStreamOut, backetStreamErr}
)

func initS3Cli() error {
	var err error
	s3cli, err = minio.New(flag.S3Address, flag.S3AccessKey, flag.S3SecretKey, true)
	if err != nil {
		return fmt.Errorf("core: s3cli: %s", err)
	}

	return initBackets()
}

func initBackets() error {
	for i := range backets {
		b := backets[i]
		err := s3cli.BucketExists(b)
		if err != nil {
			err = s3cli.MakeBucket(b, "")
			if err != nil {
				return fmt.Errorf("core: s3cli: %s", err)
			}
		}
	}

	return nil
}

func putObject(bucket, object string, r io.Reader) error {
	_, err := s3cli.PutObject(bucket, object, r, "")
	return err
}

func popObject(bucket, object string) (io.ReadCloser, error) {
	o, err := s3cli.GetObject(bucket, object)
	if err != nil {
		return nil, err
	}

	err = s3cli.RemoveObject(bucket, object)
	if err != nil {
		return nil, err
	}

	return o, nil
}

func listObjectsN(bucket, prefix string, recursive bool, n int) ([]minio.ObjectInfo, error) {
	doneCh := make(chan struct{}, 1)
	defer close(doneCh)

	i := 0
	objs := make([]minio.ObjectInfo, 0, n)
	for object := range s3cli.ListObjects(bucket, prefix, recursive, doneCh) {
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

func goToStreamIn(object string, r io.Reader) {
	goToStream(backetStreamIn, object, r)
}

func goToStreamOut(object string, r io.Reader) {
	goToStream(backetStreamOut, object, r)
}

func goToStreamErr(object string, r io.Reader) {
	goToStream(backetStreamErr, object, r)
}

func goToStream(bucket, object string, r io.Reader) {
	go func() {
		err := putObject(bucket, object, r)
		if err != nil {
			log.Println(err)
		}
	}()
}
