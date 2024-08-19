package main

import (
	"OrderManager/models"
	"OrderManager/pb"
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/reflection"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
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

var adminMap = make(map[string]bool)

func isAdmin(in string) bool {
	_, ok := adminMap[in]
	return ok
}
func init() {
	for _, admin := range Conf.Admins {
		adminMap[admin] = true
	}
	//logrus.Error(NotificationServer.rdb.Ping(context.Background()).Err())

	// 设置日志
	//logrus.SetOutput(&lumberjack.Logger{
	//	Filename:   Conf.LogPath,
	//	MaxSize:    100, // MB
	//	MaxBackups: 30,
	//	MaxAge:     0, // Disable age-based rotation
	//	Compress:   true,
	//})

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	//测试
	logrus.SetOutput(os.Stdout)
	logrus.SetFormatter(&logrus.JSONFormatter{})

	if Conf.Redis.Port == "" || Conf.Redis.Host == "" || Conf.MySQL.GormDNS == "" {
		logrus.Fatalf("服务器链接数据库错误：redisHost: %s, redisPort: %s, GormDNS: %s", Conf.Redis.Port, Conf.Redis.Host, Conf.MySQL.GormDNS)
	}

	tmpDb, err := gorm.Open(mysql.Open(Conf.MySQL.GormDNS), &gorm.Config{})
	if err != nil {
		logrus.Fatal("Failed to connect to database:", err, Conf.MySQL.GormDNS)
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
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		// 执行清理操作
		for _, cliName := range NotificationServer.clients {
			NotificationServer.rdb.Del(context.Background(), fmt.Sprintf("user:%s", cliName))
			logrus.Info("Unsubscribed client:", cliName)
		}
		os.Exit(0)
	}()
	//sigs := make(chan os.Signal, 1)
	//signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	if Conf.SMTP.Switch == "on" {
		go emailClock()
	}
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

	/*
		gRPC 反射是一种机制，允许 gRPC 客户端在不知道服务的 proto 文件的情况下，动态查询服务的描述符信息（包括服务名、方法名、消息结构等）。
		grpcurl 依赖于这个功能来调用 gRPC 方法，如果服务器没有启用反射，你需要手动提供 .proto 文件的信息。
	*/
	//启用反射服务
	reflection.Register(grpcServer)

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
