package utilities

import "github.com/spf13/viper"

func ReadConfigurations() {
	viper.SetConfigName("dnscheck.yaml")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("/etc/dnscheck/")
	err := viper.ReadInConfig() // Find and read the config file
	CheckError(err)
}
