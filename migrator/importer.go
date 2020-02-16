package migrator

import (
	"github.com/go-redis/redis"
	"log"
	"sync"
)

type Importer struct {
	Destination *redis.Client
	WorkerCount int
	Wg          *sync.WaitGroup
}

func NewImporter(client *redis.Client, workerCount int) *Importer {
	return &Importer{
		Destination: client,
		WorkerCount: workerCount,
		Wg:          &sync.WaitGroup{},
	}
}

func (importer *Importer) WriteFrom(exporter *Exporter) {
	for i := 0; i < importer.WorkerCount; i++ {
		importer.Wg.Add(1)
		go func(workerNum int) {
		listener:
			for {
				select {
				case c := <-exporter.Chan:
					for _, key := range c.Keys {
						value, err := exporter.Source.Dump(key).Result()
						if err != nil {
							log.Println(err)
							continue
						}
						ttl, err := exporter.Source.TTL(key).Result()
						if err != nil {
							log.Println(err) //we can ignore ttl errors here
						}
						if ttl < 0 {
							ttl = 0
						}
						log.Printf("[%d] importing key : %s\n", workerNum, key)
						_, err = importer.Destination.Restore(key, ttl, value).Result()
						if err != nil {
							log.Printf("error restoring %s : %s", key, err)
						}
					}

				case <-exporter.Done:
					break listener
				}

			}

			importer.Wg.Done()
		}(i)
	}
}
