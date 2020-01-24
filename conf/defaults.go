package conf

import (
	"path/filepath"
	"strings"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var (
	// Home is set in init()
	Home = ""
)

// Using init to create the starting state of the configuration when this package is first imported
func init() {
	Home, err := homedir.Dir()
	if err != nil {
		zap.S().Fatalw("error finding home directory", err)
	}

	// Search config in home directory with name ".[Executable]" (without extension).
	viper.AddConfigPath(".")
	viper.AddConfigPath(Home)
	viper.SetConfigName("." + Executable)
	viper.SetTypeByDefaultValue(true)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.SetDefault("kubeconfig", filepath.Join(Home, ".kube", "config"))
	viper.SetDefault("namespace", "")
}
