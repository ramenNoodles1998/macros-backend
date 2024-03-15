package dailyMacroTotal

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
	Date string  `json:"date"`
	Protein float64  `json:"protein"`
	Carbs float64  `json:"carbs"`
	Fat float64  `json:"fat"`
}

const tableName string = "dev-macros"
const (
    YYYYMMDD = "20060102"
)
const uuidRoman string = "123123"
const dailyMacroTotalPrefix = "DAILY_MACRO_TOTAL-"

func SetMacroLogRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/save-daily-macro-total", saveDailyMacroTotal)
	mux.HandleFunc("/api/get-daily-macro-total", getDailyMacroTotal)
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
		SortKey: dailyMacroTotalPrefix + dmt.Date,
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

func getDailyMacroTotal(w http.ResponseWriter, r *http.Request) {
	var svc *dynamodb.DynamoDB = dynamoservice.DynamoService()
	//TODO: here we just add up logs
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
						S: aws.String(dailyMacroTotalPrefix + time.Now().Format(YYYYMMDD)),
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
		fmt.Printf("Could not find Daily Macro Total")
		json.NewEncoder(w).Encode(DailyMacroTotal{ Id: "", Date: time.Now().Format(YYYYMMDD), Protein: 0, Carbs: 0, Fat: 0 })
		return
	}
    
	dmts := []DailyMacroTotal{}

	for _, item := range result.Items {
		dmtdb := DailyMacroTotalDB{}
		err = dynamodbattribute.UnmarshalMap(item, &dmtdb)
		var dmt = DailyMacroTotal{
			Id: dmtdb.PartitionKey,
			Date: dmtdb.SortKey,
			Protein: dmtdb.Protein,
			Carbs: dmtdb.Carbs,
			Fat: dmtdb.Fat,
		}

		dmt.Date, _ = strings.CutPrefix(dmt.Date, dailyMacroTotalPrefix)
		dmts = append(dmts, dmt)
		if err != nil {
			panic(fmt.Sprintf("Failed to unmarshal Record, %v", err))
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dmts[0])

}