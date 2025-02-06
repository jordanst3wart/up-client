Go client for up

https://developer.up.com.au

This is just quick client code for up bank. Largely, created with claude code generation.

The openapi spec is available at https://github.com/up-banking/api/blob/master/v1/openapi.json

# Ping example
```go
func ping() {
    client := up.NewClient("your-token-here", nil)
    pingResp, _, err := client.Utility.Ping(context.TODO())
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("API Status: %s\n", pingResp.Meta.StatusEmoji)
}
```

