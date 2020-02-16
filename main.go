package main

import (
	"flag"
	"github.com/go-redis/redis"
	"github.com/tufanbarisyildirim/redis-migrate/migrator"
	"log"
)

func main() {

	var (
		//configFile   = flag.String("config", "0", "config.toml")
		//	migrationJob = flag.String("job", "default", "the migration job name to work")
		workerCount = flag.Int("worker", 5, "concurrent worker size")
	)

	flag.Usage = func() {
		log.Printf("Usage of %s :\n", "redis-migrate")
		flag.PrintDefaults()
	}

	flag.Parse()

	source := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       2,
	})

	destination := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6388",
		Password: "",
		DB:       0,
	})

	_, err := source.Ping().Result()
	if err != nil {
		log.Fatalf("error connecting source redis: %s", err)
	}

	_, err = destination.Ping().Result()
	if err != nil {
		log.Fatalf("error connecting destination redis: %s", err)
	}

	exporter := migrator.NewExporter(source, *workerCount)
	importer := migrator.NewImporter(destination, *workerCount)
	importer.WriteFrom(exporter)

	importer.Wg.Wait()
}
