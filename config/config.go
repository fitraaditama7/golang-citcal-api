package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var (
	PORT   = 0
	APIKEY = ""
)

/*
 *
 * @desc Load all configuration in .env
 *
 */
func Load() {
	var err error
	err = godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	PORT, err = strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		PORT = 4500
	}

	APIKEY = os.Getenv("API_KEY_CITCALL")
}
