package config

// Server is a struct to use in config
type Server struct {
	Port string `mapstructure:"port"`
}

type FiscalData struct {
	URL string `mapstructure:"url"`
}

type Redis struct {
	Password string `mapstructure:"password"`
	URL      string `mapstructure:"url"`
	DB       int    `mapstructure:"db"`
}

type MongoDb struct {
	URL    string `mapstructure:"url"`
	Scheme string `mapstructure:"scheme"`
}

type Config struct {
	ENV           string     `mapstructure:"env"`
	Server        Server     `mapstructure:"server"`
	FiscalData    FiscalData `mapstructure:"fiscaldata"`
	Redis         Redis      `mapstructure:"redis"`
	MongoDbReader MongoDb    `mapstructure:"mongodb_reader"`
	MongoDbWriter MongoDb    `mapstructure:"mongodb_writer"`
}

// GlobalConfig is you use in all app
var GlobalConfig *Config

func init() {
	GlobalConfig = new(Config)
}
