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

func login(username string, password string) {
	if err := reloadSession(); err != nil {
		createSession(username, password)
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

func createSession(username string, password string) {
	inst := goinsta.New(username, password)
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
	conf := config.NewConfig()
	login(conf.Instagram.Username, conf.Instagram.Password)
	for _, donor := range conf.Instagram.Donors {
		if profile, err := insta.VisitProfile(donor); err == nil {
			user := profile.User
			log.Printf("%s has %d followers, %d posts, and %d IGTV vids \n", user.Username, user.FollowerCount, user.MediaCount, user.IGTVCount)
		}
	}
	log.Print(insta.Account.FullName)
}
