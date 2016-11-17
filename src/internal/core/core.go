package core

import (
	"encoding/json"
	"log"
	"time"

	"internal/database/minio"
	"internal/net/nats"
)

const (
	bucketStreamIn     = "stream-in"
	bucketStreamErr    = "stream-err"
	bucketStreamOut    = "stream-out"
	bucketStreamOutGeo = "stream-out.geo"

	subjectSteamIn     = "m12." + bucketStreamIn
	subjectSteamOut    = "m12." + bucketStreamOut
	subjectSteamOutGeo = "m12." + bucketStreamOutGeo

	// should be move to pref
	listN = 50
	tickD = 10 * time.Second
	trimD = 3 * 24 * time.Hour
)

type path struct {
	Bucket string `json:"bucket,omitempty"`
	Object string `json:"object,omitempty"`
}

// Init inits package
func Init() error {
	initBuckets(bucketStreamIn, bucketStreamErr, bucketStreamOut, bucketStreamOutGeo)
	sendMessage(bucketStreamOut, subjectSteamOut, tickD, listN)
	sendMessage(bucketStreamIn, subjectSteamIn, tickD, listN)
	sendMessage(bucketStreamOutGeo, subjectSteamOutGeo, tickD, listN)
	trimZLog(tickD*60, trimD)
	return nats.Subscribe(subjectSteamIn, proc)
}

func initBuckets(b ...string) error {
	var err error
	for i := range b {
		err = minio.Make(b[i])
		if err != nil {
			return err
		}
	}
	return nil
}

func sendMessage(b, s string, d time.Duration, n int) {
	_ = time.AfterFunc(d, func() {
		l, err := minio.List(b, n)
		if err != nil {
			log.Println(err)
		} else {
			for i := range l {
				p, err := encodePath(b, l[i])
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

func trimZLog(d, t time.Duration) {
	_ = time.AfterFunc(d, func() {
		err := remZlog(time.Now().Add(-1 * t).Unix())
		if err != nil {
			log.Println(err)
		}
		trimZLog(d, t)
	})

}

func encodePath(b string, o string) ([]byte, error) {
	p := path{b, o}
	return json.Marshal(p)
}

func decodePath(data []byte) (path, error) {
	p := path{}
	err := json.Unmarshal(data, &p)
	return p, err
}
