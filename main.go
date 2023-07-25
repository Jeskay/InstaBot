package main

import (
	"errors"
	config "instabot/src"
	"log"

	"github.com/Davincible/goinsta/v3"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print(".env file not found")
	}
}

var insta *goinsta.Instagram

func login() {
	if err := reloadSession(); err != nil {
		createSession()
	}
}

func reloadSession() error {
	inst, err := goinsta.Import("./goinsta-session")
	if err != nil {
		return errors.New("could not recover the session")
	}

	if inst != nil {
		insta = inst
	}

	log.Println("Restore session - SUCCESS")
	return nil
}

func createSession() {
	conf := config.NewConfig()
	inst := goinsta.New(conf.Instagram.Username, conf.Instagram.Password)
	insta = inst
	if err := insta.Login(); err != nil {
		panic(err)
	}

	if err := insta.Export("./goinsta-session"); err != nil {
		panic(err)
	}
	log.Println("Log in - SUCCESS")
}

func main() {
	login()
	log.Print(insta.Account.FullName)
}
