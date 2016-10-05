package proc

import (
	"log"
	"time"

	"internal/core/structs"
	"internal/database/minio"
	"internal/net/nats"
)

// Init inits client for NATS Server
func Init() error {
	sendMessage(structs.BucketStreamOut, structs.SubjectSteamOut, structs.TickD, structs.ListN)
	sendMessage(structs.BucketStreamIn, structs.SubjectSteamIn, structs.TickD, structs.ListN)
	return nats.Subscribe(structs.SubjectSteamIn, proc)
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
