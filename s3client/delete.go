package s3client

import (
	"flag"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"log"
	"os"
	"strings"
)

func Delete(args []string) {
	var configPath string
	deleteFlags := flag.NewFlagSet("s3 rm ", flag.ExitOnError)
	deleteFlags.StringVar(&configPath, "c", "", "config path")
	deleteFlags.Usage = func() {
		usage := `Usage of S3 Delete:

        ./s3 rm -c config.json <bucket> <file ...>

        -c config.json  config path
        bucket          bucket name

`
		fmt.Print(usage)
	}

	if err := deleteFlags.Parse(args); err != nil {
		deleteFlags.Usage()
		os.Exit(0)
	}

	if len(deleteFlags.Args()) < 2 {
		deleteFlags.Usage()
		os.Exit(0)
	}

	bucket, files := deleteFlags.Args()[0], deleteFlags.Args()[1:]

	if config, err := Load(configPath); err != nil {
		log.Fatal(err)
	} else {
		sess := NewSession(config)
		del(sess, bucket, files)
	}
}

func del(ses *session.Session, bucket string, keys []string) {
	var objects []s3manager.BatchDeleteObject

	for _, key := range keys {
		objects = append(objects, s3manager.BatchDeleteObject{
			Object: &s3.DeleteObjectInput{
				Key:    aws.String(key),
				Bucket: aws.String(bucket),
			},
		})
	}
	deleteObjectIterator := s3manager.DeleteObjectsIterator{
		Objects: objects,
	}
	deleter := s3manager.NewBatchDelete(ses)
	if err := deleter.Delete(aws.BackgroundContext(), &deleteObjectIterator); err == nil {
		log.Printf("delete \"%s\" success.", strings.Join(keys, ","))
	} else {
		log.Printf("delete \"%s\" failed. error: %s", strings.Join(keys, ","), err)
	}
}
