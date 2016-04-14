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

func initCliS3() error {
	var err error
	if s3cli, err = s3.New(flag.S3Address, flag.S3AccessKey, flag.S3SecretKey, true); err != nil {
		return fmt.Errorf("core: s3cli: %s", err)
	}

	return initBackets()
}

func initBackets() error {
	for i := range backets {
		b := backets[i]
		if err := s3cli.BucketExists(b); err != nil {
			if err = s3cli.MakeBucket(b, ""); err != nil {
				return fmt.Errorf("core: s3cli: %s", err)
			}
		}
	}

	return nil
}
