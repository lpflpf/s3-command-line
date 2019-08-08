package s3client

import (
	"flag"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func List(arguments []string) {
	var configPath, prefix string
	var humanReading bool
	downloadFlags := flag.NewFlagSet("s3 List", flag.ExitOnError)
	downloadFlags.StringVar(&configPath, "c", "", "config path")
	downloadFlags.BoolVar(&humanReading, "h", false, "human reading")
	downloadFlags.Usage = func() {
		usage := `Usage of S3 List:

        ./s3 ls -c config.json [ -h]  <bucket> <pattern>

        -c config.json    config path
        -h                human reading
        bucket            bucket name
        pattern           file pattern

        example:
            ./s3 ls -c config.json example-bucket /pattern*
            ./s3 ls -c config.json example-bucket /pattern-all-name
            ./s3 ls -c config.json example-bucket /pattern*.tar.gz

`
		fmt.Print(usage)
	}

	if err := downloadFlags.Parse(arguments); err != nil {
		downloadFlags.Usage()
		os.Exit(0)
	}

	if len(downloadFlags.Args()) < 2 {
		downloadFlags.Usage()
		os.Exit(0)
	}

	bucket, prefix := downloadFlags.Args()[0], downloadFlags.Args()[1]

	if config, err := Load(configPath); err != nil {
		log.Fatal(err)
	} else {
		sess := NewSession(config)
		prefix = strings.Trim(prefix, string([]byte{filepath.Separator}))

		listFile(sess, prefix, bucket, humanReading)
	}
}

func listFile(ses *session.Session, path string, bucket string, humanReading bool) {
	s3client := s3.New(ses)
	listObjectOutput, err := s3client.ListObjects(&s3.ListObjectsInput{
		Bucket: aws.String(bucket),
		Prefix: aws.String(path),
	})

	if err != nil {
		log.Fatal(err)
	} else {
		location := time.FixedZone("CST", 8*3600)
		if humanReading {
			fmt.Printf("%-45s%4s %-25s\n", "Filename", "Size", "Time")
		} else {
			fmt.Printf("%-45s%-15s%-25s\n", "Filename", "Size", "Time")
		}
		for _, output := range listObjectOutput.Contents {
			if humanReading {
				fmt.Printf("%-45s%4s %-25s\n",
					*output.Key,
					formatSize(*output.Size),
					output.LastModified.In(location).Format("2006-01-02 15:04:05"))
			} else {
				fmt.Printf("%-45s%-15d%-25s\n",
					*output.Key,
					*output.Size,
					output.LastModified.In(location).Format("2006-01-02 15:04:05"))
			}
		}
	}
}

func formatSize(size int64) (output string) {

	idx := -1

	fSize := float64(size)

	for fSize > 1024 {
		fSize = fSize / 1024
		idx++
	}

	units := []byte{'K', 'M', 'G'}
	switch {
	case idx == -1:
		return fmt.Sprintf("%3d", size)
	case idx >= 0 && idx <= 2:
		if fSize >= 10 {
			return fmt.Sprintf("%3.0f%c", fSize, units[idx])
		}
		return fmt.Sprintf("%1.1f%c", fSize, units[idx])
	case idx > 2:
		return fmt.Sprintf("INF")
	}

	return ""
}
