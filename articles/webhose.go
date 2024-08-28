package articles

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)

const (
	whToken  = "<Your webz.io token here>"
	whDomain = "https://webhose.io/filterWebContent"
)

// News news and blog articles
type News struct {
	Posts            []Posts `json:"posts"`
	TotalResults     int     `json:"totalResults"`
	MoreResultsAvail int     `json:"moreResultsAvailable"`
	Next             string  `json:"next"`
	RequestsLeft     int     `json:"requestsLeft"`
	Warnings         int     `json:"warnings"`
}

// Posts articles fetched
type Posts struct {
	Thread               Thread           `json:"thread"`
	UUID                 string           `json:"uuid"`
	URL                  string           `json:"url"`
	OrdInThread          int              `json:"ord_in_thread"`
	ParentURL            string           `json:"parent_url"`
	Author               string           `json:"author"`
	Published            string           `json:"published"`
	Title                string           `json:"title"`
	Text                 string           `json:"text"`
	HighlightText        string           `json:"highlightText"`
	HighlightTitle       string           `json:"highlightTitle"`
	HighlightThreadTitle string           `json:"highlightThreadTitle"`
	Language             string           `json:"language"`
	ExternalLinks        []string         `json:"external_links"`
	ExternalImages       []ExternalImages `json:"external_images"`
	Entities             Entities         `json:"entities"`
	Rating               string           `json:"rating"`
	Crawled              string           `json:"crawled"`
	Updated              string           `json:"updated"`
}

// ExternalImages for external images
type ExternalImages struct {
	URL      string   `json:"url"`
	MetaInfo string   `json:"meta_info"`
	UUID     string   `json:"uuid"`
	Label    []string `json:"label"`
	Text     string   `json:"text"`
}

// Entities entities extracted
type Entities struct {
	Persons       []Entity
	Organizations []Entity
	Locations     []Entity
}

// Entity entity extracted
type Entity struct {
	Name      string `json:"name"`
	Sentiment string `json:"sentiment"`
}

// Thread details on thread
type Thread struct {
	UUID              string   `json:"uuid"`
	URL               string   `json:"url"`
	SiteFull          string   `json:"site_full"`
	Site              string   `json:"site"`
	SiteSection       string   `json:"site_section"`
	SiteCategories    []string `json:"site_categories"`
	SectionTitle      string   `json:"section_title"`
	Title             string   `json:"title"`
	TitleFull         string   `json:"title_full"`
	Published         string   `json:"published"`
	RepliesCount      int      `json:"replies_count"`
	ParticipantsCount int      `json:"participants_count"`
	SiteType          string   `json:"site_type"`
	Country           string   `json:"country"`
	SpamScore         float32  `json:"spam_score"`
	MainImage         string   `json:"main_image"`
	PerformanceScore  int      `json:"performance_score"`
	DomainRank        int      `json:"domain_rank"`
	Reach             Reach    `json:"reach"`
	Social            Social   `json:"social"`
}

// Reach controls viewers
type Reach struct {
	PerMillion float32   `json:"per_million"`
	PageViews  PageViews `json:"page_views"`
	Updated    string    `json:"updated"`
}

// PageViews SEO related
type PageViews struct {
	PerMillion float32 `json:"per_million"`
	PerUser    float32 `json:"per_user"`
}

// Social social interactions
type Social struct {
	Facebook    Facebook
	Gplus       Shares
	Pinterest   Shares
	LinkedIn    Shares
	StumbleUpon Shares
	Vk          Shares
}

// Facebook facebook data
type Facebook struct {
	Likes    int `json:"likes"`
	Comments int `json:"comments"`
	Shares   int `json:"shares"`
}

// Shares number of social shares
type Shares struct {
	Shares int `json:"shares"`
}

// InitWebhoseRequest initializes the webhoseAPI string
func InitWebhoseRequest() News {

	data := new(News)

	// Sample AI filter 
	qFilterString := "performance_score:>0 language:english -site:"usatoday.com" category:(Arts, Culture and Entertainment) (sentiment:\"positive\" AND site_type: \"news\" AND site_category: \"entertainment\" AND site_category: \"music\")" 
	
	// Prepare the Webhose Domain
	whNewsURL, werr := url.Parse(whDomain)
	if werr != nil {
		log.Fatal(werr)
	}

	// Set a timestamp via method if you want updated news
	nTimestamp := 1722473457731

	// Query String Parameters
	params := url.Values{}
	params.Add("token", whToken)
	params.Add("sort", "crawled")
	params.Add("format", "json")
	params.Add("q", qFilterString)
	params.Add("ts", strconv.Itoa(nTimestamp))
	whNewsURL.RawQuery = params.Encode()

	// call API
	cl := &http.Client{Timeout: 10 * time.Second}

	req, err := http.NewRequest("GET", whNewsURL.String(), nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header = http.Header{
		"Content-Type":   []string{"text/plain; charset=utf-8"},
		"Accept-Charset": []string{"utf-8"},
	}

	resp, arr := cl.Do(req)
	if arr != nil {
		log.Fatal(arr)
	}

	// defer happens once news is written
	defer resp.Body.Close()

	// Unmarshall
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		log.Fatal(err)
	}

	logger := log.New(os.Stdout, "http: ", log.LstdFlags)

	return *data
}
