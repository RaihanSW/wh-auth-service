package config

import (
	"fmt"
	"github.com/cinemaxxi-team/xxi-library/nrredigo"
	"strconv"
	"time"

	// "github.com/cinemaxxi-team/xxi-library/nrredigo"
	"github.com/cinemaxxi-team/xxi-library/redis"
	redigo "github.com/gomodule/redigo/redis"
)

// Connect to host:port, host may be omitted, so ":6379" is valid.
// Will not try to AUTH with any given password (password@host:port).
func newRedisConnectionTo(connectionString string) (redigo.Conn, error) {
	return redigo.Dial("tcp", connectionString)
}

// Create a new connection pool given a host:port string.
// A password may be supplied as well, on the form "password@host:port".
func NewRedisConnectionPool(cfg *ApplicationConfig) *redis.Pool {
	conString := fmt.Sprintf("%s:%s", cfg.RedisCfg.Address, cfg.RedisCfg.Port)

	pool := redis.InitRedis(redis.Opts{
		MaxIdle:     cfg.RedisCfg.MaxIdleConnections,
		MaxActive:   cfg.RedisCfg.MaxActiveConnections,
		IdleTimeout: cfg.RedisCfg.IdleTimeout * time.Second,
		Dial: func() (redigo.Conn, error) {
			c, err := newRedisConnectionTo(conString)
			if err != nil {
				return nil, err
			}

			//check password if there is password set
			if cfg.RedisCfg.Password != "" {
				if _, err := c.Do("AUTH", cfg.RedisCfg.Password); err != nil {
					c.Close()
					return nil, err
				}
			}

			if cfg.RedisCfg.DB > 0 {
				if _, err := c.Do("SELECT", strconv.Itoa(cfg.RedisCfg.DB)); err != nil {
					c.Close()
					return nil, err
				}
			}

			return c, err
		},
		TestOnBorrow: func(c redigo.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	})

	pool.AddHook(nrredigo.NewHook(&nrredigo.Config{
		Host:         cfg.RedisCfg.Address,
		PortPathOrID: cfg.RedisCfg.Port,
		DBName:       "",
	}))

	return pool

	//pool := ConnectionPool(*redisPool)
	//return &pool
}
