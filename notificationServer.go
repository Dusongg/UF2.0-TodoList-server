package main

import (
	"OrderManager/pb"
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	"log"
	"regexp"
	"sync"
	"time"
)

type notificationServer struct {
	pb.UnimplementedNotificationServiceServer
	clients map[string]pb.NotificationService_SubscribeServer
	mu      sync.Mutex
	rdb     *redis.Client
	ctx     context.Context
}

func (ns *notificationServer) checkIfLoggedIn(userName string) bool {
	_, ok := ns.clients[userName]
	return ok
}

var NotificationServer = &notificationServer{
	clients: make(map[string]pb.NotificationService_SubscribeServer),
	rdb:     redis.NewClient(&redis.Options{Addr: fmt.Sprintf("%s:%s", Conf.Redis.Host, Conf.Redis.Port)}),
	ctx:     context.Background(),
}

func (ns *notificationServer) updateDatabaseAndNotify(updateData string) {
	log.Println(updateData)
	msg := fmt.Sprintf("%s: %s -- refresh to view", time.Now().Format("2006-01-02 15:04:05"), updateData)
	if err := ns.rdb.Publish(context.Background(), "updates", msg).Err(); err != nil {
		logrus.Warningf("Error updating database: %v", err)
	}
}

type modDeadlineInfo struct {
	patchNo     string
	reqNo       string
	newDeadline string
	user        string
}

// Subscribe 多个客户端调用
func (ns *notificationServer) Subscribe(req *pb.SubscriptionRequest, stream pb.NotificationService_SubscribeServer) error {
	ns.mu.Lock()
	ns.clients[req.ClientId] = stream
	logrus.Infof("%s has subscribed\n", req.ClientId)
	ns.mu.Unlock()

	//TODO: 没有做持久化
	pubsub := ns.rdb.Subscribe(ns.ctx, "updates")
	ch := pubsub.Channel()

	go func() {
		<-stream.Context().Done()
		pubsub.Close() // 关闭 Redis 订阅以退出 for 循环
	}()

	for msg := range ch {

		re := regexp.MustCompile(`<([^>]*)>`)
		matches := re.FindAllStringSubmatch(msg.Payload, -1)
		if len(matches) != 2 {
			continue
		}
		from := matches[0][1]
		to := matches[1][1]
		if from == req.ClientId || (to != "ALL" && to != req.ClientId) {
			//log.Printf("from:%s, to:%s, name: %s", from, to, req.ClientId)
			continue
		}
		if err := stream.Send(&pb.Notification{Message: msg.Payload}); err != nil {
			return err
		}
	}

	ns.mu.Lock()
	delete(ns.clients, req.ClientId)
	logrus.Info("Unsubscribed client:", req.ClientId)
	ns.mu.Unlock()

	return nil
}
