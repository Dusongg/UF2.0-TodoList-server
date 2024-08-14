package config

import "os"

var (
	//Native_MySQL_DSN = "root:123123@tcp(127.0.0.1:3306)/OrderManager"
	//DSN = "root:root1234@tcp(127.0.0.1:13306)/test_docker_mysql"
	RedisHost = os.Getenv("REDIS_HOST")
	RedisPort = os.Getenv("REDIS_PORT")
	GormDNS   = os.Getenv("GORM_DNS")
)
