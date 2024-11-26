package url

import (
	"fmt"
	"os"
	"time"

	"github.com/rancher-government-carbide/rancher-cluster-id-finder/pkg/flags"
	"github.com/rancher-government-carbide/rancher-cluster-id-finder/pkg/kubernetes"
	"github.com/spf13/cobra"
)

var UrlCmd = &cobra.Command{
	Use:   "url",
	Short: "find url for rancher",
	RunE: func(cmd *cobra.Command, args []string) error {
		kc, err := kubernetes.NewKubeClient(flags.KubeconfigFile)
		if err != nil {
			return err
		}

		var url string
		currentRetries := 0

		url, err = kc.GetRancherURL()
		for (url == "") && (currentRetries < flags.Retries) {
			fmt.Println("Rancher URL Not found. Retrying..")
			time.Sleep(time.Duration(flags.Interval) * time.Second)
			url, err = kc.GetRancherURL()
			currentRetries += 1
		}

		if err != nil {
			return err
		}

		if flags.ConfigMapName != "" {
			err = kc.WriteConfigMap(url, flags.ConfigMapName, flags.ConfigMapNamespace, flags.ConfigMapKey)
			if err != nil {
				return err
			}
		}

		if flags.WriteFile != "" {
			err = os.WriteFile(flags.WriteFile, []byte(url), 0666)
		}

		return nil
	},
}
