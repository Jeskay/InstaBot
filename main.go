package main

import (
	"instabot/src/extensions"
	"instabot/src/widgets"
	"log"
	"strings"
	"time"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print(".env file not found")
	}
}

func main() {
	application := app.New()

	settingsWindow := application.NewWindow("Settings")

	instagram_mode := binding.NewString()

	username_value := binding.NewString()
	disable_start_button := binding.NewBool()
	subs_per_hour_value := binding.NewInt()
	subs_per_hour_value.Set(10)

	subs_rate_entry := extensions.NewNumericalEntryWithData(binding.IntToString(subs_per_hour_value))
	mode_check := widgets.NewModeCheck(instagram_mode)
	spread_entry_1 := widgets.NewSpreadEntry(0, 10)
	spread_entry_2 := widgets.NewSpreadEntry(0, 10)

	bot_mode := widget.NewRadioGroup([]string{widgets.SubscribeMode, widgets.UnsubscribeMode}, func(s string) {
		instagram_mode.Set(s)
	})
	bot_mode.SetSelected(widgets.SubscribeMode)
	form_settings := widget.NewForm(
		widget.NewFormItem("Подписок в час", subs_rate_entry),
		widget.NewFormItem("Базовый интервал", spread_entry_1.Container),
		widget.NewFormItem("Часовой интервал", spread_entry_2.Container),
		widget.NewFormItem("Режим работы", bot_mode),
		widget.NewFormItem("Условия работы", mode_check.Container),
	)
	password_entry := widget.NewPasswordEntry()
	form_account := widget.NewForm(
		widget.NewFormItem("Имя пользователя", widget.NewEntryWithData(username_value)),
		widget.NewFormItem("Пароль", password_entry),
	)

	donors_entry := widget.NewMultiLineEntry()
	start_btn := widget.NewButton("Начать", func() {
		donors := strings.Split(donors_entry.Text, "\n")
		var subs_per_hour int = 10
		if value, err := subs_per_hour_value.Get(); err == nil {
			subs_per_hour = value
		}
		var sub_condition bool = false
		var unsub_condition bool = false
		if value, err := mode_check.SubConditionValue.Get(); err == nil {
			sub_condition = value
		}
		if value, err := mode_check.UnsubConditionValue.Get(); err == nil {
			unsub_condition = value
		}
		if !spread_entry_1.Spread.IsValid() {
			spread_entry_1.Spread.Max = spread_entry_1.Spread.Min + 3
		}
		if !spread_entry_2.Spread.IsValid() {
			spread_entry_2.Spread.Max = spread_entry_2.Spread.Min + 3
		}
		settings := &Settings{
			subs_per_hour: subs_per_hour,
			sub_interval:  spread_entry_1.Spread,
			hour_interval: spread_entry_2.Spread,
			condition_1:   sub_condition,
			condition_2:   unsub_condition,
		}
		bot := NewBot(donors, settings)
		if mode, err := instagram_mode.Get(); err == nil {
			disable_start_button.Set(true)
			var username string = ""
			if value, err := username_value.Get(); err == nil {
				username = value
			}
			go func() {
				bot.Login(username, password_entry.Text)
				if mode == widgets.SubscribeMode {
					time.Sleep(2 * time.Second)
					bot.StartFollowingMode()
					disable_start_button.Set(false)
				} else {
					time.Sleep(2 * time.Second)
					bot.StartCleaningMode()
					disable_start_button.Set(false)
				}
			}()

		}
	})
	disable_start_button.AddListener(binding.NewDataListener(func() {
		log.Println("Disabled button state changed")
		if value, err := disable_start_button.Get(); err == nil {
			if value && !start_btn.Disabled() {
				start_btn.Disable()
			} else if start_btn.Disabled() {
				start_btn.Enable()
			}
		}
	}))
	settingsContent := container.New(
		layout.NewVBoxLayout(),
		widget.NewLabel("Настройки"),
		form_settings,
		form_account,
		widget.NewLabel("Список доноров"),
		donors_entry,
		start_btn,
	)
	settingsWindow.SetContent(settingsContent)
	settingsWindow.ShowAndRun()

}
