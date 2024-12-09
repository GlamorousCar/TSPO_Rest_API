package config

import (
	"os"
)

type Configuration struct {
	ServerAddr       string
	ServerPort       string
	LogLevel         string //INFO, DEBUG, WARNING
	DBFlavor         string
	DBHost           string
	DBPort           string
	DBUser           string
	DBPass           string
	DBName           string
	JWTSecret        string
	JWTRefreshSecret string
}

// Construct() использует метод os.LookupEnv() для получения значений переменных окружения.
// Если переменная окружения не установлена, метод устанавливает значение по умолчанию.
func (c *Configuration) Construct() {
	c.ServerAddr = os.Getenv("API_SERVER_ADDR")

	c.DBHost = os.Getenv("DB_HOST")
	c.DBFlavor = "postgres"
	c.DBPort = os.Getenv("DB_PORT")
	c.DBUser = os.Getenv("DB_USER")
	c.DBPass = os.Getenv("DB_PASSWORD")
	c.DBName = os.Getenv("DB_NAME")
	c.JWTSecret = "your-secret-key"
	c.JWTRefreshSecret = "your-refresh-secret-key"

}
