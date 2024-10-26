package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	_ "github.com/lib/pq"
)

type Result struct {
	Song   *Platform `json:"song"`
	Album  *Graph    `json:"album"`
	Artist *Graph    `json:"artist"`
}

type Platform struct {
	Melon *Song `json:"melon"`
	Genie *Song `json:"genie"`
	Flo   *Song `json:"flo"`
	Bugs  *Song `json:"bugs"`
	Vibe  *Song `json:"vibe"`
}

type Song struct {
	Name   string `json:"name"`
	Artist string `json:"artist"`
	Image  string `json:"image"`
}

type Graph struct {
	Label []string `json:"label"`
	Data  []int    `json:"data"`
}

func HandleLambdaEvent(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	result, err := fetchData()
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500}, err
	}

	responseBody, err := json.Marshal(result)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500}, err
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(responseBody),
	}, nil
}

func fetchData() (*Result, error) {
	db, err := connection()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	song, err := getTopSong(db)
	if err != nil {
		return nil, err
	}

	album, err := getTopItems(db, "album")
	if err != nil {
		return nil, err
	}

	artist, err := getTopItems(db, "artist")
	if err != nil {
		return nil, err
	}

	return &Result{
		Song:   song,
		Album:  album,
		Artist: artist,
	}, nil
}

func getTopSong(db *sql.DB) (*Platform, error) {
	rows, err := db.Query(`
		SELECT m.name, a.artist, a.image, r.types
		FROM music m 
		JOIN album a ON m.album_id = a.id
		JOIN ranks r ON r.music_id = m.id
		WHERE r.ranks = 1 
		ORDER BY r.crdate 
		LIMIT 5
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var platform Platform
	for rows.Next() {
		var name, artist, image, types string

		err := rows.Scan(&name, &artist, &image, &types)
		if err != nil {
			return nil, err
		}

		song := Song{name, artist, image}
		switch types {
		case "melon":
			platform.Melon = &song
		case "genie":
			platform.Genie = &song
		case "flo":
			platform.Flo = &song
		case "bugs":
			platform.Bugs = &song
		case "vibe":
			platform.Vibe = &song
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &platform, nil
}

func getTopItems(db *sql.DB, itemType string) (*Graph, error) {
	var col string
	if itemType == "album" {
		col = "a.name"
	} else if itemType == "artist" {
		col = "a.artist"
	}

	rows, err := db.Query(fmt.Sprintf(`
		SELECT %[1]s, COUNT(*)
		FROM music m
		JOIN album a ON m.album_id = a.id
		JOIN ranks r ON r.music_id = m.id
		WHERE r.crdate = (SELECT MAX(crdate) FROM ranks)
		GROUP BY %[1]s
		ORDER BY count DESC
		LIMIT 5
	`, col))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var label []string
	var data []int
	for rows.Next() {
		var name string
		var cnt int

		err := rows.Scan(&name, &cnt)
		if err != nil {
			return nil, err
		}

		label = append(label, name)
		data = append(data, cnt)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &Graph{label, data}, nil
}

func connection() (*sql.DB, error) {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s",
		os.Getenv("host"), os.Getenv("port"), os.Getenv("user"), os.Getenv("password"), os.Getenv("dbname"),
	)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func main() {
	lambda.Start(HandleLambdaEvent)
}
