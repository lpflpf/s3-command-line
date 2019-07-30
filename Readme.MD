## 一个S3 Golang 客户端的实现

目前实现如下功能：
  - 文件的上传
  - 文件的下载
  - 文件的删除
  - 文件大小，更新时间查看

### session config

```json
{
  "id": "",
  "secret": "",
  "endPoint": "",
  "region": ""
}
```

```
Usage:
        s3 <command> [arguments]

The commands are:

        upload    upload file to s3
        download  download file from s3
        delete    delete file in s3
        list      list files in s3

Use "s3 <command> " for more information about a command.
```

### upload

```
Usage of S3 Uploader:

        ./s3 upload -c config.json [-d directory] <bucket> <file ...>

        -c config.json  config path
        -d directory    local directory
        bucket          bucket name


```

### download

```
Usage of S3 Downloader:

        ./s3 download -c config.json [-d directory] [-x] <bucket> <file ...>

        -c config.json  config path
        -d directory    local directory
        -x              set output is stdout
        bucket          bucket name
```

### delete

```
Usage of S3 Downloader:

        ./s3 delete -c config.json <bucket> <file ...>

        -c config.json  config path
        bucket          bucket name

```

### list

```

Usage of S3 List:

        ./s3 list -c config.json [ -h]  <bucket> <pattern>

        -c config.json    config path
        -h                human reading
        bucket            bucket name
        pattern           file pattern
```