package cmd

import (
	"fmt"
	"os"

	"github.com/rancher-government-carbide/rancher-cluster-id-finder/cmd/id"
	"github.com/rancher-government-carbide/rancher-cluster-id-finder/cmd/url"
	"github.com/rancher-government-carbide/rancher-cluster-id-finder/pkg/flags"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "rcidf",
	Short: "rcidf: rancher cluster id finder. also finds rancher url!",
	Long:  "rcidf is used to find the cluster id of a downstream rancher cluster. it can also grab a rancher url",
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&flags.Debug, "debug", false, "enable debug logging")
	rootCmd.PersistentFlags().StringVar(&flags.KubeconfigFile, "kubeconfig", "", "path to kubeconfig file")
	rootCmd.PersistentFlags().StringVar(&flags.ConfigMapName, "configmap-name", "", "name of configmap to create")
	rootCmd.PersistentFlags().StringVar(&flags.ConfigMapNamespace, "configmap-namespace", "", "namespace of configmap")
	rootCmd.PersistentFlags().StringVar(&flags.ConfigMapKey, "configmap-key", "", "key in configmap")
	rootCmd.PersistentFlags().StringVar(&flags.WriteFile, "write-file", "", "path to write output")
	rootCmd.PersistentFlags().IntVar(&flags.Retries, "retries", 10, "number of retries to acquire metadata before failure")
	rootCmd.PersistentFlags().IntVar(&flags.Interval, "interval", 5, "internal in seconds between retries")

	rootCmd.AddCommand(id.IdCmd)
	rootCmd.AddCommand(url.UrlCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
