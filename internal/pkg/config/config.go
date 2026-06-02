package config

import (
	"go-modular/internal/pkg/jwt"
	"log"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	filename string
}

func NewConfig(filename string) Config {
	return Config{filename: filename}
}
func (c *Config) Initialize() error {

	viper.SetConfigFile(c.filename)
	viper.SetConfigType("env")

	viper.AutomaticEnv()
	err := viper.ReadInConfig()

	if err != nil {
		// If the config file doesn't exist, fall back to environment variables
		// (e.g. when running in Docker with env_file in docker-compose.yml)
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil
		}
		// Also handle OS-level "file not found" errors
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	return nil
}

func checkKey(key string) {
	if !viper.IsSet(key) {
		log.Fatalf("Configuration key %s not found; aborting \n", key)
		os.Exit(1)
	}
}

func GetString(key string) string {
	checkKey(key)
	return viper.GetString(key)
}

func GetInt(key string) int {
	checkKey(key)
	return viper.GetInt(key)
}

func GetBool(key string) bool {
	checkKey(key)
	return viper.GetBool(key)
}

func GetJWTService() jwt.JWT {
	signatureKey := GetString("JWT_SIGNATURE_KEY")
	if signatureKey == "" {
		panic("JWT signature key not found in configuration")
	}
	return jwt.NewJWTImpl(signatureKey, 7)
}
