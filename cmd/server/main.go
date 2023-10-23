package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/gin-gonic/gin"
)

type Values struct {
	Id          string    `dynamodbav:"id"`
	Value       string    `dynamodbav:"value"`
	UpdateCount int       `dynamodbav:"updateCount"`
	LastUpdated time.Time `dynamodbav:"lastUpdated"`
}

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.LoadHTMLGlob("cmd/server/dist/index.html")
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.GET("/decay", func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		value := getDecayValue()
		c.JSON(http.StatusOK, gin.H{"value": value})
	})

	r.GET("/", func(c *gin.Context) {
		c.HTML((http.StatusOK), "index.html", nil)
	})
	return r
}

func main() {
	router := setupRouter()
	router.Static("/assets", "./cmd/server/dist/assets")
	router.Run(":8000")
}

func getDecayValue() (value string) {
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

	valueAsFloat, err := strconv.ParseFloat(values.Value, 64)

	if valueAsFloat < 1 || err != nil {
		return "0"
	}
	return values.Value
}
