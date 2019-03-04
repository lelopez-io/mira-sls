package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var (
	API_KEY      = os.Getenv("API_KEY")
	ErrorBackend = errors.New("Something went wrong")
)

// CorsHeaders for response
var CorsHeaders = map[string]string{
	"Access-Control-Allow-Origin":      "*",
	"Access-Control-Allow-Credentials": "true",
	"Access-Control-Allow-Headers":     "Content-Type,Authorization",
	"Access-Control-Allow-Methods":     "PUT, POST, GET, DELETE, OPTIONS, ANY",
}

type Request struct {
	ID int `json:"id"`
}

type MovieDBResponse struct {
	Movies []Movie `json:"results"`
}

type Movie struct {
	ID          int     `json:id`
	Title       string  `json:"title"`
	Description string  `json:"overview"`
	Cover       string  `json:"poster_path"`
	ReleaseDate string  `json:"release_date"`
	Runtime     int     `json:"runtime"`
	VoteAverage float64 `json:"vote_average"`
	Genres      []struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"genres"`
}

func Handler(request Request) (events.APIGatewayProxyResponse, error) {
	url := fmt.Sprintf("https://api.themoviedb.org/3/discover/movie?api_key=%s", API_KEY)

	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: "Error", StatusCode: 500}, nil
	}

	if request.ID > 0 {
		q := req.URL.Query()
		q.Add("with_genres", strconv.Itoa(request.ID))
		req.URL.RawQuery = q.Encode()
	}

	resp, err := client.Do(req)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: "Error", StatusCode: 500}, nil
	}
	defer resp.Body.Close()

	var data MovieDBResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return events.APIGatewayProxyResponse{Body: "Error", StatusCode: 500}, nil
	}

	jsonBody, _ := json.Marshal(data.Movies)
	stringBody := string(jsonBody) + "\n"

	return events.APIGatewayProxyResponse{Headers: CorsHeaders, Body: stringBody, StatusCode: 200}, nil

}

func main() {
	lambda.Start(Handler)
}
