package dailyMacroTotal

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	dynamoservice "github.com/ramenNoodles1998/macros-backend/internal/dynamo-service"
)

type DailyMacroTotalDB struct {
	//id
	PartitionKey string
	//date
	SortKey string
	Protein float64 
	Carbs float64
	Fat float64
}

type DailyMacroTotal struct {
	Id string  `json:"id"`
	Date string  `json:"type"`
	Protein float64  `json:"protein"`
	Carbs float64  `json:"carbs"`
	Fat float64  `json:"fat"`
}

const tableName string = "dev-macros"
const (
    YYYYMMDD = "20060102"
)
const uuidRoman string = "123123"
const dailyMacroTotalSuffix = "DAILY_MACRO_TOTAL-"

func SetMacroLogRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/save-daily-macro-total", saveDailyMacroTotal)
}

func saveDailyMacroTotal(w http.ResponseWriter, r *http.Request) {
	var dmt DailyMacroTotal

    err := json.NewDecoder(r.Body).Decode(&dmt)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

	//TODO: date problems keep track of front end to backedn times and compensate.
	//TODO: must get macro total first and add to it.
	var svc *dynamodb.DynamoDB = dynamoservice.DynamoService()
	if len(dmt.Date) == 0 {
		dmt.Date = time.Now().Format(YYYYMMDD);
	}

	var dmtdb = DailyMacroTotalDB {
		PartitionKey: uuidRoman,
		SortKey: dailyMacroTotalSuffix + dmt.Date,
		Protein: dmt.Protein,
		Carbs: dmt.Carbs,
		Fat: dmt.Fat,
	}

	rdmt := DailyMacroTotal{
		Id: uuidRoman,
		Date: dmt.Date,
		Protein: dmt.Protein,
		Carbs: dmt.Carbs,
		Fat: dmt.Fat,
	}

	av, err := dynamodbattribute.MarshalMap(dmtdb)
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

	json.NewEncoder(w).Encode(rdmt)
}