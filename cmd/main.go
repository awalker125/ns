package main

import (
	"fmt"
	"os"

	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	if len(os.Args) > 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s [namespace]\n", os.Args[0])
		os.Exit(1)
	}

	// Empty string unsets the namespace in kubeconfig, reverting kubectl to "default"
	var targetNamespace string
	if len(os.Args) == 2 {
		targetNamespace = os.Args[1]
	}

	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	configAccess := clientcmd.NewDefaultPathOptions()

	config, err := loadingRules.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading kubeconfig: %v\n", err)
		os.Exit(1)
	}

	currentContextName := config.CurrentContext
	if currentContextName == "" {
		fmt.Fprintln(os.Stderr, "Error: no current-context is set")
		os.Exit(1)
	}

	currentContext, exists := config.Contexts[currentContextName]
	if !exists {
		fmt.Fprintf(os.Stderr, "Error: current-context %q does not exist in kubeconfig\n", currentContextName)
		os.Exit(1)
	}

	currentContext.Namespace = targetNamespace

	if err := clientcmd.ModifyConfig(configAccess, *config, true); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing kubeconfig: %v\n", err)
		os.Exit(1)
	}

	if targetNamespace == "" {
		fmt.Printf("Context %q namespace unset (now using \"default\").\n", currentContextName)
	} else {
		fmt.Printf("Context %q active namespace set to %q.\n", currentContextName, targetNamespace)
	}
}
