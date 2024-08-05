package main

import (
	"OrderManager/config"
	"OrderManager/pb"
	"context"
	"database/sql"
	"log"
	"net"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
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

func unaryInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	// Retrieve client information
	p, ok := peer.FromContext(ctx)
	if ok {
		log.Printf("Received request from: %s", p.Addr)
	}
	return handler(ctx, req)
}

func main() {
	go emailClock()
	//go testSendEmail()
	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(unaryInterceptor))
	pb.RegisterServiceServer(grpcServer, Server)
	listener, err := net.Listen("tcp", ":8001")
	if err != nil {
		log.Fatal("服务监听失败", err)
	} else {
		log.Println("正在监听端口：", listener.Addr())
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
