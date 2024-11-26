package id

import (
	"fmt"
	"os"
	"time"

	"github.com/rancher-government-carbide/rancher-cluster-id-finder/pkg/flags"
	"github.com/rancher-government-carbide/rancher-cluster-id-finder/pkg/kubernetes"
	"github.com/spf13/cobra"
)

var IdCmd = &cobra.Command{
	Use:   "id",
	Short: "find id for rancher cluster",
	RunE: func(cmd *cobra.Command, args []string) error {
		kc, err := kubernetes.NewKubeClient(flags.KubeconfigFile)
		if err != nil {
			return err
		}

		var rancherClusterId string
		var valid bool
		currentRetries := 0

		rancherClusterId, valid = kc.CheckLocalCluster()
		if !valid {
			rancherClusterId, valid = kc.GetClusterIDFromConfigMap()
			if !valid {
				rancherClusterId, valid = kc.GetClusterIDFromSecret()
				if !valid {
					rancherClusterId, valid = kc.GetClusterIDFromAnnotations()
				}
			}
		}

		for (rancherClusterId == "") && (currentRetries < flags.Retries) {
			fmt.Println("Rancher URL Not found. Retrying..")
			time.Sleep(time.Duration(flags.Interval) * time.Second)
			rancherClusterId, valid = kc.CheckLocalCluster()
			if !valid {
				rancherClusterId, valid = kc.GetClusterIDFromConfigMap()
				if !valid {
					rancherClusterId, valid = kc.GetClusterIDFromSecret()
					if !valid {
						rancherClusterId, valid = kc.GetClusterIDFromAnnotations()
					}
				}
			}
			currentRetries += 1
		}

		if !valid {
			return fmt.Errorf("ERROR: Retries exceeded. Could not get Cluster ID.")
		}

		if flags.ConfigMapName != "" {
			err = kc.WriteConfigMap(rancherClusterId, flags.ConfigMapName, flags.ConfigMapNamespace, flags.ConfigMapKey)
			if err != nil {
				return err
			}
		}

		if flags.WriteFile != "" {
			err = os.WriteFile(flags.WriteFile, []byte(rancherClusterId), 0666)
		}

		fmt.Print(rancherClusterId)

		return nil
	},
}
