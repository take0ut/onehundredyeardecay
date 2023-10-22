package main

import (
	"context"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type Values struct {
	Id          string    `dynamodbav:"id"`
	Value       string    `dynamodbav:"value"`
	UpdateCount int       `dynamodbav:"updateCount"`
	LastUpdated time.Time `dynamodbav:"lastUpdated"`
}

func main() {
	// TODO is an empty context which hasn't been evaluated yet.
	cfg, err := config.LoadDefaultConfig(context.TODO())

	if err != nil {
		log.Fatal(err)
	}

	values := Values{}

	tableName := os.Getenv("TABLE_NAME")

	client := dynamodb.NewFromConfig(cfg)

	pk, err := attributevalue.Marshal("decay")

	if err != nil {
		log.Fatal("Somehow failed to marshal the primary key.")
	}

	result, err := client.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]types.AttributeValue{
			"id": pk,
		},
	})

	log.Printf("Found values: %+v", result.Item)

	if (err != nil) || (result.Item == nil) {
		log.Fatal(err)
	}

	err = attributevalue.UnmarshalMap(result.Item, &values)

	if (err != nil) || (values.Id == "") {
		log.Fatal(err)
	}

	values.UpdateCount++

	values.LastUpdated = time.Now()

	valueAsFloat, err := strconv.ParseFloat(values.Value, 64)
	log.Printf("Found values: %+v", values)

	if err != nil {
		log.Fatal(err)
	}

	if valueAsFloat == 0 {
		log.Fatal("It has decayed.")
	} else {
		valueAsFloat -= 0.000001
	}

	values.Value = strconv.FormatFloat(valueAsFloat, 'f', -1, 64)

	item, err := attributevalue.MarshalMap(values)

	if err != nil {
		log.Fatal(err)
	}

	client.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: &tableName,
		Item:      item,
	})

	log.Printf("Updated values: %+v", values)

	os.Exit(0)

}
