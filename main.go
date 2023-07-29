package main

import (
	"errors"
	config "instabot/src/utils"
	"log"
	"time"

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
			log.Println("Visiting profile - ", user.Username)
			log.Printf("%s has %d followers, %d posts, and %d IGTV vids \n", user.Username, user.FollowerCount, user.MediaCount, user.IGTVCount)
			donorSubs := user.Followers("")
			for i := 0; i < donorSubs.PageSize; i++ {
				time.Sleep(time.Duration(conf.Instagram.SubInterval) * time.Second)
				donorSubs.Next()
				sub := donorSubs.Users[i]
				sub_profile, err := sub.VisitProfile()
				if err != nil {
					log.Println("Profile unavailable")
					log.Println(err)
					continue
				}

				if sub_profile.Friendship.Following {
					log.Println("User is already followed")
					continue
				}

				if conf.Instagram.Condition_2 && sub_profile.Friendship.FollowedBy {
					log.Println("User is following core account")
					continue
				}

				log.Println("Visited profile with ", sub_profile.Feed.NumResults, "posts")
				// if err := sub.Follow(); err != nil {
				// 	log.Println("Following profile - FAILED")
				// 	log.Println(err)
				// 	continue
				// }
				log.Println("Following profile - SUCCESS")

			}
		} else {
			log.Println(err)
		}
	}
	log.Print(insta.Account.FullName)
}
