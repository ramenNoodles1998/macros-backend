package macrolog

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	dynamoservice "github.com/ramenNoodles1998/macros-backend/internal/dynamo-service"
)

type MacroLogDB struct {
	//uuid
	PartitionKey string
	//date
	SortKey string
	Protein float64 
	Carbs float64
	Fat float64
}

type MacroLog struct {
	Id string
	Date string
	Protein float64 
	Carbs float64
	Fat float64
}

type Macro struct {
	Protein float64 
	Carbs float64 
	Fat float64 
}

const tableName string = "dev-macros"
const (
    YYYYMMDD = "2006-01-02"
)
const uuidRoman string = "123123"

func SaveMacroLog(w http.ResponseWriter, r *http.Request) {
	var m Macro 

    err := json.NewDecoder(r.Body).Decode(&m)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

	var svc *dynamodb.DynamoDB = dynamoservice.DynamoService()

	var macroLog = MacroLogDB {
		PartitionKey: uuidRoman,
		SortKey: time.Now().String(),
		Protein: m.Protein,
		Carbs: m.Carbs,
		Fat: m.Fat,
	}

	av, err := dynamodbattribute.MarshalMap(macroLog)
	if err != nil {
		fmt.Printf("Got error marshalling item: %s", err)
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
}


func GetMacroLogId(w http.ResponseWriter, r *http.Request) {
	//TODO: sortkey need Date-id
	id := r.URL.Query().Get("id")
	id = strings.Join(strings.Split(id, "%"), " ")
	fmt.Printf("%s", id)
	var svc *dynamodb.DynamoDB = dynamoservice.DynamoService()
	//gets todays log
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
				ComparisonOperator: aws.String("EQ"),
				AttributeValueList: []*dynamodb.AttributeValue {
					{
						S: aws.String("2024-02-29 16:36:03.6702179 -0700 MST m=+10.632405301"),
					},
				},
			},
		},
	})

	if err != nil {
		fmt.Printf("Got error calling GetItem: %s", err)
		return
	}

	if result.Items == nil {
		fmt.Printf("Could not find Logs")
		return
	}
   fmt.Println("here 2 %d", len(result.Items)) 
	logs := []MacroLog{}

	for _, item := range result.Items {
		logDB := MacroLogDB{}
		err = dynamodbattribute.UnmarshalMap(item, &logDB)
		var log = MacroLog {
			Id: logDB.PartitionKey,
			Date: logDB.SortKey,
			Protein: logDB.Protein,
			Carbs: logDB.Carbs,
			Fat: logDB.Fat,
		}

		logs = append(logs, log)
		if err != nil {
			panic(fmt.Sprintf("Failed to unmarshal Record, %v", err))
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(logs[0])
}


func GetMacroLog(w http.ResponseWriter, r *http.Request) {
	var svc *dynamodb.DynamoDB = dynamoservice.DynamoService()
	//gets todays log
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
						S: aws.String(time.Now().Format(YYYYMMDD)),
					},
				},
			},
		},
	})

	if err != nil {
		fmt.Printf("Got error calling GetItem: %s", err)
		return
	}

	if result.Items == nil {
		fmt.Printf("Could not find Logs")
		return
	}
    
	logs := []MacroLog{}

	for _, item := range result.Items {
		logDB := MacroLogDB{}
		err = dynamodbattribute.UnmarshalMap(item, &logDB)
		var log = MacroLog {
			Id: logDB.PartitionKey,
			Date: logDB.SortKey,
			Protein: logDB.Protein,
			Carbs: logDB.Carbs,
			Fat: logDB.Fat,
		}

		logs = append(logs, log)
		if err != nil {
			panic(fmt.Sprintf("Failed to unmarshal Record, %v", err))
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(logs)
}


func DeleteMacroLog(w http.ResponseWriter, r *http.Request) {
	var m MacroLog 

    err := json.NewDecoder(r.Body).Decode(&m)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

	fmt.Printf("%p+", m)

	var svc *dynamodb.DynamoDB = dynamoservice.DynamoService()

	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"PartitionKey": {
				S: aws.String(m.Id),
			},
			"SortKey": {
				S: aws.String(m.Date),
			},
		},
		TableName: aws.String(tableName),
	}

	_, err = svc.DeleteItem(input)
	if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
	}
}