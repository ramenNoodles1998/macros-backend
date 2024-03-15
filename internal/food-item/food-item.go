package fooditem

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	dynamoservice "github.com/ramenNoodles1998/macros-backend/internal/dynamo-service"
)

type FoodItemDB struct {
	//uuid
	PartitionKey string
	//name
	SortKey string
	Protein float64
	Carbs   float64
	Fat     float64
	Serving string
}

type FoodItem struct {
	Id      string  `json:"id"`
	Name    string  `json:"name"`
	Protein float64  `json:"protein"`
	Carbs   float64  `json:"carbs"`
	Fat     float64  `json:"fat"`
	Serving string  `json:"serving"`
}

const tableName string = "dev-macros"
const uuidRoman string = "123123"
const namePrefix = "NAME-"

func SetFoodItemRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/save-food-item", saveFoodItem)
	mux.HandleFunc("/api/get-food-items", getFoodItems)
	mux.HandleFunc("/api/get-food-item-by-id", getFoodItemId)
	mux.HandleFunc("/api/delete-food-item", deleteFoodItem)
}

func saveFoodItem(w http.ResponseWriter, r *http.Request) {
	var fi FoodItem 

    err := json.NewDecoder(r.Body).Decode(&fi)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

	foodItem := FoodItemDB{
		PartitionKey: fi.Id,
		SortKey: namePrefix + fi.Name,
		Protein: fi.Protein,
		Carbs: fi.Carbs,
		Fat: fi.Fat,
		Serving: fi.Serving,
	}

	returnFoodItem := FoodItem{
		Id: fi.Id,
		Name: fi.Name,
		Protein: fi.Protein,
		Carbs: fi.Carbs,
		Fat: fi.Fat,
		Serving: fi.Serving,
	}

	var svc *dynamodb.DynamoDB = dynamoservice.DynamoService()

	av, err := dynamodbattribute.MarshalMap(foodItem)
	if err != nil {
		fmt.Printf("Got error marshalling food item: %s", err)
		return
	}

	input := &dynamodb.PutItemInput{
		Item: av,
		TableName: aws.String(tableName),
	}

	_, err = svc.PutItem(input)

	if err != nil {
		fmt.Printf("Got error calling PutItem: %s", err)
		return
	}

	json.NewEncoder(w).Encode(returnFoodItem)
}

func getFoodItems(w http.ResponseWriter, r *http.Request) {
	var svc *dynamodb.DynamoDB = dynamoservice.DynamoService()
	result, err := svc.Query(&dynamodb.QueryInput{
		TableName: aws.String(tableName),
		KeyConditions: map[string]*dynamodb.Condition{
			"PartitionKey": {
				ComparisonOperator: aws.String("EQ"),
				AttributeValueList: []*dynamodb.AttributeValue {
					{
						S: aws.String(uuidRoman),
					},
				},
			},
			"SortKey": {
				ComparisonOperator: aws.String("BEGINS_WITH"),
				AttributeValueList: []*dynamodb.AttributeValue {
					{
						S: aws.String("NAME"),
					},
				},
			},
			
		},
	})

	if err != nil {
		fmt.Printf("Got error calling GetItem: %s", err)
		return
	}

	if result.Items == nil || len(result.Items) == 0 {
		fmt.Printf("Could not find Logs")
		return
	}
    
	foodItems := []FoodItem{}

	for _, item := range result.Items {
		fiDB := FoodItemDB{}
		err = dynamodbattribute.UnmarshalMap(item, &fiDB)
		var fi = FoodItem{
			Id: fiDB.PartitionKey,
			Protein: fiDB.Protein,
			Carbs: fiDB.Carbs,
			Fat: fiDB.Fat,
			Serving: fiDB.Serving,
		}

		fi.Name, _ = strings.CutPrefix(fiDB.SortKey, namePrefix)
		foodItems = append(foodItems, fi)
		if err != nil {
			panic(fmt.Sprintf("Failed to unmarshal Record, %v", err))
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(foodItems)
}

func getFoodItemId(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	var svc *dynamodb.DynamoDB = dynamoservice.DynamoService()
	result, err := svc.Query(&dynamodb.QueryInput{
		TableName: aws.String(tableName),
		Limit: aws.Int64(1),
		KeyConditions: map[string]*dynamodb.Condition{
			"PartitionKey": {
				ComparisonOperator: aws.String("EQ"),
				AttributeValueList: []*dynamodb.AttributeValue {
					{
						S: aws.String(uuidRoman),
					},
				},
			},
			"SortKey": {
				ComparisonOperator: aws.String("EQ"),
				AttributeValueList: []*dynamodb.AttributeValue {
					{
						S: aws.String(namePrefix + id),
					},
				},
			},
			
		},
	})

	if err != nil {
		fmt.Printf("Got error calling GetItem: %s", err)
		return
	}

	if result.Items == nil || len(result.Items) == 0 {
		fmt.Printf("Could not find Logs")
		return
	}
    
	foodItems := []FoodItem{}

	for _, item := range result.Items {
		fiDB := FoodItemDB{}
		err = dynamodbattribute.UnmarshalMap(item, &fiDB)
		var fi = FoodItem{
			Id: fiDB.PartitionKey,
			Protein: fiDB.Protein,
			Carbs: fiDB.Carbs,
			Fat: fiDB.Fat,
			Serving: fiDB.Serving,
		}

		fi.Name, _ = strings.CutPrefix(fiDB.SortKey, namePrefix)
		foodItems = append(foodItems, fi)
		if err != nil {
			panic(fmt.Sprintf("Failed to unmarshal Record, %v", err))
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(foodItems[0])
}

func deleteFoodItem(w http.ResponseWriter, r *http.Request) {
	var fi FoodItem 

    err := json.NewDecoder(r.Body).Decode(&fi)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

	var svc *dynamodb.DynamoDB = dynamoservice.DynamoService()

	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"PartitionKey": {
				S: aws.String(fi.Id),
			},
			"SortKey": {
				S: aws.String(namePrefix + fi.Name),
			},
		},
		TableName: aws.String(tableName),
	}

	_, err = svc.DeleteItem(input)
	if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
	}

	json.NewEncoder(w).Encode(fi)
}
