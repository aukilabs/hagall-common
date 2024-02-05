# metrics

A package that provides reusable metrics accross Aukilabs Go projects.

## Install

```sh
go get -u github.com/aukilabs/hagall-common/metrics
```

## HTTP

### Server

```go
var h http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// A standard http handler.
})

h = metrics.HTTPHandler(h)
http.Handle("/", h)
```

### Client

```go
client := http.Client{
	Transport: HTTPTransport(http.DefaultTransport),
}

res, err := client.Get("https://ted.wushu")
```
