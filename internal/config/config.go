package config

const configFileName = ".gatorconfig.json"

type Config struct {
	db_url            string
	current_user_name string
}
