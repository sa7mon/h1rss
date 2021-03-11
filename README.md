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

* **bounty** (one of the following)
  * true
  * false
* **state** (pipe-separated list)
  * duplicate
  * informative
  * not-applicable
  * resolved

Example: `http://localhost:8000/rss?bounty=true&state=resolved|informative`

## Screenshots

The feed as rendered by `miniflux`

<img src="https://user-images.githubusercontent.com/3712226/110651862-71815e00-8181-11eb-9459-bcdfecf5e327.png" width="800" />

<img src="https://user-images.githubusercontent.com/3712226/110651857-70e8c780-8181-11eb-8bba-b2da08db7332.png" width="800" />

## License

MIT
