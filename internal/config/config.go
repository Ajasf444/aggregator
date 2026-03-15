package config

const configFileName = ".gatorconfig.json"

type Config struct {
	db_url            string
	current_user_name string
}

func Read(location string) Config {
	return Config{}
}
