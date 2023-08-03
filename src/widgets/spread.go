package widgets

import (
	"instabot/src/extensions"
	"instabot/src/utils"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

type SpreadEntry struct {
	Spread    *utils.Spread
	MinValue  binding.Int
	MaxValue  binding.Int
	Container *fyne.Container
}

func (e *SpreadEntry) updateMin() {
	if value, err := e.MinValue.Get(); err == nil {
		e.Spread.Min = value
	}
}

func (e *SpreadEntry) updateMax() {
	if value, err := e.MaxValue.Get(); err == nil {
		e.Spread.Max = value
	}
}

func NewSpreadEntry(min int, max int) *SpreadEntry {
	spread_entry := &SpreadEntry{
		Spread:   utils.NewSpread(min, max),
		MinValue: binding.NewInt(),
		MaxValue: binding.NewInt(),
	}
	spread_entry.MinValue.Set(min)
	spread_entry.MaxValue.Set(max)

	container := container.NewVBox(
		container.NewGridWithColumns(
			5,
			widget.NewLabel("Min"),
			extensions.NewNumericalEntryWithData(binding.IntToString(spread_entry.MinValue)),
		),
		container.NewGridWithColumns(
			5,
			widget.NewLabel("Max"),
			extensions.NewNumericalEntryWithData(binding.IntToString(spread_entry.MaxValue)),
		),
	)
	spread_entry.Container = container
	spread_entry.MinValue.AddListener(binding.NewDataListener(spread_entry.updateMin))
	spread_entry.MaxValue.AddListener(binding.NewDataListener(spread_entry.updateMax))

	return spread_entry
}
