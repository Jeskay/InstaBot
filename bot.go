package main

import (
	"errors"
	utils "instabot/src/utils"
	"log"
	"time"

	"github.com/Davincible/goinsta/v3"
)

type Bot struct {
	instance     *goinsta.Instagram
	donors       *utils.IterableList[string]
	subs         *utils.IterableList[*goinsta.User]
	settings     *Settings
	hour_timeout time.Time
	sub_counter  int
}

type Settings struct {
	subs_per_hour int
	sub_interval  *utils.Spread
	hour_interval *utils.Spread
	condition_1   bool
	condition_2   bool
}

func NewBot(donors []string, settings *Settings) *Bot {
	return &Bot{
		instance: nil,
		donors:   utils.NewIterableList(donors),
		settings: settings,
	}
}

func (bot *Bot) Login(username, password string) {
	if err := bot.reloadSession(username); err != nil {
		bot.createSession(username, password)
	}
}

func (bot *Bot) reloadSession(username string) error {
	inst, err := goinsta.Import("./goinsta-session")
	if err != nil {
		return errors.New("could not recover the session")
	}
	if inst.Account.Username != username {
		return errors.New("session does not belong to the account")
	}
	if inst != nil {
		bot.instance = inst
	}

	return nil
}

func (bot *Bot) createSession(username string, password string) {
	inst := goinsta.New(username, password)
	bot.instance = inst
	if err := bot.instance.Login(); err != nil {
		panic(err)
	}

	if err := bot.instance.Export("./goinsta-session"); err != nil {
		panic(err)
	}
}

func (bot *Bot) resetSubPerHourTimeout() {
	bot.hour_timeout = time.Now().Add(1 * time.Hour)
	bot.sub_counter = 0
}

func (bot *Bot) hourlySleep() {
	hour_interval_duration := time.Duration(bot.settings.hour_interval.Rand()) * time.Second
	sleep_duration := time.Until(bot.hour_timeout)
	time.Sleep(sleep_duration)
	time.Sleep(hour_interval_duration)
	bot.resetSubPerHourTimeout()
}

func (bot *Bot) pageIteration(iteration func() error) {
	for bot.subs.Next() {
		time.Sleep(time.Duration(bot.settings.sub_interval.Rand()) * time.Second)
		log.Println(bot.hour_timeout.Clock())
		log.Println(time.Now())
		if bot.hour_timeout.Before(time.Now()) || bot.sub_counter >= bot.settings.subs_per_hour {
			bot.hourlySleep()
		}
		if err := iteration(); err != nil {
			log.Println(err)
			continue
		}
		bot.sub_counter++
	}
}

func (bot *Bot) subIteration() error {
	sub_profile, err := bot.subs.Current(true).VisitProfile()
	if err != nil {
		return err
	}
	if sub_profile.Friendship.Following {
		return errors.New("user is already followed")
	}
	if bot.settings.condition_1 && sub_profile.Friendship.FollowedBy {
		return errors.New("user is following core account")
	}
	if err := bot.subs.Current(true).Follow(); err != nil {
		return err
	}
	log.Println("Following profile of ", bot.subs.Current(true).Username, " - SUCCESS")
	return nil
}

func (bot *Bot) cleanIteration() error {
	profile, err := bot.subs.Current(true).VisitProfile()
	if err != nil {
		return err
	}
	if !bot.settings.condition_2 || !profile.Friendship.FollowedBy {
		log.Println("Checking profile of ", bot.subs.Current(true).Username, "- FAIL")
		err = profile.User.Unfollow()
		if err == nil {
			log.Println("Unfollowing profile of ", bot.subs.Current(true).Username, "- SUCCESS")
		}
		return err
	}
	log.Println("Checking profile of ", bot.subs.Current(true).Username, "- SUCCESS")
	return nil
}

func (bot *Bot) StartCleaningMode() {
	bot.resetSubPerHourTimeout()
	followers := bot.instance.Account.Following("", goinsta.DefaultOrder)
	for followers.Next() {
		bot.subs = utils.NewIterableList[*goinsta.User](followers.Users)
		bot.pageIteration(bot.cleanIteration)
	}
}

func (bot *Bot) StartFollowingMode() {
	bot.resetSubPerHourTimeout()
	for bot.donors.Next() {
		if profile, err := bot.instance.VisitProfile(bot.donors.Current(false)); err == nil {
			followers := profile.User.Followers("")
			for followers.Next() {
				bot.subs = utils.NewIterableList[*goinsta.User](followers.Users)
				bot.pageIteration(bot.subIteration)
			}
		} else {
			log.Println(err)
		}
	}
}
