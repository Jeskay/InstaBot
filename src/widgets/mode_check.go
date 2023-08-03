package widgets

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

type ModeCheck struct {
	Mode                  binding.String
	SubConditionValue     binding.Bool
	UnsubConditionValue   binding.Bool
	sub_condition_check   *widget.RadioGroup
	unsub_condition_check *widget.RadioGroup
	Container             *fyne.Container
}

const (
	SubscribeMode             string = "Подписка"
	UnsubscribeMode           string = "Отписка"
	SubscribeConditionTrue    string = "Подписка в ответ"
	SubscribeConditionFalse   string = "Подписка на всех"
	UnsubscribeConditionTrue  string = "Отписка от неподписавшихся"
	UnsubscribeConditionFalse string = "Отписка от всех"
)

func (m *ModeCheck) updateMode() {
	if value, err := m.Mode.Get(); err == nil {
		if value == SubscribeMode {
			m.sub_condition_check.Enable()
			m.unsub_condition_check.Disable()
		} else {
			m.sub_condition_check.Disable()
			m.unsub_condition_check.Enable()
		}
	}
}

func NewModeCheck(mode binding.String) *ModeCheck {
	modeCheck := &ModeCheck{
		Mode:                mode,
		SubConditionValue:   binding.NewBool(),
		UnsubConditionValue: binding.NewBool(),
	}
	modeCheck.sub_condition_check = widget.NewRadioGroup([]string{SubscribeConditionFalse, SubscribeConditionTrue}, func(s string) {
		modeCheck.SubConditionValue.Set(s == SubscribeConditionTrue)
	})
	modeCheck.unsub_condition_check = widget.NewRadioGroup([]string{SubscribeConditionFalse, UnsubscribeConditionTrue}, func(s string) {
		modeCheck.UnsubConditionValue.Set(s == UnsubscribeConditionTrue)
	})
	modeCheck.Container = container.NewVBox(
		modeCheck.sub_condition_check,
		modeCheck.unsub_condition_check,
	)
	modeCheck.Mode.AddListener(binding.NewDataListener(modeCheck.updateMode))
	return modeCheck
}
