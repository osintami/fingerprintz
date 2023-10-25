// Copyright © 2023 OSINTAMI. This is not yours.
package common

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDateStuf(t *testing.T) {
	stuff := NewDateStuff()
	today := time.Now()

	yesterday := today.AddDate(0, 0, -1)
	date := yesterday.Format("2006-01-02")
	dayAgo := stuff.AgoStringToDate("1.day.ago")
	assert.Equal(t, date, dayAgo)

	lastMonth := today.AddDate(0, -1, 0)
	date = lastMonth.Format("2006-01-02")
	monthAgo := stuff.AgoStringToDate("1.month.ago")
	assert.Equal(t, date, monthAgo)

	lastYear := today.AddDate(-1, 0, 0)
	date = lastYear.Format("2006-01-02")
	yearAgo := stuff.AgoStringToDate("1.year.ago")
	assert.Equal(t, date, yearAgo)
}

func TestGetDaysFromDate(t *testing.T) {
	today := time.Now()
	yesterday := today.AddDate(0, 0, -1)

	// expected date format
	days, err := GetDaysFromDate(yesterday.Format("2006-01-02"))
	assert.Nil(t, err)
	assert.Equal(t, 2, days)

	// date and time format
	days, err = GetDaysFromDate(yesterday.Format("2006-01-02 15:04:05"))
	assert.Nil(t, err)
	assert.Equal(t, 2, days)

	// universal time format
	days, err = GetDaysFromDate(yesterday.Format("2006-01-02T15:04:05Z7:00"))
	assert.Nil(t, err)
	assert.Equal(t, 2, days)

	// error paths
	days, err = GetDaysFromDate("xxx")
	assert.NotNil(t, err)
	assert.Equal(t, 0, days)

	days, err = GetDaysFromDate("")
	assert.Equal(t, ErrEmptyDate, err)
	assert.Equal(t, 0, days)
}

func TestTimeInHours(t *testing.T) {
	hourAgo := time.Now().Add((-time.Duration(time.Hour)))
	assert.Equal(t, 1, TimeInHours(hourAgo))
}
