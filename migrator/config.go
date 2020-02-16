package migrator

import (
	"fmt"
	"github.com/go-redis/redis"
)

type Config struct {
	Servers map[string]Server
	Jobs    map[string]Job
}

type Server struct {
	Name     string
	Addr     string
	Password string
	Db       int
}

type Job struct {
	From string
	To   string
}

func (c *Config) CreateClient(serverName string) *redis.Client {
	if server, ok := c.Servers[serverName]; ok {
		return redis.NewClient(&redis.Options{
			Addr:     server.Addr,
			Password: server.Password,
			DB:       server.Db,
		})
	}
	panic(fmt.Errorf("server not found in config toml %s", serverName))
}
