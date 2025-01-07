package configs

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/renebizelli/goexpert/desafios/multithreading/internal/utils"

	"github.com/spf13/viper"
)

type (
	Config struct {
		WebServer WebServerConfig
		Services  ServicesConfig
	}

	WebServerConfig struct {
		Port string
	}

	ServicesConfig struct {
		ViacepUrl    string
		BrasilApiUrl string
		Timeout      int
	}
)

func LoadConfig(path string) Config {

	viper.SetConfigName("app_config")
	viper.SetConfigType("dotenv")
	viper.AddConfigPath(path)
	viper.SetConfigFile(".env")

	viper.AutomaticEnv()

	err := viper.ReadInConfig()

	utils.PanicIfError(err, "Load config file error")

	var config Config

	valueRef := reflect.ValueOf(&config)
	typeRef := reflect.TypeOf(&config)

	for i := 0; i < typeRef.Elem().NumField(); i++ {

		mainProperty := typeRef.Elem().Field(i)

		for j := 0; j < mainProperty.Type.NumField(); j++ {

			property := mainProperty.Type.Field(j)

			envPropertyName := strings.ToLower(strings.Join([]string{mainProperty.Name, property.Name}, "_"))

			value := viper.Get(envPropertyName)

			fmt.Printf("Property %s: %s\n", utils.CyanText(envPropertyName), value)

			if value != nil {

				refValueProperty := valueRef.Elem().FieldByName(mainProperty.Name).FieldByName(property.Name)

				switch property.Type.Kind() {
				case reflect.String:
					refValueProperty.SetString(value.(string))
				case reflect.Int64:
				case reflect.Int:
					v, e := strconv.ParseInt(value.(string), 0, 64)
					utils.PanicIfError(e, fmt.Sprintf("Invalid value for key: %s, value: %s", envPropertyName, value))
					refValueProperty.SetInt(v)
				default:
					fmt.Printf("Property %s: %s\n", utils.CyanText(envPropertyName), property.Type.Kind())
				}

			} else {
				fmt.Printf("Property %s is nil\n", utils.CyanText(envPropertyName))
			}
		}
	}

	return config
}
