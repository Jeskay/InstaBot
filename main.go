package main

import (
	config "instabot/src/utils"
	"log"

	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print(".env file not found")
	}
}

func main() {
	conf := config.NewConfig()
	settings := &Settings{
		subs_per_hour: 10,
		sub_interval:  config.NewSpread(3, 5),
		hour_interval: config.NewSpread(1, 10),
		condition_1:   true,
	}
	inst := NewBot(conf.Instagram.Donors, settings)
	inst.Login(conf.Instagram.Username, conf.Instagram.Password)
	inst.Start()
}
