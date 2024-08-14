package config

import "github.com/spf13/viper"

type Conf struct {
	DBUri                    string `mapstructure:"DB_URI"`
	DBDatabase               string `mapstructure:"DB_DATABASE"`
	DBSessionCollection      string `mapstructure:"DB_SESSION_COLLECTION"`
	DBPatientCollection      string `mapstructure:"DB_PATIENT_COLLECTION"`
	DBProfessionalCollection string `mapstructure:"DB_PROFESSIONAL_COLLECTION"`
	WebServerPort            string `mapstructure:"WEB_SERVER_PORT"`
	OtelExporterEndpoint     string `mapstructure:"OTEL_EXPORTER_OTLP_ENDPOINT"`
}

func LoadConfig(paths []string) (*Conf, error) {
	var cfg *Conf
	viper.SetConfigName("app_config")
	viper.SetConfigType("env")
	for _, path := range paths {
		viper.AddConfigPath(path)
	}
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	err = viper.Unmarshal(&cfg)
	if err != nil {
		panic(err)
	}
	return cfg, err
}
