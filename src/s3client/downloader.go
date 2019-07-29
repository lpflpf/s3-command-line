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
	"path/filepath"
	"strings"
)

func Downloader(args []string) {
	var configPath, dir string
	downloadFlags := flag.NewFlagSet("s3 downloader", flag.ExitOnError)
	downloadFlags.StringVar(&configPath, "c", "", "config path")
	downloadFlags.StringVar(&dir, "o", "./", "local directory")
	downloadFlags.Usage = func() {
		usage := `Usage of S3 Downloader:

        ./s3 download -c config.json [-o directory] <bucket> <file ...>

        -c config.json  config path
        -o directory    local directory
        bucket          bucket name

`
		fmt.Print(usage)
	}

	if err := downloadFlags.Parse(args); err != nil {
		downloadFlags.Usage()
		os.Exit(0)
	}

	if len(downloadFlags.Args()) < 2 {
		downloadFlags.Usage()
		os.Exit(0)
	}

	bucket, files := downloadFlags.Args()[0], downloadFlags.Args()[1:]

	if config, err := Load(configPath); err != nil {
		log.Fatal(err)
	} else {
		sess := NewSession(config)
		dir = strings.Trim(dir, string([]byte{filepath.Separator}))

		for i := 0; i < len(files); i++ {
			download(sess, dir, bucket, files[i])
		}
	}
}

func download(ses *session.Session, path string, bucket string, key string) {
	if isExist, err := fileExist(path); !isExist {
		if err := os.MkdirAll(path, 0775); err != nil {
			log.Panic("create dir failed\n", err)
		}
	} else if err != nil {
		log.Fatalf("cannot access dir:%s\n", path)
	}
	filename := path + "/" + filepath.Base(key)

	if fd, err := os.Create(filename); err != nil {
		log.Fatalf("file \"%s\" open failed. error:%v\n", filename, err)
	} else {
		downloader := s3manager.NewDownloader(ses)
		_, err := downloader.Download(fd, &s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(key),
		})

		if err != nil {
			log.Printf("download \"%s\" failed, error:%v", filename, err)
		} else {
			log.Printf("download \"%s\" success.", filename)
		}
	}
}
