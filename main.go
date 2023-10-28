package main

import (
	"flag"

	genericserver "k8s.io/apiserver/pkg/server"
	"k8s.io/component-base/logs"
	"k8s.io/klog/v2"

	"github.com/kubernetescode-aaserver/cmd"
)

func main() {
	stopCh := genericserver.SetupSignalHandler()
	command := cmd.NewCommandStartServer(stopCh)
	command.Flags().AddGoFlagSet(flag.CommandLine)
	logs.InitLogs()
	defer logs.FlushLogs()
	if err := command.Execute(); err != nil {
		klog.Fatal(err)
	}
}
