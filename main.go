package main

import (
	"flag"
	"time"

	"github.com/go-redis/redis"
	"github.com/golang/glog"
)

// NewRedisClient initialises and return the Redis client
func NewRedisClient(url string, password string) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     url,
		Password: password,
		DB:       0,
	})

	_, err := client.Ping().Result()
	if err != nil {
		glog.Errorf("NewRedisClient: cannot initialise the Redis client due to error %s", err)
		return nil, err
	}
	return client, nil
}

func main() {
	flag.Parse()

	r, err := NewRedisClient("localhost:6379", "password")
	if err != nil {
		glog.Errorf("main: Redis client can not be initialised due to error %s", err)
		panic(err)
	}
	defer r.Close()

	keys := []string{"test@1", "test@2", "test@3", "test@1", "test@3", "test@2", "test@4"}
	cmds, err := r.Pipelined(func(pipe redis.Pipeliner) error {
		for _, k := range keys {
			pipe.SetNX(k, true, time.Duration(1*time.Second))
		}
		return nil
	})

	// Get the results of the pipelined commands
	for i, c := range cmds {
		cmd, _ := c.(*redis.BoolCmd)
		glog.Infof("SetNX for key '%s' (index: %d) returned '%v'", cmd.Args()[1], i, cmd.Val())
	}
}
