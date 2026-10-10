package healer

import (
	"context"
	"fmt"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type AutoHealer struct {
	client    *kubernetes.Clientset
	namespace string
}

func NewAutoHealer(client *kubernetes.Clientset, namespace string) *AutoHealer {
	return &AutoHealer{
		client:    client,
		namespace: namespace,
	}
}

func (h *AutoHealer) CheckAndHeal(ctx context.Context) error {
	pods, err := h.client.CoreV1().Pods(h.namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return fmt.Errorf("failed to list pods: %w", err)
	}

	fmt.Printf("[%s] 🔍 Scanning %d pods in namespace '%s'...\n", time.Now().Format("15:04:05"), len(pods.Items), h.namespace)

	for _, pod := range pods.Items {
		for _, cs := range pod.Status.ContainerStatuses {
			needsHealing := false
			reason := ""

			if cs.State.Waiting != nil && (cs.State.Waiting.Reason == "CrashLoopBackOff" || cs.State.Waiting.Reason == "Error") {
				needsHealing = true
				reason = cs.State.Waiting.Reason
			} else if cs.RestartCount > 10 {
				needsHealing = true
				reason = fmt.Sprintf("ExcessiveRestarts (%d)", cs.RestartCount)
			}

			if needsHealing {
				fmt.Printf("⚠️ Detected unhealthy pod: %s (Reason: %s). Remediating...\n", pod.Name, reason)
				err := h.client.CoreV1().Pods(pod.Namespace).Delete(ctx, pod.Name, metav1.DeleteOptions{})
				if err != nil {
					fmt.Printf("❌ Failed to delete/restart pod %s: %v\n", pod.Name, err)
				} else {
					fmt.Printf("✅ Pod %s remediated (deleted for controller replacement).\n", pod.Name)
				}
				break
			}
		}
	}

	return nil
}
