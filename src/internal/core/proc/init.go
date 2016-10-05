package proc

import (
	"log"
	"time"

	"internal/database/minio"
	"internal/net/nats"
)

const (
	bucketStreamIn  = "stream-in"
	bucketStreamOut = "stream-out"
	bucketStreamErr = "stream-err"

	subjectSteamIn  = "m12." + bucketStreamIn
	subjectSteamOut = "m12." + bucketStreamOut

	// should be move to pref
	listN = 100
	tickD = 10 * time.Second
)

// Init inits client for NATS Server
func Init() error {
	sendMessage(bucketStreamOut, subjectSteamOut, tickD, listN)
	sendMessage(bucketStreamIn, subjectSteamIn, tickD, listN)
	return nats.Subscribe(subjectSteamIn, proc)
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
