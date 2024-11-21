package main

import (
	"context"
	"encoding/json"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/bedrockagentruntime"
	"github.com/aws/aws-sdk-go-v2/service/bedrockagentruntime/types"
)

type Body struct {
	Prompt string `json:"prompt"`
}

func HandleLambdaEvent(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var body Body
	requestBody := request.Body
	err := json.Unmarshal([]byte(requestBody), &body)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 400}, err
	}

	response, err := lyricsSearch(body)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500}, err
	}

	responseBody, err := json.Marshal(response)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500}, err
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(responseBody),
	}, nil
}

func lyricsSearch(body Body) (*string, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("ap-northeast-2"))
	if err != nil {
		return nil, err
	}
	bedrock := bedrockagentruntime.NewFromConfig(cfg)

	response, err := bedrock.RetrieveAndGenerate(context.TODO(), &bedrockagentruntime.RetrieveAndGenerateInput{
		Input: &types.RetrieveAndGenerateInput{Text: &body.Prompt},
		RetrieveAndGenerateConfiguration: &types.RetrieveAndGenerateConfiguration{
			Type: types.RetrieveAndGenerateType("KNOWLEDGE_BASE"),
			KnowledgeBaseConfiguration: &types.KnowledgeBaseRetrieveAndGenerateConfiguration{
				KnowledgeBaseId: aws.String(os.Getenv("kb_id")),
				ModelArn:        aws.String("arn:aws:bedrock:ap-northeast-2::foundation-model/anthropic.claude-3-sonnet-20240229-v1:0"),
			},
		},
	})
	if err != nil {
		return nil, err
	}

	return response.Output.Text, nil
}

func main() {
	lambda.Start(HandleLambdaEvent)
}
