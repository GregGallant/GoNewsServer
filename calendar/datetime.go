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

// GetServerTimestamp returns current time (NOW) as timestamp
func GetServerTimestamp() int64 {
	tn := time.Now()
	timeDate := time.Date(tn.Year(), tn.Month(), tn.Day(), tn.Hour(), tn.Minute(), tn.Second(), tn.Nanosecond(), tn.Location())
	timeStamp := timeDate.UnixMilli()
	return timeStamp
}

// GetFutureTimestamp gets timestamp six hours ahead
func GetFutureTimestamp() int64 {
	tn := time.Now()
	timeDate := time.Date(tn.Year(), tn.Month(), tn.Day(), tn.Hour(), tn.Minute(), tn.Second(), tn.Nanosecond(), tn.Location())
	futureTime := addHoursToDate(timeDate, 24)
	timeStamp := futureTime.UnixMilli()
	return timeStamp
}

// GetPastTimestamp gets a timestamp from specified number of hours ago
func GetPastTimestamp(ghours int) int64 {
	tn := time.Now()
	timeDate := time.Date(tn.Year(), tn.Month(), tn.Day(), tn.Hour(), tn.Minute(), tn.Second(), tn.Nanosecond(), tn.Location())
	futureTime := removeHoursFromDate(timeDate, ghours)
	timeStamp := futureTime.UnixMilli()
	return timeStamp
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


// addHoursToDate adds specified number of hours to date
func addHoursToDate(cdate time.Time, ihours int) time.Time {
	durationStr := strconv.Itoa(ihours) + "h"
	ghours, err := time.ParseDuration(durationStr)
	if err != nil {
		fmt.Printf("%v %s", err, " :addHours")
	}
	nextNews := cdate.Add(ghours)
	return nextNews
}

// removeHoursFromDate removes specified hours to date
func removeHoursFromDate(cdate time.Time, ihours int) time.Time {
	durationStr := strconv.Itoa(ihours) + "h"
	ghours, err := time.ParseDuration(durationStr)
	if err != nil {
		fmt.Printf("%v %s", err, " :removeHours")
	}
	pastDate := cdate.Add(-(ghours))
	return pastDate
}

// DateIsExpired checks expired dates
func DateIsExpired(someDate string) bool {
	parsedDate := convertStringDate(someDate)
	forwardCheck := addHoursToDate(parsedDate, 24)
	date2Check := GetServerDate()

	return date2Check.After(forwardCheck)
}
