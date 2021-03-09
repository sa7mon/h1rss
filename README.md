# H1RSS

An RSS feed generator for HackerOne Hacktivity

## Usage

The HTTP server will listen on port 8000 for requests to /rss

```
  -bind string
        Address and port to bind to (default ":8000")
  -interval int
        Minutes to wait between scrapes (default 120)

```

**Build from source**

```
go build -o h1rss main.go
./h1rss -interval 60 -bind localhost:8181
```

## License

MIT