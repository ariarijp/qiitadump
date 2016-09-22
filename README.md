# qiitadump

The qiitadump command can generate output in [JSON Lines](http://jsonlines.org/) format.

## Usage

```shell
$ go get github.com/ariarijp/qiitadump
$ qiitadump -h
Usage of ./qiitadump:
  -endpoint string
    	Endpoint (default "/api/v2/authenticated_user/items")
  -host string
    	Host (default "qiita.com")
  -limit int
    	Limit the number of items (default 20)
  -token string
    	Access token
  -without-private
    	Dump without private items (default true)
$ qiitadump -token YOUR_ACCESS_TOKEN > dump.json
```

## License

MIT

## Author

[Takuya Arita](https://github.com/ariarijp)
