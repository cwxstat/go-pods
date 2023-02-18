package context

import (
	"fmt"
	"k8s.io/client-go/tools/clientcmd"
	"os"
)

func Context() {
	// Use the default kubeconfig file location
	kubeconfig := clientcmd.NewDefaultClientConfigLoadingRules().GetDefaultFilename()

	// Load the configuration from the kubeconfig file
	config, err := clientcmd.LoadFromFile(kubeconfig)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading kubeconfig: %v\n", err)
		os.Exit(1)
	}

	// Get the current context from the configuration
	currentContext := config.CurrentContext

	fmt.Printf("Current Kubernetes context: %s\n", currentContext)
}

func All() {

	// Use the default kubeconfig file location
	kubeconfig := clientcmd.NewDefaultClientConfigLoadingRules().GetDefaultFilename()

	// Load the configuration from the kubeconfig file
	config, err := clientcmd.LoadFromFile(kubeconfig)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading kubeconfig: %v\n", err)
		os.Exit(1)
	}

	// Get the names of all available contexts from the configuration
	contexts := make([]string, 0, len(config.Contexts))
	for name := range config.Contexts {
		contexts = append(contexts, name)
	}

	fmt.Println("Available Kubernetes contexts:")
	for _, context := range contexts {
		fmt.Println(context)
	}

}
