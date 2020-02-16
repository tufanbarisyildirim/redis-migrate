package migrator

import (
	"github.com/go-redis/redis"
	"log"
)

type Exporter struct {
	Source *redis.Client
	Chan   chan Cursor
	Done   chan struct{}
}

func NewExporter(client *redis.Client, buf int) *Exporter {
	exp := &Exporter{
		Source: client,
		Chan:   make(chan Cursor, buf),
		Done:   make(chan struct{}),
	}
	go exp.Read()
	return exp
}

func (e *Exporter) Read() {
	crs := uint64(0)
	for {
		keys, cursor, err := e.Source.Scan(crs, "*", 1).Result()
		log.Printf("exported cursor : %d", cursor)
		if err != nil {
			log.Fatalf("scanning failed: %v", err)
		}
		e.Chan <- Cursor{Keys: keys}
		if cursor == 0 {
			break
		}
		crs = cursor
	}
	close(e.Done)
}
