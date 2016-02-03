udpchan [![Circle CI](https://circleci.com/gh/Cistern/udpchan.svg?style=svg&circle-token=18c5b364cdae51455ba9a70c4809686f1cf41a63)](https://circleci.com/gh/Cistern/udpchan) [![GoDoc](https://godoc.org/github.com/PreetamJinka/udpchan?status.png)](https://godoc.org/github.com/PreetamJinka/udpchan)
===
A tiny channel wrapper around UDP connections

Usage
---
It's pretty simple: call `Connect` or `Listen` and get a `[]byte` channel back!

```go
inbound, err := Listen(":9999", nil)
if err != nil {
	// handle err
}

outbound, err := Connect(":9999")
if err != nil {
	// handle err
}

message := []byte("foo")

// Send a message over UDP
outbound <- message

// Receive a message over UDP
read := <-inbound // = []byte("foo")
```

License
---
MIT
