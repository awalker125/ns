package main

import (
	"fmt"
	"os"

	"k8s.io/client-go/tools/clientcmd"
)

const (
	colorRed   = "\033[31m"
	colorGreen = "\033[32m"
	colorReset = "\033[0m"
)

func main() {
	if len(os.Args) > 2 {
		fmt.Fprintln(os.Stderr, formatMessage("error", colorRed, fmt.Sprintf("Usage: %s [namespace]", os.Args[0])))
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
		fmt.Fprintln(os.Stderr, formatMessage("error", colorRed, fmt.Sprintf("loading kubeconfig: %v", err)))
		os.Exit(1)
	}

	currentContextName := config.CurrentContext
	if currentContextName == "" {
		fmt.Fprintln(os.Stderr, formatMessage("error", colorRed, "no current-context is set"))
		os.Exit(1)
	}

	currentContext, exists := config.Contexts[currentContextName]
	if !exists {
		fmt.Fprintln(os.Stderr, formatMessage("error", colorRed, fmt.Sprintf("current-context %q does not exist in kubeconfig", currentContextName)))
		os.Exit(1)
	}

	currentContext.Namespace = targetNamespace

	if err := clientcmd.ModifyConfig(configAccess, *config, true); err != nil {
		fmt.Fprintln(os.Stderr, formatMessage("error", colorRed, fmt.Sprintf("writing kubeconfig: %v", err)))
		os.Exit(1)
	}

	if targetNamespace == "" {
		fmt.Println(formatMessage("info", colorGreen, "setting default namespace to: default"))
	} else {
		fmt.Println(formatMessage("info", colorGreen, fmt.Sprintf("setting default namespace to: %q.", targetNamespace)))
	}
}

func formatMessage(level, color, message string) string {
	return colorize(color, level) + ": " + message
}

func colorize(color, message string) string {
	if _, disabled := os.LookupEnv("NO_COLOR"); disabled {
		return message
	}
	return color + message + colorReset
}
