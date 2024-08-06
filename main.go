package main

import (
	"OrderManager/config"
	"OrderManager/models"
	"OrderManager/pb"
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"net"
)

const (
	TASK_STATE_WAIT   = 0
	TASK_STATE_ING    = 1
	TASK_STATE_FINISH = 2

	EMERGENCY_LEVEL_0 = 0
	EMERGENCY_LEVEL_1 = 1
	EMERGENCY_LEVEL_2 = 2
)

var db *gorm.DB

type TaskInfo = models.TaskInfo

type PatchsInfo = models.PatchsInfo
type UserInfo = models.UserInfo

func init() {
	tmpDb, err := gorm.Open(mysql.Open(config.GORM_DNS), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	db = tmpDb
	err = db.AutoMigrate(&TaskInfo{}, &PatchsInfo{}, &UserInfo{})
	if err != nil {
		log.Fatal(err)
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
		log.Printf("Received request from:%s", p.String())
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
}
