package utils

import (
	"fmt"
	"strings"
	"time"
)

type MomentDuration string

const (
	Nanosecond  MomentDuration = "Nanosecond"
	Microsecond MomentDuration = "Microsecond"
	Millisecond                = "Millisecond"
	Second                     = "Second"
	Minute                     = "Minute"
	Hour                       = "Hour"
	Day                        = "Day"
	Month                      = "Month"
	Year                       = "Year"
)

type Moment struct {
	t         time.Time
	tz_offset int
}

func NewMoment(t time.Time) *Moment {
	return &Moment{t: t}
}
func (m *Moment) WithTZOffset(tz_offset int) *Moment {
	m.tz_offset = tz_offset
	return m
}

func (m *Moment) ToDate() time.Time {
	return m.t
}
func (m *Moment) ToUnix() int64 {
	return m.ToDate().Unix()
}

func (m *Moment) StartOfDate() *Moment {
	m.pretz()
	ti := time.Date(m.t.Year(), m.t.Month(), m.t.Day(), 0, 0, 0, 0, m.t.Location())
	m.t = ti
	m.posttz()
	return m
}

func (m *Moment) MidOfDay() *Moment {
	m.pretz()
	ti := time.Date(m.t.Year(), m.t.Month(), m.t.Day(), 12, 0, 0, 0, m.t.Location())
	m.t = ti
	m.posttz()
	return m
}

func (m *Moment) EndOfDate() *Moment {
	m.pretz()
	ti := time.Date(m.t.Year(), m.t.Month(), m.t.Day(), 23, 59, 59, 999999999, m.t.Location())
	m.t = ti
	m.posttz()
	return m
}

func (m *Moment) StartOfMonth() *Moment {
	m.pretz()
	ti := time.Date(m.t.Year(), m.t.Month(), 1, 0, 0, 0, 0, m.t.Location())
	m.t = ti
	m.posttz()
	return m
}

func (m *Moment) EndOfMonth() *Moment {
	m.pretz()
	som_ti := time.Date(m.t.Year(), m.t.Month(), 1, 0, 0, 0, 0, m.t.Location())
	ti := som_ti.AddDate(0, 1, 0).Add(-time.Nanosecond)
	m.t = ti
	m.posttz()
	return m
}

func (m *Moment) StartOfYear() *Moment {
	m.pretz()
	ti := time.Date(m.t.Year(), 1, 1, 0, 0, 0, 0, m.t.Location())
	m.t = ti
	m.posttz()
	return m
}

func (m *Moment) EndOfYear() *Moment {
	m.pretz()
	som_ti := time.Date(m.t.Year(), 1, 1, 0, 0, 0, 0, m.t.Location())
	ti := som_ti.AddDate(1, 0, 0).Add(-time.Nanosecond)
	m.t = ti
	m.posttz()
	return m
}

func (m *Moment) Add(duration MomentDuration, amount int64) *Moment {
	switch duration {
	case Nanosecond:
		m.t = m.t.Add(time.Nanosecond * time.Duration(amount))
	case Microsecond:
		m.t = m.t.Add(time.Microsecond * time.Duration(amount))
	case Millisecond:
		m.t = m.t.Add(time.Millisecond * time.Duration(amount))
	case Second:
		m.t = m.t.Add(time.Second * time.Duration(amount))
	case Minute:
		m.t = m.t.Add(time.Minute * time.Duration(amount))
	case Hour:
		m.t = m.t.Add(time.Hour * time.Duration(amount))
	case Day:
		m.t = m.t.AddDate(0, 0, int(amount))
	case Month:
		m.t = m.t.AddDate(0, int(amount), 0)
	case Year:
		m.t = m.t.AddDate(int(amount), 0, 0)
	default:
		return m
	}

	return m
}

func (m *Moment) Diff(mo *Moment) time.Duration {
	return m.t.Sub(mo.t).Abs()
}

func (m *Moment) IsAfter(mo *Moment) bool {
	return m.t.After(mo.t)
}

func (m *Moment) IsAfterEqual(mo *Moment) bool {
	return m.t.Equal(mo.t) || m.t.After(mo.t)
}

func (m *Moment) IsBefore(mo *Moment) bool {
	return m.t.Before(mo.t) && !m.IsEqual(mo)
}

func (m *Moment) IsBeforeEqual(mo *Moment) bool {
	return m.t.Before(mo.t) || m.IsEqual(mo)
}

func (m *Moment) IsEqual(mo *Moment) bool {
	return m.t.Equal(mo.t)
}

func (m *Moment) IsBetween(from *Moment, to *Moment) bool {
	return m.IsAfterEqual(from) && m.IsBefore(to)
}

func (m *Moment) Format(layout string) string {
	layout = strings.ReplaceAll(layout, "A", TernaryOp(m.t.Hour()/12 >= 1, "PM", "AM"))
	layout = strings.ReplaceAll(layout, "hh", TernaryOp(m.t.Hour()%12 < 10, fmt.Sprintf("0%d", m.t.Hour()%12), fmt.Sprintf("%d", m.t.Hour()%12)))
	layout = strings.ReplaceAll(layout, "HH", TernaryOp(m.t.Hour() < 10, fmt.Sprintf("0%d", m.t.Hour()), fmt.Sprintf("%d", m.t.Hour())))
	layout = strings.ReplaceAll(layout, "mm", TernaryOp(m.t.Minute() < 10, fmt.Sprintf("0%d", m.t.Minute()), fmt.Sprintf("%d", m.t.Minute())))
	layout = strings.ReplaceAll(layout, "DD", TernaryOp(m.t.Day() < 10, fmt.Sprintf("0%d", m.t.Day()), fmt.Sprintf("%d", m.t.Day())))
	layout = strings.ReplaceAll(layout, "MM", TernaryOp(m.t.Month() < 10, fmt.Sprintf("0%d", m.t.Month()), fmt.Sprintf("%d", m.t.Month())))
	layout = strings.ReplaceAll(layout, "YYYY", TernaryOp(m.t.Year() < 10, fmt.Sprintf("0%d", m.t.Year()), fmt.Sprintf("%d", m.t.Year())))
	return layout
}

func (m *Moment) Weekday() int {
	m.pretz()
	return int(m.t.Weekday())
}

func (m *Moment) pretz() {
	if m.tz_offset != 0 {
		m.Add(Minute, -int64(m.tz_offset))
	}
}
func (m *Moment) posttz() {
	if m.tz_offset != 0 {
		m.Add(Minute, int64(m.tz_offset))
	}
}
