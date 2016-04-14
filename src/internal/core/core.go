package core

import (
	"net/http"

	"internal/database/s3util"

	"golang.org/x/net/context"
)

const (
	backetStreamIn  = "stream-in"
	backetStreamOut = "stream-out"
	backetStreamErr = "stream-out"
)

var backets = [...]string{backetStreamIn, backetStreamOut, backetStreamErr}

// Handler is func for processing data from api.
type Handler func(context.Context, *http.Request) (interface{}, error)

func Init() error {
	return nil
}

func initBackets() error {
	l, err := s3util.LsB()
	if err != nil {
		return err
	}

	m := make(map[string]struct{}, len(l))
	for i := range l {
		m[l[i]] = struct{}{}
	}

	for i := range backets {
		b := backets[i]
		if _, ok := m[b]; !ok {
			if err = s3util.MkB(b); err != nil {
				return err
			}
		}
	}

	return nil
}

//
//	if err != nil {
//		return nil, err
//	}

// s3util.Get and check content type
