package dynamoservice

import (
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/dynamodb"
)

var sess *session.Session = session.Must(session.NewSessionWithOptions(session.Options{
	SharedConfigState: session.SharedConfigEnable,
}))
var svc *dynamodb.DynamoDB = dynamodb.New(sess)

func DynamoService() *dynamodb.DynamoDB {
	return svc 
}