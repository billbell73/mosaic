// fetches all images from Amazon s3 bucket and saves name
// and red-green-blue averages in local database
package main

import (
	"database/sql"
	"image"
	"log"
	"os"

	"github.com/billbell73/mosaic/lib/imageAnalyser"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"

	_ "image/jpeg"

	_ "github.com/lib/pq"
)

var awsBucket string

func init() {
	awsBucket = os.Getenv("AWS_BUCKET")
	if awsBucket == "" {
		log.Fatal("$AWS_BUCKET must be set")
	}
}

func fetchImageNames(svc *s3.S3) []string {
	params := &s3.ListObjectsInput{
		Bucket: aws.String(awsBucket),
	}

	resp, err := svc.ListObjects(params)
	if err != nil {
		log.Fatal(err)
	}

	var imageNames []string

	for _, key := range resp.Contents {
		imageNames = append(imageNames, *key.Key)
	}
	return imageNames
}

func saveToDb(stmt *sql.Stmt, imageName string, avg [3]int) {
	_, err := stmt.Exec(imageName, avg[0], avg[1], avg[2])
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	stmt, err := db.Prepare("INSERT INTO averages(name, red, green, blue) VALUES($1, $2, $3, $4)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	svc := s3.New(session.New())
	imageNames := fetchImageNames(svc)

	for _, name := range imageNames {
		obj, err := svc.GetObject(&s3.GetObjectInput{
			Bucket: aws.String(awsBucket),
			Key:    aws.String(name),
		})
		if err != nil {
			log.Fatal(err)
		}

		img, _, err := image.Decode(obj.Body)
		if err != nil {
			log.Fatal(err)
		}
		avg := imageAnalyser.AverageRGB(img, imageAnalyser.TotalRGB)

		saveToDb(stmt, name, avg)
	}
}
