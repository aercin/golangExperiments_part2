package configs

import (
	"fmt"
	"go-poc/configs/abstractions"
	"os"
	"reflect"
	"strings"

	"github.com/spf13/viper"
)

type config struct {
	Postgres     postgresOptions
	HttpServer   httpServerOptions
	RabbitMQ     rabbitMqOptions
	MessageRelay messageRelayOptions
	Log          logOptions
}

type postgresOptions struct {
	Host         string
	Port         int64
	UserName     string
	Password     string
	DatabaseName string
}

type httpServerOptions struct {
	Port int
}

type rabbitMqOptions struct {
	BrokerAddress  string
	ProduceTimeout int
	ProduceQueue   string
	ConsumeQueue   string
}

type messageRelayOptions struct {
	CycleTime int
}

type logOptions struct {
	Path string
}

func NewConfig() abstractions.Config {

	environment := os.Getenv("APP_ENV")

	if environment == "" {
		panic("APP_ENV environment variable is not set")
	}

	viper.SetConfigName(fmt.Sprintf("appsettings.%s", environment))
	viper.SetConfigType("json")
	viper.AddConfigPath(".")
	viper.AddConfigPath("../../configs")

	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	var options *config
	if err := viper.Unmarshal(&options); err != nil {
		panic(fmt.Errorf("unable to decode into struct: %v", err))
	}

	return options
}

func (c *config) GetValue(key string) any {

	keys := strings.Split(key, ":")

	if len(keys) == 0 {
		panic("config - invalid key format")
	}

	var current reflect.Value = reflect.ValueOf(c).Elem()

	for _, k := range keys {
		current = current.FieldByName(k)
		if !current.IsValid() {
			panic(fmt.Errorf("field %s not found", k))
		}
	}

	return current.Interface()
}
