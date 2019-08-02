package s3client

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func Downloader(args []string) {
	var configPath, dir string
	var standardOutput bool

	downloadFlags := flag.NewFlagSet("s3 downloader", flag.ExitOnError)
	downloadFlags.StringVar(&configPath, "c", "", "config path")
	downloadFlags.StringVar(&dir, "d", "./", "local directory")
	downloadFlags.BoolVar(&standardOutput, "x", false, "stdout")
	downloadFlags.Usage = func() {
		usage := `Usage of S3 Downloader:

        ./s3 download -c config.json [-d directory] [-x] <bucket> <file ...>

        -c config.json  config path
        -d directory    local directory
        -x              set output is stdout
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
		if standardOutput {
			for _, file := range files {
				downloadByStdout(sess, bucket, file)
			}
		} else {
			dir = strings.Trim(dir, string([]byte{filepath.Separator}))
			for _, file := range files {
				download(sess, dir, bucket, file)
			}
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
		err := s3download(ses, fd, bucket, key)
		if err != nil {
			log.Printf("download \"%s\" failed, error:%v", filename, err)
		} else {
			log.Printf("download \"%s\" success.", filename)
		}
	}
}

func downloadByStdout(ses *session.Session, bucket string, key string) {
	//err := s3download(ses, os.Stdout, bucket, key)
	s3client := s3.New(ses)

	output, err := s3client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(bucket),

		Key: aws.String(key),
	})

	if err != nil {
		log.Printf("output %s failed, error:%v", key, err)
		return
	} else {
		log.Printf("output %s success.", key)
	}
	bufReader := bufio.NewReader(output.Body)

	for {
		data, err := bufReader.ReadBytes('\n')
		fmt.Printf("%s", string(data))
		if err == io.EOF {
			break
		}
	}
}

func s3download(ses *session.Session, writer io.WriterAt, bucket string, key string) error {
	downloader := s3manager.NewDownloader(ses)
	_, err := downloader.Download(writer, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	return err
}
