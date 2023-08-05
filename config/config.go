package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

const ENV = ".env"
const MONGO_URI = "MONGO_URI"
const POSTGRES_URI = "POSTGRES_URI"

type Config struct {
	MongoUri    string
	PostgresUri string
}

func Load() Config {
	err := godotenv.Load(ENV)
	if err != nil {
		log.Fatal(err)
	}
	cfg := Config{
		MongoUri:    ReadFromEnvFile(MONGO_URI),
		PostgresUri: ReadFromEnvFile(POSTGRES_URI),
	}
	return cfg
}

func ReadFromEnvFile(key string) string {
	return os.Getenv(key)
}
