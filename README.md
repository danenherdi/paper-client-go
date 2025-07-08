# paper-client-go

The Go [PaperCache](https://papercache.io) client. The client supports all commands described in the wire protocol on the homepage.

## Example
```go
client, err := Connect("paper://127.0.0.1:3145")

if err != nil {
  // handle error
}

client.Set("key", "value", 0)
got, err := client.Get("key")

if err != nil {
  // handle error
}
```
