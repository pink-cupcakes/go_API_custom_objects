package workers

import (
    "fmt"
    "encoding/json"
    "os"

    "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

func worker(id int, q struct{Jobs chan []byte; Results chan []byte; svc *dynamodb.DynamoDB}) {
    fmt.Println("Started worker ", id)
    for j := range q.Jobs {
        var request map[string]interface{}
        fmt.Println("worker", id, "started  job", j)
        json.Unmarshal([]byte(j), &request)
        av, err := dynamodbattribute.MarshalMap(request)
        if err != nil {
            fmt.Println("Got error marshalling new movie item:")
            fmt.Println(err.Error())
            os.Exit(1)
        }
        tableName := "CustomFields"
        input := &dynamodb.PutItemInput{
            Item:      av,
            TableName: aws.String(tableName),
        }
        _, err = q.svc.PutItem(input)
        if err != nil {
            fmt.Println(input)
            fmt.Println("Got error calling PutItem:")
            fmt.Println(err.Error())
            os.Exit(1)
        }
        fmt.Println("worker", id, "finished job", j)
        q.Results <- j
    }
}

func Start() struct{Jobs chan []byte; Results chan []byte; svc *dynamodb.DynamoDB} {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Create DynamoDB client
	svc := dynamodb.New(sess, &aws.Config{Endpoint: aws.String("http://localhost:8000")})
    type Queues struct {
        Jobs    chan []byte
        Results chan []byte
        svc 	*dynamodb.DynamoDB

    }
    q := Queues{make(chan []byte, 10000), make(chan []byte, 10000), svc}
    for w := 1; w <= 500; w++ {
        go worker(w, q)
    }

    return q
}