# This is a sample plugin structure.

This code sample use AWS Go SDK to provide an idea about how to write a plugin for Kubernetes cluster auto scalar.

Plugin should have `ModifyEventAnalyzer` and `DeleteEventAnalyzer` main function. Main `autoscalar` will search for these two functions in the plugin.

Build the plugin,
```
go build -buildmode=plugin -o AWS.so main.go
```

Read more about AWS Go SDK here,
[AWS SDK for Go Developer Guide](https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/welcome.html)

[AWS SDK for Go API Reference](https://docs.aws.amazon.com/sdk-for-go/api/)
