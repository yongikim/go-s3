package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/joho/godotenv"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}

	region := os.Getenv("AWS_REGION")
	bucket := os.Getenv("AWS_BUDGET")
	path := os.Getenv("AWS_CSV_PATH")
	access_key := os.Getenv("AWS_ACCESS_KEY")
	secret_key := os.Getenv("AWS_SECRET_KEY")

	svc := s3.New(session.New(), &aws.Config{
		Region: aws.String(region),
		Credentials: credentials.NewStaticCredentialsFromCreds(credentials.Value{
			AccessKeyID:     access_key,
			SecretAccessKey: secret_key,
		}),
	})

	obj, _ := svc.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(path),
	})
	defer obj.Body.Close()

	buf := new(bytes.Buffer)
	buf.ReadFrom(obj.Body)
	reader := transform.NewReader(strings.NewReader(buf.String()), japanese.ShiftJIS.NewDecoder())
	csv_reader := csv.NewReader(reader)

	var line []string
	i := 0
	for {
		i++

		line, err = csv_reader.Read()
		if err != nil {
			break
		}

		if i == 1 {
			s := ""
			for j, content := range line {
				s += strconv.Itoa(j) + ": " + content + ", "
			}
			fmt.Println(s)
			continue
		}

		fmt.Println(line)
	}
}
