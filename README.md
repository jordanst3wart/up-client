Go client for up

https://developer.up.com.au

This is just quick client code for up bank. Largely, created with claude code generation.

The openapi spec is available at https://github.com/up-banking/api/blob/master/v1/openapi.json

```go
client := upclient.NewClient("your-token-here")
pingResp, err := client.Ping()
if err != nil {
    log.Fatal(err)
}
fmt.Printf("API Status: %s\n", pingResp.Meta.StatusEmoji)
```

```go
client := up.NewClient("your-token-here", nil)

// List accounts
accounts, resp, err := client.Accounts.List(
```