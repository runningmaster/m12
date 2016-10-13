package core

import (
	"encoding/json"
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
	trimD = 3 * 24 * time.Hour
)

type path struct {
	Bucket string `json:"bucket,omitempty"`
	Object string `json:"object,omitempty"`
}

// Init inits client for NATS Server
func Init() error {
	sendMessage(BucketStreamOut, SubjectSteamOut, tickD, listN)
	sendMessage(BucketStreamIn, SubjectSteamIn, tickD, listN)
	trimZLog(tickD, trimD)
	return nats.Subscribe(SubjectSteamIn, proc)
}

func sendMessage(b, s string, d time.Duration, n int) {
	_ = time.AfterFunc(d, func() {
		l, err := minio.List(b, n)
		if err != nil {
			log.Println("err: minio:", err)
		} else {
			for i := range l {
				p, err := EncodePath(b, l[i])
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

func EncodePath(b string, o string) ([]byte, error) {
	p := path{b, o}
	return json.Marshal(p)
}

func DecodePath(data []byte) (string, string, error) {
	p := path{}
	err := json.Unmarshal(data, &p)
	return p.Bucket, p.Object, err
}
