package util

import (
	"github.com/spf13/viper"
	"log"
)

// Name for the config file.
const ViperConfigName = "madprobe"

// Type of the config file.
const ViperConfigType = "yml"

// Path to search the config file.
const ViperConfigPath = "."

func init() {
	configurationFile()
	viper.AutomaticEnv()
}

func configurationFile() {
	viper.SetConfigName(ViperConfigName)
	viper.SetConfigType(ViperConfigType)
	viper.AddConfigPath(ViperConfigPath)
	if err := viper.ReadInConfig(); err != nil {
		log.Printf("[WARNING] an error occured when reading viper configuration file. got: [%v]\n", err)
	}
}
