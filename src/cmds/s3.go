package main

import (
	"fmt"
	"os"
	"s3client"
)

func main() {
	if len(os.Args) == 1 {
		Usage()
	}

	var arguments []string

	if len(os.Args) > 2 {
		arguments = os.Args[2:]
	}

	switch os.Args[1] {
	case "upload":
		s3client.Uploader(arguments)
	case "ls":
	case "download":
		s3client.Downloader(arguments)
	case "delete":
		s3client.Delete(arguments)
	case "list":
		s3client.List(arguments)
	default:
		Usage()
	}
}

func Usage() {
	usage := `
Usage:
        s3 <command> [arguments]

The commands are:

        upload    upload file to s3
        download  download file from s3
        delete    delete file in s3
        list      list files in s3

Use "s3 <command> " for more information about a command.

`

	fmt.Print(usage)
	os.Exit(0)
}
