# logs

A package that provides enriched logs:

- Level support
- JSON format
- Customizable JSON encoder
- Customizable output

## Install

```sh
go get -u github.com/aukilabs/hagall-common/logs
```

## Usage

### Info

```go
logs.New().Info("hi")
logs.WithTag("to", "ted").Info("hi")
```

### Error

```go
logs.Error(errors.New("an error"))

logs.WithTag("type", "simulation").
	Error(errors.New("an error"))
```

_Note that errors created with `errors` package enriches logs by default._

### Debug

```go
logs.New().Warn("hi")
logs.WithTag("to", "ted").Warn("hi")
```

### Warning

```go
logs.New().Warn("hi")
logs.WithTag("to", "ted").Warn("hi")
```

### Enrich With Tags

```go
logs.WithTag("method", "GET").
    WithTag("path", "/cookies").
    Info("http request")
```

### Set A Custom Logger

```go
logs.SetLogger(func(e Entry) {
    fmt.Println(e) // or whatever you want to use.
})
```

### Set A Custom JSON Encoder

```go
logs.Encoder = json.Marshal // or whathever JSON encoder you want to use.
```
