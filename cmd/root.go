package cmd

import (
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/xortim/peruse/conf"
	"go.uber.org/zap"

	"github.com/xortim/peruse/k8sclient"
)

var (
	cfgFile string
)

func newRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Version: conf.GitVersion,
		Use:     conf.Executable,
		RunE:    rootRun,
	}

	cobra.OnInitialize(initConfig)

	cmd.AddCommand(
		newVersionCmd(),
		newServCmd(),
	)

	cmd.PersistentFlags().StringVarP(&cfgFile, "configfile", "c", "", "ConfigFile to use instead of the default locations")
	cmd.PersistentFlags().String("kubeconfig", filepath.Join(conf.Home, ".kube", "config"), "Fully qualified path to the kubeconfig file")
	cmd.PersistentFlags().StringP("namespace", "n", "", "Limit the action to this namespace")

	cmd.MarkFlagRequired("kubeconfig")

	cmd.MarkFlagFilename("configfile")
	cmd.MarkFlagFilename("kubeconfig")

	viper.BindPFlags(cmd.PersistentFlags())

	return cmd
}

func rootRun(cmd *cobra.Command, args []string) error {
	k8s, err := k8sclient.NewClient("", viper.GetString("kubeconfig"))
	dips, err := k8sclient.GetDeploymentIngressPaths(k8s, viper.GetString("namespace"))
	if err != nil {
		return err
	}

	dips.FPrintTable(os.Stdout)
	return nil
}

func initConfig() {
	// If a config file is found, read it in.
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
		err := viper.ReadInConfig()
		if err != nil {
			zap.S().Fatalf("could not read config file %s ERROR: %s\n", cfgFile, err.Error())
		}
	}

	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		zap.S().Debugf("error loading config file: %s", err.Error())
	} else {
		zap.S().Debugf("using config file: %s", viper.ConfigFileUsed())
	}
}
