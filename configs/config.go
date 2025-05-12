package configs

import (
	"fmt"

	"github.com/spf13/viper"
)

type conf struct {
	DBDriver          string `mapstructure:"DB_DRIVER"`
	DBHost            string `mapstructure:"DB_HOST"`
	DBPort            string `mapstructure:"DB_PORT"`
	DBUser            string `mapstructure:"DB_USER"`
	DBPassword        string `mapstructure:"DB_PASSWORD"`
	DBName            string `mapstructure:"DB_NAME"`
	WebServerPort     string `mapstructure:"WEB_SERVER_PORT"`
	GRPCServerPort    string `mapstructure:"GRPC_SERVER_PORT"`
	GraphQLServerPort string `mapstructure:"GRAPHQL_SERVER_PORT"`
}

// LoadConfig loads configuration from a specific file path.
func LoadConfig(configFilePath string) (*conf, error) {
	var cfg *conf
	viper.SetConfigFile(configFilePath)
	viper.SetConfigType("env")

	err := viper.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("error reading config file %s: %w", configFilePath, err)
	}

	viper.AutomaticEnv()

	err = viper.Unmarshal(&cfg)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling config: %w", err)
	}
	return cfg, nil
}
