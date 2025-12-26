package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
	"github.com/shuail0/prediction-aggregator/pkg/exchange"
)

var log = logrus.New()

func main() {
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	log.SetLevel(logrus.InfoLevel)

	log.Info("🚀 Prediction Aggregator Starting...")

	// 创建上下文
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 监听终止信号
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		log.Info("📡 Received shutdown signal")
		cancel()
	}()

	// 演示：创建 Polymarket 客户端
	client, err := exchange.New("polymarket")
	if err != nil {
		log.Fatalf("❌ Failed to create exchange: %v", err)
	}

	log.Infof("✅ Created %s exchange", client.Name())
	log.Infof("🔗 Supported chains: %v", client.SupportedChains())

	// TODO: 加载配置
	// TODO: 连接到交易所
	// TODO: 启动策略

	fmt.Println("\n📊 Supported Platforms:")
	for _, p := range exchange.SupportedPlatforms() {
		fmt.Printf("  - %s\n", p)
	}

	<-ctx.Done()
	log.Info("👋 Prediction Aggregator Stopped")
}
