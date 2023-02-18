package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	// Load the Kubernetes configuration from the default location or a specified path.
	kubeconfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		clientcmd.NewDefaultClientConfigLoadingRules(),
		&clientcmd.ConfigOverrides{})

	config, err := kubeconfig.ClientConfig()
	if err != nil {
		log.Fatalf("Error loading Kubernetes configuration: %v", err)
	}

	// Create a Kubernetes client using the configuration.
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("Error creating Kubernetes client: %v", err)
	}

	// List all the pods in the "dev" namespace.
	pods, err := clientset.CoreV1().Pods("").List(context.Background(), metav1.ListOptions{})
	if err != nil {
		log.Fatalf("Error listing pods: %v", err)
	}

	// Iterate over the pods and print their associated events and logs.
	for _, pod := range pods.Items {
		fmt.Printf("%s: Pod %s: status: %v,\n", pod.Namespace, pod.Name, pod.Status.Phase)

		if pod.Status.Phase == corev1.PodRunning {
			fmt.Printf("Pod %s:\n", pod.Name)

			// Get the logs for the pod's containers.
			podLogs := clientset.CoreV1().Pods(pod.Namespace).GetLogs(pod.Name, &corev1.PodLogOptions{})
			logs, err := podLogs.Stream(context.Background())
			if err != nil {
				log.Printf("Error getting logs for pod %s: %v", pod.Name, err)
				continue
			}
			defer logs.Close()
			buf := new(bytes.Buffer)
			_, err = io.Copy(buf, logs)
			if err != nil {
				log.Printf("Error reading logs for pod %s: %v", pod.Name, err)
				continue
			}
			fmt.Printf("  Logs: %s\n", buf.String())
		}

	}

	fmt.Println("\n\n\nEvents:")
	events, err := clientset.CoreV1().Events("").List(context.Background(), metav1.ListOptions{})
	if err != nil {
		log.Fatalf("Error listing events: %v", err)
	}

	// Print the events.
	for _, event := range events.Items {
		fmt.Printf("Event: %s %s %s\n", event.Reason, event.Type, event.Message)
	}
}
