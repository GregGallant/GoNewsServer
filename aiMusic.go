package music

import (
	"context"
	"errors"
	"fmt"
	"github.com/predictionguard/go-client"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"strings"
	"time"

	md "github.com/JohannesKaufmann/html-to-markdown"
)

// VectorizedChunk is the struct holding the vectorized chunk
type VectorizedChunk struct {
	ID       int       `json:"id"`
	Chunk    string    `json:"chunk"`
	Vector   []float64 `json:"vector"`
	Metadata string    `json:"metadata"`
}

var nvidiaAPIKey = "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"

var embeddingModel = "bge-m3"

var host = "https://api.predictionguard.com"
var apiKey = os.Getenv("PGKEY") // get non-test key

// VectorizedChunks is an array of vectorizedChunks
type VectorizedChunks []VectorizedChunk

type eInputs []client.EmbeddingInput

// Log all this stuff to see what it is actually doing
var datalog *log.Logger

// qAPromptTemplate is a basic qa string template
func qAPromptTemplate(context, question string) string {
	return fmt.Sprintf(`Context: "%s" Question: "%s"`, context, question)
}

// AIStart starts somewhere with ai (embed function)
func AIStart() {

	datalog = log.New(os.Stdout, "http: ", log.LstdFlags)
	datalog.Println("AI Start called")

	// Test website
	website := "https://en.wikipedia.org/wiki/Aretha_Franklin"

	log.Println("AI Start website: " + website)

	chunks, err := websiteChunks(website, "", "")
	if err != nil {
		log.Fatal(err) // TODO? Fail gracefully instead
	}

	vectorizedChunks := VectorizedChunks{}
	for i, chunk := range chunks {
		fmt.Printf("Embedding chunk %d of %d\n", i+1, len(chunks))

		vectorizedChunk, err := embed("", chunk)
		if err != nil {
			log.Fatal(err) // TODO? Fail gracefully instead
		}

		vectorizedChunks = append(vectorizedChunks, *vectorizedChunk)
		vectorizedChunks[i].ID = i
		vectorizedChunks[i].Metadata = chunk
	}

	// Begin test prompt
	input := "Who is Aretha Franklin"

	// Embed the input question
	embedding, err := embed("", input)
	if err != nil {
		log.Fatal(err)
	}

	// Search...
	chunk, err := search(vectorizedChunks, *embedding)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Vectorized Chunks search called")

	// Response
	if err := run(input, string(chunk)); err != nil {
		log.Fatalln(err)
	}
}

// run handles response output
func run(query, queryContext string) error {

	logger := func(ctx context.Context, msg string, v ...any) {
		datalog.Printf("Using context: %v", ctx)
		s := fmt.Sprintf("run msg: %s", msg)
		for i := 0; i < len(v); i = i + 2 {
			s = s + fmt.Sprintf(", %s: %v", v[i], v[i+1])
		}
	}

	cln := client.New(logger, host, apiKey)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	input := client.ChatSSEInput{
		Model: "bge-m3",
		Messages: []client.ChatInputMessage{
			{
				Role:    client.Roles.System,
				Content: "Read the context provided...",
			},
			{
				Role: client.Roles.User,
				Content: qAPromptTemplate(
					string(queryContext),
					query,
				),
			},
		},
		MaxTokens:   client.Ptr(1000),
		Temperature: client.Ptr[float32](0.3),
		TopP:        client.Ptr(0.1),
	}

	ch := make(chan client.ChatSSE, 1000)

	err := cln.ChatSSE(ctx, input, ch)
	if err != nil {
		return fmt.Errorf("error on channel: %w", err)
	}

	for resp := range ch {
		for _, choice := range resp.Choices {
			fmt.Println(choice.Delta.Content)
		}
	}

	return nil
}

// search searches the vectorized chunks
func search(chunks VectorizedChunks, embedding VectorizedChunk) (string, error) {
	outChunk := ""

	var maxSimilarity float64

	for _, c := range chunks {
		distance, err := cosineSimilarity(c.Vector, embedding.Vector)
		if err != nil {
			return "", err
		}
		if distance > maxSimilarity {
			outChunk = c.Chunk
			maxSimilarity = distance
		}
	}

	return outChunk, nil
}

// cosineSimilarity calculates the cosine similarity between two vectors
func cosineSimilarity(a []float64, b []float64) (cosine float64, err error) {

	count := 0

	lengthA := len(a)
	lengthB := len(b)

	if lengthA > lengthB {
		count = lengthA
	} else {
		count = lengthB
	}

	sumA := 0.0
	s1 := 0.0
	s2 := 0.0
	for k := 0; k < count; k++ {
		if k >= lengthA {
			//s2 += math.Pow(b[k], 2)
			s2 += b[k] * b[k]
			continue
		}
		if k >= lengthB {
			//s1 += math.Pow(a[k], 2)
			s1 += a[k] * a[k]
			continue
		}
		sumA += a[k] * b[k]
		//s1 += math.Pow(a[k], 2)
		s1 += a[k] * a[k]
		//s2 += math.Pow(b[k], 2)
		s2 += b[k] * b[k]

	}

	// No null vectors
	if s1 == 0 || s2 == 0 {
		return 0.0, errors.New("vectors should not be null (all zeros)")
	}

	// Final result
	cosineSim := sumA / math.Sqrt(s2)

	return cosineSim, nil

}

func embed(imageLink string, text string) (*VectorizedChunk, error) {
	logger := func(ctx context.Context, msg string, v ...any) {
		s := fmt.Sprintf("embed msg: %s", msg)

		log.Printf("%s : %v : %v", s, ctx, v)
	}
	fmt.Println("First AI")
	cln := client.New(logger, host, apiKey)

	// Loading in a context to be passed into the pg client embedding, but for what?
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Images
	var image client.ImageNetwork
	if imageLink != "" {
		imageParsed, err := client.NewImageNetwork(imageLink)
		if err != nil {
			return nil, fmt.Errorf("error: %w", err)
		}
		image = imageParsed
	}

	//	var embeddingInputs eInputs
	/*
		embeddingInputs = []client.EmbeddingInput{
			{
				Text:  text,
				Image: image,
			},
		}
	*/
	embeddingInputs := eInputs{
		{
			Text:  text,
			Image: image,
		},
	}

	if imageLink != "" {
		embeddingInputs[0].Image = image
	}

	// New prediction guard way of calling embeddings
	//var input client.EmbeddingInputTypes
	input := embeddingInputs
	input.EmbedInputType()

	// This takes three arguments, the context, the model and the input - we need to specify a model if necessary
	resp, err := cln.Embedding(ctx, embeddingModel, input)
	if err != nil {
		return nil, fmt.Errorf("error: %w", err)
	}

	return &VectorizedChunk{
		Chunk:  text,
		Vector: resp.Data[0].Embedding,
	}, nil

}

// EmbedInputType exists to satisfy the interface
func (i eInputs) EmbedInputType() {}

// websiteChunks loads the website and splits it into chunks - start string and end string are optional
func websiteChunks(website string, start string, end string) ([]string, error) {

	converter := md.NewConverter("", true, nil)

	// get the website
	res, err := http.Get(website)
	if err != nil {
		return nil, err
	}

	content, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	cErr := res.Body.Close()
	if cErr != nil {
		return nil, err
	}

	html := string(content)

	// Convert html to markdown
	markdown, err := converter.ConvertString(html)
	if err != nil {
		return nil, err
	}

	// Split markdown on any provided start / end string
	if start != "" {
		markdownRemaining := strings.Split(markdown, start)[1:]
		markdown = strings.Join(markdownRemaining, "")
	}
	if end != "" {
		markdown = strings.Split(markdown, end)[0]
	}

	// Split the text into chunk sizes w/ overlap
	chunks := characterTextSplitter(markdown, 100, 10)
	return chunks, nil
}

// characterTextSplitter takes a string and splits it into chunks of given size, w/ overlap
func characterTextSplitter(text string, splitSize int, overlapSize int) []string {

	// Create a slice to hold chunks
	chunks := []string{}

	// Split the text into tokens via whitespace
	tokens := strings.Split(text, " ")

	// Loop over the tokens creating chunks of splitSize with an overlap of overlapSize
	for i := 0; i < len(tokens); i += splitSize - overlapSize {
		end := i + splitSize - overlapSize
		if end > len(tokens) {
			end = len(tokens)
		}
		chunks = append(chunks, strings.Join(tokens[i:end], " "))
	}

	return chunks
}
