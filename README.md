# kong-go
Simple Go client for Kong Admin API

### Usage:

```go
kapi := NewKongClient(http.DefaultClient, "http://192.168.99.100:30081")

apis, err := kapi.GetAPIs()
...
api, err := kapi.AddAPI(KongAPI{
		Name:         "kong-test-api",
		Request_path: "/kong-test-api",
		Upstream_url: "http://www.targetprocess.com",
	})
...
api, err := kapi.GetAPI("kong-test-api")
...
kapi.DeleteAPI("kong-test-api")
```
