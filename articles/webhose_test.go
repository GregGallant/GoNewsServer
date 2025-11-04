package articles

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestInitWebhoseRequest tests basic news functionality
func TestInitWebhoseRequest(t *testing.T) {

	newsType := "default"

	result := InitWebhoseRequest(newsType)

	assert.NotNil(t, result, "News is available")
	assert.IsType(t, result, News{}, "Returned value is of type News")

}
