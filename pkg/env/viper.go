package env

import (
	"errors"
	"github.com/novando/go-cinema/pkg/logger"
	"github.com/spf13/viper"
	"strings"
)

// InitViper
// Initialize Viper to use the config file as env variable
func InitViper(path string, l ...*logger.Logger) error {
	log := logger.Call()
	if len(l) > 0 {
		log = l[0]
	}
	var configName string
	splitPaths := strings.Split(path, "/")
	if len(splitPaths) > 0 {
		for i := 0; i < len(splitPaths); i++ {
			configName = splitPaths[i]
		}
	}
	splitNames := strings.Split(configName, ".")
	if len(splitNames) < 2 {
		err := errors.New("failed to parse config name")
		if log != nil {
			log.Fatalf(err.Error())
		}
		return err
	}
	formatName := splitNames[len(splitNames)-1]
	viper.SetConfigName(strings.TrimRight(configName, "."+formatName))
	viper.SetConfigType(formatName)
	viper.AddConfigPath(strings.TrimRight(path, configName))
	err := viper.ReadInConfig()
	if err != nil && log != nil {
		logger.Call().Infof("Configs file: %v", err)
	}
	return err
}
