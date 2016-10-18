package core

import (
	"encoding/json"
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
	trimD = 3 * 24 * time.Hour
)

type path struct {
	Bucket string `json:"bucket,omitempty"`
	Object string `json:"object,omitempty"`
}

// Init inits client for NATS Server
func Init() error {
	sendMessage(bucketStreamOut, subjectSteamOut, tickD, listN)
	sendMessage(bucketStreamIn, subjectSteamIn, tickD, listN)
	trimZLog(tickD, trimD)
	return nats.Subscribe(subjectSteamIn, proc)
}

func sendMessage(b, s string, d time.Duration, n int) {
	_ = time.AfterFunc(d, func() {
		l, err := minio.List(b, n)
		if err != nil {
			log.Println("err: minio:", err)
		} else {
			for i := range l {
				p, err := encodePath(b, l[i])
				if err != nil {
					log.Println("err: minio:", err)
				}
				err = nats.Publish(s, p)
				if err != nil {
					log.Println("err: nats:", err)
				}
			}
		}
		sendMessage(b, s, d, n)
	})
}

func trimZLog(d, t time.Duration) {
	_ = time.AfterFunc(d, func() {
		err := remZlog(time.Now().Add(-1 * t).Unix())
		if err != nil {
			log.Println("err: redis:", err)
		}
		trimZLog(d, t)
	})

}

func encodePath(b string, o string) ([]byte, error) {
	p := path{b, o}
	return json.Marshal(p)
}

func decodePath(data []byte) (string, string, error) {
	p := path{}
	err := json.Unmarshal(data, &p)
	return p.Bucket, p.Object, err
}
