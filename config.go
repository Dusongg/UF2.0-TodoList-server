package main

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"os"
	"strconv"
	"strings"
	"time"
)

var Conf = NewConfig("./config/config.json")

type Config struct {
	//Host    string `json:"host"`
	//Port    int    `json:"port"`
	//DBUser  string `json:"db_user"`
	//DBPass  string `json:"db_pass"`
	//DBName  string `json:"db_name"`

	Redis struct {
		Host string `json:"host"`
		Port string `json:"port"`
	}

	MySQL struct {
		GormDNS string `json:"gorm_dns"`
	}

	//MySQLRootPassword string `json:"mysql_root_password"`
	//MySQLDatabase string `json:"mysql_database"`
	//MySQLUser string `json:"mysql_user"`
	//MySQLUserPassword string `json:"mysql_user_password"`

	SMTP struct {
		Host           string `json:"host"`
		Port           string `json:"port"`
		Sender         string `json:"sender"`
		SenderPassword string `json:"email_sender_password"`
		SendTimePoint1 string `json:"send_time_point_1"`
		SendTimePoint2 string `json:"send_time_point_2"`
		P1             time.Time
		P2             time.Time
	}

	//日志存放路径
	LogPath string `json:"log_path"`

	//管理员用户
	Admins []string `json:"admins"`
}

func NewConfig(filePath string) *Config {
	var config Config
	data, err := os.ReadFile(filePath)
	if err != nil {
		logrus.Fatal(err)
	}
	err = json.Unmarshal(data, &config)
	if err != nil {
		logrus.Fatal(err)

	}
	config.overrideWithEnvVars()
	return &config
}

func (conf *Config) overrideWithEnvVars() {
	now := time.Now()
	if value, exists := os.LookupEnv("REDIS_HOST"); exists {
		conf.Redis.Port = value
	}
	if value, exists := os.LookupEnv("REDIS_PORT"); exists {
		conf.Redis.Host = value
	}
	if value, exists := os.LookupEnv("GORM_DNS"); exists {
		conf.MySQL.GormDNS = value
	}
	if value, exists := os.LookupEnv("EMAIL_SENDER"); exists {
		conf.SMTP.Sender = value
	}
	if value, exists := os.LookupEnv("EMAIL_SENDER"); exists {
		conf.SMTP.Sender = value
	}
	if value, exists := os.LookupEnv("EMAIL_SENDER"); exists {
		conf.SMTP.Sender = value
	}
	if value, exists := os.LookupEnv("EMAIL_SENDER_PASSWORD"); exists {
		conf.SMTP.SenderPassword = value
	}
	if value, exists := os.LookupEnv("SEND_MAIL_TIME_POINT1"); exists {
		hour, _ := strconv.Atoi(value)
		conf.SMTP.P1 = time.Date(now.Year(), now.Month(), now.Day(), hour, 0, 0, 0, now.Location())
	} else {
		hour, _ := strconv.Atoi(conf.SMTP.SendTimePoint1)
		conf.SMTP.P1 = time.Date(now.Year(), now.Month(), now.Day(), hour, 0, 0, 0, now.Location())
	}
	if value, exists := os.LookupEnv("SEND_MAIL_TIME_POINT2"); exists {
		hour, _ := strconv.Atoi(value)
		conf.SMTP.P2 = time.Date(now.Year(), now.Month(), now.Day(), hour, 0, 0, 0, now.Location())
	} else {
		hour, _ := strconv.Atoi(conf.SMTP.SendTimePoint2)
		conf.SMTP.P2 = time.Date(now.Year(), now.Month(), now.Day(), hour, 0, 0, 0, now.Location())
	}
	if value, exists := os.LookupEnv("LOG_PATH"); exists {
		conf.LogPath = value
	}
	if value, exists := os.LookupEnv("ADMINS"); exists {
		admins := strings.Split(value, ";")
		conf.Admins = admins
	}
}

/*
func main() {
    // 加载配置文件
    config, err := loadConfigFromFile("config.json")
    if err != nil {
        log.Fatalf("Error loading config file: %v", err)
    }

    // 覆盖配置文件中的值（如果环境变量存在）
    overrideWithEnvVars(config)

    // 使用最终的配置
    fmt.Printf("App Name: %s\n", config.AppName)
    fmt.Printf("Server Host: %s\n", config.Host)
    fmt.Printf("Database User: %s\n", config.DBUser)
}

*/
