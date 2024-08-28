// Package calendar handles dates and times
package calendar

import (
	"fmt"
	"time"
)

// GetServerDate generic Time method func Date(year int, month Month, day, hour, min, sec, nsec int, loc *Location) Time
func GetServerDate() time.Time {
	tn := time.Now()
	timeDate := time.Date(tn.Year(), tn.Month(), tn.Day(), tn.Hour(), tn.Minute(), tn.Second(), tn.Nanosecond(), tn.Location())
	return timeDate
}

// convertStringDate converts date to string
func convertStringDate(fileDate string) time.Time {
	layout := "2006-01-02 15:04:05 +0000 UTC"
	parsedDate, err := time.Parse(layout, fileDate)
	if err != nil {
		fmt.Println(err)
	}
	return parsedDate
}

// add6HoursToDate adds six hours to the date
func add6HoursToDate(cdate time.Time) time.Time {
	nextNews := cdate.Add(time.Hour * 6)
	return nextNews
}

// DateIsExpired checks for expired date
func DateIsExpired(someDate string) bool {
	parsedDate := convertStringDate(someDate)
	forwardCheck := add6HoursToDate(parsedDate)
	date2Check := GetServerDate()

	return date2Check.After(forwardCheck)
}
