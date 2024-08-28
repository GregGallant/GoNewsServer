package calendar

import (
	"fmt"
	"time"
)

// getServerTime generic Time method func Date(year int, month Month, day, hour, min, sec, nsec int, loc *Location) Time
func GetServerDate() time.Time {
	tn := time.Now()
	timeDate := time.Date(tn.Year(), tn.Month(), tn.Day(), tn.Hour(), tn.Minute(), tn.Second(), tn.Nanosecond(), tn.Location())
	return timeDate
}

func convertStringDate(fileDate string) time.Time {
	layout := "2006-01-02 15:04:05 +0000 UTC"
	parsedDate, err := time.Parse(layout, fileDate)
	if err != nil {
		fmt.Println(err)
	}
	return parsedDate
}

func add6HoursToDate(cdate time.Time) time.Time {
	nextNews := cdate.Add(time.Hour * 6)
	return nextNews
}

func DateIsExpired(someDate string) bool {
	parsedDate := convertStringDate(someDate)
	forwardCheck := add6HoursToDate(parsedDate)
	date2Check := GetServerDate()

	return date2Check.After(forwardCheck)
}

/*
From the playground:

func main() {
        //ourYear := GetServerDateTime()
	testYear := GetServerDateTime().String()

	layout := "2006-01-02 15:04:05 +0000 UTC"

	parsedTime, err := time.Parse(layout, testYear)

	if err != nil {
    		fmt.Println(err)
	}


	year2000 := time.Date(2010, 1, 1, 0, 0, 0, 0, time.UTC)
	//year3000 := time.Date(3000, 1, 1, 0, 0, 0, 0, time.UTC)

	isYear3000AfterYear2000 := parsedTime.After(year2000) // True
	isYear2000AfterYear3000 := year2000.After(parsedTime) // False

	fmt.Printf("serverYear(2009) is After(year2010) = %v\n", isYear3000AfterYear2000)
	fmt.Printf("year2010 is After(ourYear) = %v\n", isYear2000AfterYear3000)
	fmt.Printf("Our Year is: %v\n", parsedTime)
}

// GetServerDateTime
func GetServerDateTime() time.Time {
	tn := time.Now()
	timeDate := time.Date(tn.Year(), tn.Month(), tn.Day(), tn.Hour(), tn.Minute(), tn.Second(), tn.Nanosecond(), tn.Location())

	return timeDate

}

 */