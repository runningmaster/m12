package core

import (
	"fmt"
	"internal/flag"

	s3 "github.com/minio/minio-go"
)

const (
	backetStreamIn  = "stream-in"
	backetStreamOut = "stream-out"
	backetStreamErr = "stream-err"
)

var (
	s3cli   *s3.Client
	backets = [...]string{backetStreamIn, backetStreamOut, backetStreamErr}
)

func initS3Cli() error {
	var err error
	s3cli, err = s3.New(flag.S3Address, flag.S3AccessKey, flag.S3SecretKey, true)
	if err != nil {
		return fmt.Errorf("core: s3cli: %s", err)
	}
	testListBackets()
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

func testListBackets() {
	doneCh := make(chan struct{})
	defer close(doneCh)

	n := 0
	objectCh := s3cli.ListObjects(backetStreamIn, "", false, doneCh)
	for object := range objectCh {
		if object.Err != nil {
			fmt.Println(object.Err)
			return
		}
		n++
		fmt.Println(n, object, object.ContentType)
		if n == 5 {
			doneCh <- struct{}{}
		}
		// send to nats
	}
	fmt.Println(n)
}
