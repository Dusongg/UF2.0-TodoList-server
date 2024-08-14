package main

import (
	"OrderManager/config"
	"OrderManager/models"
	"OrderManager/pb"
	"context"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
	"gopkg.in/natefinch/lumberjack.v2"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
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
	if config.GormDNS == "" {
		config.GormDNS = "root:123123@tcp(127.0.0.1:3306)/OrderManager?charset=utf8mb4&parseTime=True&loc=Local" //本地测试
	}
	if config.RedisHost == "" {
		config.RedisHost = "redis"
	}
	if config.RedisPort == "" {
		config.RedisPort = "6379"
	}
	//log.SetFlags(log.LstdFlags | log.Lshortfile)

	tmpDb, err := gorm.Open(mysql.Open(config.GormDNS), &gorm.Config{})
	if err != nil {
		logrus.Fatal("Failed to connect to database:", err, config.GormDNS)
	}
	db = tmpDb
	err = db.AutoMigrate(&TaskInfo{}, &PatchsInfo{}, &UserInfo{})
	if err != nil {
		logrus.Fatal(err)
	}
}

// 获得客户端ip端口
func unaryInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	// Retrieve client information
	p, ok := peer.FromContext(ctx)
	if ok {
		logrus.Infof("Received request from:%s", p.String())
	}
	return handler(ctx, req)
}

func main() {
	logrus.SetOutput(&lumberjack.Logger{
		Filename:   "./logs/app.log",
		MaxSize:    100, // MB
		MaxBackups: 30,
		MaxAge:     0, // Disable age-based rotation
		Compress:   true,
	})

	//sigs := make(chan os.Signal, 1)
	//signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go emailClock()
	//go func() {
	//	for {
	//		time.Sleep(3 * time.Second)
	//		log.Println("client nums: ", len(NotificationServer.clients))
	//	}
	//}()
	//go testSendEmail()
	//grpcServer := grpc.NewServer(grpc.UnaryInterceptor(unaryInterceptor))
	grpcServer := grpc.NewServer()
	defer grpcServer.GracefulStop()

	pb.RegisterServiceServer(grpcServer, Server)

	pb.RegisterNotificationServiceServer(grpcServer, NotificationServer)

	//go NotificationServer.publishUpdates() //启动 Redis 订阅者

	listener, err := net.Listen("tcp", ":8001")
	if err != nil {
		logrus.Info("服务监听失败", err)
	} else {
		logrus.Info("正在监听端口：", listener.Addr())
	}
	if err := grpcServer.Serve(listener); err != nil {
		logrus.Fatal(err)
	}

	//<-sigs
}
