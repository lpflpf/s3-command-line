package s3client

import (
	"flag"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func Uploader(args []string) {
	var configPath, dir string
	uploaderFlag := flag.NewFlagSet("s3 uploader", flag.ExitOnError)

	uploaderFlag.StringVar(&configPath, "c", "", "config path")
	uploaderFlag.StringVar(&dir, "d", "", "s3 directory.")
	uploaderFlag.Usage = func() {
		usage := `Usage of S3 Uploader:

        ./s3 upload -c config.json [-d directory] <bucket> <file ...>

        -c config.json  config path
        -d directory    local directory
        bucket          bucket name

`
		fmt.Print(usage)
	}
	if err := uploaderFlag.Parse(args); err != nil {
		uploaderFlag.Usage()
		os.Exit(0)
	}

	if len(uploaderFlag.Args()) < 2 {
		uploaderFlag.Usage()
		os.Exit(0)
	}

	bucket, files := uploaderFlag.Args()[0], uploaderFlag.Args()[1:]

	if config, err := Load(configPath); err != nil {
		log.Fatal(err)
	} else {
		sess := NewSession(config)
		if dir != "" {
			dir = strings.Trim(dir, string([]byte{filepath.Separator}))
		}
		for i := 0; i < len(files); i++ {
			filename := filepath.Base(files[i])
			if dir != "" {
				filename = dir + "/" + filename
			}
			uploadFile(sess, files[i], bucket, filename)
		}
	}
}

func uploadFile(ses *session.Session, path string, bucket string, key string) {
	if isExist, err := fileExist(path); !isExist || err != nil {
		log.Fatalf("file %s cannot be find.\n", path)
	}

	if fd, err := os.Open(path); err != nil {
		log.Fatalf("file %s open failed.\n", path)
	} else {
		uploader := s3manager.NewUploader(ses)
		_, err := uploader.Upload(&s3manager.UploadInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(key),
			Body:   fd,
		})

		if err != nil {
			log.Printf("upload %s failed, error:%v", path, err)
		} else {
			log.Printf("upload %s success!\n", path)
		}
	}
}
