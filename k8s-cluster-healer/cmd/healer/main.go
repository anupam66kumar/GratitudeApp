package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"k8s-cluster-healer/pkg/healer"
	"k8s-cluster-healer/pkg/k8s"
)

func main() {
	fmt.Println("==================================================")
	fmt.Println("🚀 Initializing Kubernetes Cluster Auto-Healer...")
	fmt.Println("==================================================")

	client, err := k8s.GetKubeClient()
	if err != nil {
		fmt.Printf("❌ Failed to obtain Kubernetes client: %v\n", err)
		os.Exit(1)
	}

	namespace := os.Getenv("TARGET_NAMESPACE")
	if namespace == "" {
		namespace = "default"
	}

	autoHealer := healer.NewAutoHealer(client, namespace)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	// Initial scan run
	if err := autoHealer.CheckAndHeal(ctx); err != nil {
		fmt.Printf("⚠️ Scan error: %v\n", err)
	}

	fmt.Printf("🔄 Auto-Healer active. Monitoring namespace '%s' every 10s. Press Ctrl+C to stop.\n", namespace)

	for {
		select {
		case <-ctx.Done():
			fmt.Println("\n🛑 Gracefully shutting down Auto-Healer daemon.")
			return
		case <-ticker.C:
			if err := autoHealer.CheckAndHeal(ctx); err != nil {
				fmt.Printf("⚠️ Scan error: %v\n", err)
			}
		}
	}
}
