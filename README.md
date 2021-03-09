# H1RSS

An RSS feed generator for HackerOne Hacktivity

## Usage

The HTTP server will listen on port 8000 for requests to /rss

```
 -interval int
        Minutes to wait between scrapes (default 120)
```

**Build from source**

```
go build -o h1rss main.go
./h1rss -interval 60
```

## License

MIT