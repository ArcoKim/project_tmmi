package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	_ "github.com/lib/pq"
)

type Song struct {
	Name   string  `json:"name"`
	Artist string  `json:"artist"`
	Album  string  `json:"album"`
	Image  string  `json:"image"`
	Melon  *string `json:"melon,omitempty"`
	Genie  *string `json:"genie,omitempty"`
	Flo    *string `json:"flo,omitempty"`
	Bugs   *string `json:"bugs,omitempty"`
	Vibe   *string `json:"vibe,omitempty"`
}

func HandleLambdaEvent(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	input, ok := request.QueryStringParameters["input"]
	if !ok {
		return events.APIGatewayProxyResponse{StatusCode: 400}, fmt.Errorf("missing query parameter")
	}

	songs, err := basicSearch(input)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500}, err
	}

	responseBody, err := json.Marshal(songs)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500}, err
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(responseBody),
	}, nil
}

func basicSearch(input string) ([]Song, error) {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s",
		os.Getenv("host"), os.Getenv("port"), os.Getenv("user"), os.Getenv("password"), os.Getenv("dbname"),
	)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	searchPattern := fmt.Sprintf("%%%s%%", strings.ToLower(input))

	rows, err := db.Query(`
		SELECT m.name, a.artist, a.name album, a.image, m.melon, m.genie, m.flo, m.bugs, m.vibe
		FROM music m 
		JOIN album a ON m.album_id = a.id
		WHERE LOWER(m.name) LIKE $1 OR LOWER(a.artist) LIKE $1 OR LOWER(a.name) LIKE $1
	`, searchPattern)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var songs []Song
	for rows.Next() {
		var song Song
		var melon, genie, flo, bugs, vibe sql.NullString

		err := rows.Scan(&song.Name, &song.Artist, &song.Album, &song.Image,
			&melon, &genie, &flo, &bugs, &vibe,
		)
		if err != nil {
			return nil, err
		}

		song.Melon = transformString(melon)
		song.Genie = transformString(genie)
		song.Flo = transformString(flo)
		song.Bugs = transformString(bugs)
		song.Vibe = transformString(vibe)

		songs = append(songs, song)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return songs, nil
}

func transformString(str sql.NullString) *string {
	if str.Valid {
		return &str.String
	}
	return nil
}

func main() {
	lambda.Start(HandleLambdaEvent)
}