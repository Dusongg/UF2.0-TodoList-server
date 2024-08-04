package main

import (
	"OrderManager/config"
	"OrderManager/pb"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"google.golang.org/grpc"
	"log"
	"net"
	"time"
)

const (
	TASK_STATE_WAIT   = 0
	TASK_STATE_ING    = 1
	TASK_STATE_FINISH = 2

	EMERGENCY_LEVEL_0 = 0
	EMERGENCY_LEVEL_1 = 1
	EMERGENCY_LEVEL_2 = 2
)

var (
	db *sql.DB
)

func init() {
	db, _ = sql.Open("mysql", config.DSN)

	db.SetMaxOpenConns(25)                 // 最大打开连接数
	db.SetMaxIdleConns(25)                 // 最大闲置连接数
	db.SetConnMaxLifetime(5 * time.Minute) // 连接的最大生命周期
	if err := db.Ping(); err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}
}

func main() {
	go emailClock()
	grpcServer := grpc.NewServer()
	pb.RegisterServiceServer(grpcServer, Server)
	listener, err := net.Listen("tcp", ":8001")
	if err != nil {
		log.Fatal("服务监听失败", err)
	}
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatal(err)
	}
	defer db.Close()
}

type task struct {
	comment            string
	taskId             string
	emergencyLevel     int32
	deadline           string
	principal          string
	reqNo              string
	estimatedWorkHours float32
	state              string
	typeId             int32
}

type patch struct {
	patchNo    string
	reqNo      string
	describe   string
	clientName string
	deadline   string
	reason     string
	sponsor    string
}
