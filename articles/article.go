// Package articles for news
package articles

import (
	"encoding/json"
	"errors"
	gcal "gallantone.com/main/calendar"
	"html"
	"io"
	"log"
	"os"
	"regexp"
	"slices"
	"strings"
)


// Article is each base news data
type Article struct {
	Title   string `json:"Title"`
	Desc    string `json:"desc"`
	Content string `json:"content"`
}

// Original filepath
const filepath = "<YourfilepathHere>"

// Article title compare length
const defaultLen = 25

const newsDateFile = "<YourfilenameHere>"
const newsFile = "<YourfilenameHere>"

// Articles is the base map of articles
var Articles []Article

// PrepareNewsService prepares the writes
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
	theNews, err := io.ReadAll(newsfile)
	if err != nil {
		log.Printf("unable to read date file: %v", err)
	}

	return theNews
}

// newsCheck checks for updated news
func newsCheck() {

	logger := log.New(os.Stdout, "http: ", log.LstdFlags)
	lastNews := filepath + newsFile

	// Existance checks
	lastDate := filepath + newsDateFile

	logger.Println("Starting newsCheck()...")

	// Check date file
	if existVar, err := os.Stat(lastDate); errors.Is(err, os.ErrNotExist) {
		if err != nil {
			logger.Println("Last news DATEFILE doesn't exist, should create here...")
			writeNewsDateFile()
		}

		// As if...
		logger.Printf(" Something, something, crash, ugh: %v", existVar)

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
	theDate, err := io.ReadAll(datefile)
	if err != nil {
		log.Fatalf("unable to read date file: %v", err)
	}

	logger.Println("Attempting to check if date is expired")

	if theDate != nil {
		if gcal.DateIsExpired(string(theDate)) {
			logger.Println("Date is expired, attempting to print new news")
			writeNewsDateFile()
			printLatestNews()
		}
	}
}

// printLatestNews prints if files old or non-existant
func printLatestNews() {

	//logger := log.New(os.Stdout, "http: ", log.LstdFlags)

	newsfile := filepath + newsFile

	oNews := InitWebhoseRequest()

	for i, eachPost := range oNews.Posts {

		// Since the range changes, handle the updated range
		if i == len(oNews.Posts) {
			break
		}

		// Emojis and Unicode (all of this will change within interface soon...)
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

		r := regexp.MustCompile(`[\x{FFFD}]+`)
		oNews.Posts[i].Thread.Title = r.ReplaceAllString(oNews.Posts[i].Thread.Title, ``)


		oNews.Posts[i].Thread.MainImage = oNews.cleanImageURL(i)

		// Create rune of string
		runeTitle := []rune(oNews.Posts[i].Title)

		// Run through array, remove duplicates based on constant, maybe use ceil compare or assure certain title length
		// ex: lenComp = lengthTitle - int(math.Ceil(float64(lengthTitle)/3))
		looseCompare := string(runeTitle[:defaultLen])

		postsToDelete := oNews.createDelArr(i, looseCompare)

		oNews.Posts = oNews.delPosts(postsToDelete)

		//logger.Printf("Posts to delete: %v", postsToDelete)
		/*
			for j := len(postsToDelete); j > 0; j-- {
				logger.Printf("Posts to delete: %v", postsToDelete[j-1])
				oNews.Posts = slices.Delete(oNews.Posts, postsToDelete[j-1], postsToDelete[j-1]+1)
			}
		*/
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

// still want oNews as result... reorg as interface
//func organizeNews(nm NewsModel) {}
//func cleanPostTitle(oNews News, i int) string { }

func (n News) cleanImageURL(i int) string {

	ampReg := regexp.MustCompile(`&amp;`)
	n.Posts[i].Thread.MainImage = ampReg.ReplaceAllString(n.Posts[i].Thread.MainImage, `&`)

	return n.Posts[i].Thread.MainImage
}

// delPosts deletes posts from news Posts array
func (n News) delPosts(p []int) []Posts {
	for j := len(p); j > 0; j-- {
		//logger.Printf("Posts to delete: %v", p[j-1])
		n.Posts = slices.Delete(n.Posts, p[j-1], p[j-1]+1)
	}

	return n.Posts
}

// createDelArr creates a deletion array using articles, i count, lc as compared article title
// Todo: find common duplicate words within title and count for deletion
func (n News) createDelArr(i int, lc string) []int {

	var postsToDelete []int

	for m, k := range n.Posts {
		if m != i {
			looseRune := []rune(k.Title)
			lastLooseCompare := string(looseRune[:defaultLen])
			//logger.Println("Comparing strings at " + strconv.Itoa(m) + "|" + strconv.Itoa(i) + " with: " + looseCompare + " vs. " + lastLooseCompare)
			if strings.Compare(lc, lastLooseCompare) == 0 {
				postsToDelete = append(postsToDelete, m)
				//oNews.Posts = slices.Delete(oNews.Posts, m, m+1)
			}
		}
	}

	return postsToDelete
}

//}

// writeNewsDateFile creates file holding the last time news was updated
func writeNewsDateFile() {
	logger := log.New(os.Stdout, "http: ", log.LstdFlags)
	newsDatefile := filepath + newsDateFile
	lastNewsDate := gcal.GetServerDate().String()

	logger.Println("curerent date format: " + lastNewsDate)

	// Check if directory exists, if not create it
	if _, err := os.Stat(filepath); errors.Is(err, os.ErrNotExist) {
		if err != nil {
			dirErr := os.Mkdir(filepath, 0750)
			if dirErr != nil && !os.IsExist(err) {
				logger.Println("wtf no dir")
				log.Fatal(dirErr)
			}
		}
	}

	ferr := os.WriteFile(newsDatefile, []byte(lastNewsDate), 0667)
	if ferr != nil {
		logger.Println("DATE FILE NOT CREATED...")
		log.Fatal(ferr)
	}
	logger.Println("DATE FILE CREATED...")
}
