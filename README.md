# H1RSS

An RSS feed generator for HackerOne Hacktivity

Public instance: https://h1rss.badtech.xyz/rss

## Running

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

## Feed Subscription

By default, the HTTP server will listen on port 8000 for requests to `/rss`

The items can be filtered via query parameters:

* **bounty**
    * true
    * false

Example: `http://localhost:8000/rss?bounty=true`

## License

MIT