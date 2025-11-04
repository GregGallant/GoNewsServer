package articles

import (
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
)

type MockNews struct {
	mock.Mock
}

// GetNews creates mock news
func (m *MockNews) GetNews() News {

	n := News{
		Posts: []Posts{
			{
				URL:  "https://www.blah<&>.com",
				Text: "Blah blah$&%$%&#( Blah",
			},
		},
	}

	return n
}

// TestCleanUrls ensures urls have un-allowed, non-unicode special characters removed
func TestCleanUrls(t *testing.T) {

	testIndex := 0

	mNews := new(MockNews)

	news := mNews.GetNews()

	require.NotNil(t, news.Posts, "Posts should not be nil")
	require.NotNil(t, news, "News should at least be type news and not nil")
	require.IsType(t, news, News{}, "news should be type News")

	result := news.cleanUrls(testIndex)

	expected := "https://www.blah&lt;&amp;&gt;.com"

	require.Equal(t, expected, result, "URL should match")

}

// TestCleanArText cleans article text of un-allowed non-unicode text and some wacky emojis
func TestCleanArText(t *testing.T) {

	mNews := new(MockNews)

	news := mNews.GetNews()

	require.NotNil(t, news.Posts, "Posts should not be nil")
	require.NotNil(t, news, "News should at least be type news and not nil")
	require.IsType(t, news, News{}, "news should be type News")

	result := news.cleanArText(0)

	expected := "Blah blah Blah"

	require.Equal(t, expected, result, "Article text should match")

}
