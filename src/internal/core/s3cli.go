package core

import (
	"internal/flag"

	s3 "github.com/minio/minio-go"
)

const (
	backetStreamIn  = "stream-in"
	backetStreamOut = "stream-out"
	backetStreamErr = "stream-out"
)

var (
	s3cli   *s3.Client
	backets = [...]string{backetStreamIn, backetStreamOut, backetStreamErr}
)

func initCliS3() error {
	var err error
	if s3cli, err = s3.New(flag.S3Address, flag.S3AccessKey, flag.S3SecretKey, true); err != nil {
		return err
	}
	return initBackets()
}

func initBackets() error {
	l, err := s3cli.ListBuckets()
	if err != nil {
		return err
	}

	m := make(map[string]struct{}, len(l))
	for i := range l {
		m[l[i].Name] = struct{}{}
	}

	for i := range backets {
		b := backets[i]
		if _, ok := m[b]; !ok {
			if err = s3cli.MakeBucket(b, ""); err != nil {
				return err
			}
		}
	}

	return nil
}
