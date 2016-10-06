package core

import (
	"log"
	"time"

	"internal/database/minio"
	"internal/net/nats"
)

const (
	BucketStreamIn  = "stream-in"
	BucketStreamOut = "stream-out"
	BucketStreamErr = "stream-err"

	SubjectSteamIn  = "m12." + BucketStreamIn
	SubjectSteamOut = "m12." + BucketStreamOut

	// should be move to pref
	listN = 100
	tickD = 10 * time.Second
)

// Init inits client for NATS Server
func Init() error {
	sendMessage(BucketStreamOut, SubjectSteamOut, tickD, listN)
	sendMessage(BucketStreamIn, SubjectSteamIn, tickD, listN)
	return nats.Subscribe(SubjectSteamIn, proc)
}

func sendMessage(b, s string, d time.Duration, n int) {
	_ = time.AfterFunc(d, func() {
		l, err := minio.List(b, n)
		if err != nil {
			log.Println(err)
		} else {
			for i := range l {
				p, err := minio.Pair(b, l[i])
				if err != nil {
					log.Println(err)
				}
				err = nats.Publish(s, p)
				if err != nil {
					log.Println(err)
				}
			}
		}
		sendMessage(b, s, d, n)
	})
}
