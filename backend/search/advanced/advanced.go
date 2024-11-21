package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	_ "github.com/lib/pq"
)

type Body struct {
	Prompt string `json:"prompt"`
}

type ClaudeRequest struct {
	Version   string    `json:"anthropic_version"`
	MaxTokens int       `json:"max_tokens"`
	System    string    `json:"system"`
	Messages  []Message `json:"messages"`
}

type Message struct {
	Role     string    `json:"role"`
	Contents []Content `json:"content"`
}

type Content struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type ClaudeResponse struct {
	Contents []Content `json:"content"`
}

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
	var body Body
	requestBody := request.Body
	err := json.Unmarshal([]byte(requestBody), &body)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 400}, err
	}

	songs, err := advancedSearch(body.Prompt)
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

func advancedSearch(prompt string) (*[]Song, error) {
	system := `
	Here are the table schemas for the PostgreSQL database <postgre_schemas>.   
	<postgre_schemas>
		<postgre_schema>
			CREATE TABLE music (
				id SERIAL PRIMARY KEY, -- unique music id
				name VARCHAR(50) NOT NULL, -- name of the song
				album_id INTEGER NOT NULL, -- id of album table
				lyrics TEXT NOT NULL, lyrics of the song
				melon VARCHAR(15) UNIQUE, -- Song number on Melon Chart
				genie VARCHAR(15) UNIQUE, -- Song number on Genie Chart
				flo VARCHAR(15) UNIQUE, -- Song number on FLO Chart
				bugs VARCHAR(15) UNIQUE, -- Song number on Bugs Chart
				vibe VARCHAR(15) UNIQUE, -- Song number on Vibe Chart
				CONSTRAINT fk_album_id FOREIGN KEY(album_id) REFERENCES album(id) ON DELETE CASCADE ON UPDATE CASCADE
			); 
		</postgre_schema>
		<postgre_schema>
			CREATE TABLE album (
				id SERIAL PRIMARY KEY, -- unique album id
				name VARCHAR(50) NOT NULL, -- name of the album
				artist VARCHAR(50) NOT NULL, -- artist of the album
				image VARCHAR(200) NOT NULL -- Image link in album
			);
		</postgre_schema>
		<postgre_schema>
			CREATE TYPE chart_type AS ENUM ('melon', 'genie', 'flo', 'bugs', 'vibe');
			CREATE TABLE ranks (
				id SERIAL PRIMARY KEY, -- unique ranks id
				ranks smallint NOT NULL, -- music chart rankings
				crdate DATE NOT NULL, -- date of music chart
				types chart_type NOT NULL, -- Type of music chart
				music_id INTEGER NOT NULL, -- id of music table
				CONSTRAINT fk_music_id FOREIGN 	KEY(music_id) REFERENCES music(id) ON DELETE CASCADE ON UPDATE CASCADE
			);
		</postgre_schema>
	</postgre_schemas>

	Answer Format : SQL statement to query "name", "artist", "album", "image", "melon", "genie", "flo", "bugs", and "vibe" column.
	Answer Example : "SELECT m.name, a.artist, a.name album, a.image, m.melon, m.genie, m.flo, m.bugs, m.vibe FROM music m JOIN album a ON a.id=m.album_id;"
	Note : You must never return a description. You must only return queries so that they can be executed immediately.
	`

	body, err := json.Marshal(ClaudeRequest{
		Version:   "bedrock-2023-05-31",
		MaxTokens: 500,
		System:    system,
		Messages: []Message{
			{Role: "user", Contents: []Content{
				{Type: "text", Text: prompt},
			}},
			{Role: "assistant", Contents: []Content{
				{Type: "text", Text: "SELECT"},
			}},
		},
	})
	if err != nil {
		return nil, err
	}

	sdkConfig, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("ap-northeast-2"))
	if err != nil {
		return nil, err
	}

	bedrock := bedrockruntime.NewFromConfig(sdkConfig)
	output, err := bedrock.InvokeModel(context.TODO(), &bedrockruntime.InvokeModelInput{
		ModelId:     aws.String("anthropic.claude-3-sonnet-20240229-v1:0"),
		ContentType: aws.String("application/json"),
		Body:        body,
	})
	if err != nil {
		return nil, err
	}

	var response ClaudeResponse
	if err := json.Unmarshal(output.Body, &response); err != nil {
		return nil, err
	}

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s",
		os.Getenv("host"), os.Getenv("port"), os.Getenv("user"), os.Getenv("password"), os.Getenv("dbname"),
	)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	query := "SELECT " + response.Contents[0].Text
	fmt.Println(query)

	rows, err := db.Query(query)
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

	return &songs, nil
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
