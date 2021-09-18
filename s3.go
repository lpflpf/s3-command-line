package main

import (
	"fmt"
	"os"

	"github.com/lpflpf/s3-command-line/s3client"
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
	case "put":
		s3client.Uploader(arguments)
	case "get":
		s3client.Downloader(arguments)
	case "rm":
		s3client.Delete(arguments)
	case "ls":
		s3client.List(arguments)
	case "geturl":
		s3client.GetUrl(arguments)
	default:
		Usage()
	}
}

func Usage() {
	usage := `
Usage:
        s3 <command> [arguments]

The commands are:

        put       upload file to s3
        get       download file from s3
        geturl    get download file url
        rm        delete file in s3
        ls        list files in s3

Use "s3 <command> " for more information about a command.

`

	fmt.Print(usage)
	os.Exit(0)
}
