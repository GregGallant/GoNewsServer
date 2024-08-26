package articles

import (
	"encoding/json"
	"errors"
	gcal "gallantone.com/main/calendar"
	"html"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strconv"
)

type Article struct {
	Title   string `json:"Title"`
	Desc    string `json:"desc"`
	Content string `json:"content"`
}

const filepath = "<yournewsdirpath>/news/"
const newsDateFile = "lastNewsUpdate.txt"
const newsFile = "genericNews.txt"

var Articles []Article

func PrepareNewsService() []byte {
	newsCheck()

	theNewsFile := filepath + newsFile
	newsfile, err := os.Open(theNewsFile)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err = newsfile.Close(); err != nil {
			log.Fatal(err)
		}
	}()
	theNews, err := ioutil.ReadAll(newsfile)
	if err != nil {
		log.Fatalf("unable to read date file: %v", err)
	}

	return theNews
}

func newsCheck() {

	logger := log.New(os.Stdout, "http: ", log.LstdFlags)
	lastNews := filepath + newsFile

	// Existance checks
	lastDate := filepath + newsDateFile

	logger.Println("Starting newsCheck()...")

	// Check date file
	//if _, err := os.Stat(lastDate); os.IsNotExist(err) {
	if existVar, err := os.Stat(lastDate); errors.Is(err, os.ErrNotExist) {
		if err != nil {
			logger.Println("Last news DATEFILE doesn't exist, should create here...")
			writeNewsDateFile()
		}

		fileSizeInfo := existVar.Size()
		logger.Println(" DATEFILE exists with stat filesize: " + strconv.FormatInt(fileSizeInfo, 10))
		//printLatestNews()
	}

	if _, err := os.Stat(lastNews); errors.Is(err, os.ErrNotExist) {
		if err != nil {
			logger.Println("generic NEWSFILE doesn't exist, should create here...")
			printLatestNews()
		}
	}

	logger.Println("Exist checks complete in newsCheck()...")

	// News and date file exists, check date expirations

	datefile, err := os.Open(lastDate)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err = datefile.Close(); err != nil {
			log.Fatal(err)
		}
	}()
	theDate, err := ioutil.ReadAll(datefile)
	if err != nil {
		log.Fatalf("unable to read date file: %v", err)
	}

	if theDate != nil {
		if gcal.DateIsExpired(string(theDate)) {
			writeNewsDateFile()
			printLatestNews()
		}
	}
}

// printLatestNews prints if files old or non-existant
func printLatestNews() {
	newsfile := filepath + newsFile

	oNews := InitWebhoseRequest()
	for i, eachPost := range oNews.Posts {

		// Emojis and Unicode
		var emojiRx = regexp.MustCompile(`[\x{10440}-\x{1F7FF}|[\x{1EF3}-\x{26FF}]|[\x{0021}-\x{0029}]`)
		textToEscape := emojiRx.ReplaceAllString(eachPost.Text, ``)

		// line breaks
		var lbreaks = regexp.MustCompile(`\n`)
		textToEscape = lbreaks.ReplaceAllString(textToEscape, `<br style="margin:0 0 20px 0;"/>`)

		urlToEscape := eachPost.URL
		for j, eachElink := range eachPost.ExternalLinks {
			oNews.Posts[i].ExternalLinks[j] = html.EscapeString(eachElink)
		}
		oNews.Posts[i].Text = textToEscape
		oNews.Posts[i].URL = html.EscapeString(urlToEscape)

	}


	newshose, err := json.Marshal(oNews)
	if err != nil {
		log.Fatal(err)
	}

	newshoseString := string(newshose)
	newsBytes := []byte(newshoseString)

	ferr := os.WriteFile(newsfile, newsBytes, 0667)
	if ferr != nil {
		log.Fatal(ferr)
	}
}

// writeNewsDateFile creates file holding the last time news was updated
func writeNewsDateFile() {
	logger := log.New(os.Stdout, "http: ", log.LstdFlags)
	newsDatefile := filepath + newsDateFile
	lastNewsDate := gcal.GetServerDate().String()

	ferr := os.WriteFile(newsDatefile, []byte(lastNewsDate), 0667)
	if ferr != nil {
		logger.Println("DATE FILE NOT CREATED...")
		log.Fatal(ferr)
	}
	logger.Println("DATE FILE SHOULD BE CREATED...")
}
