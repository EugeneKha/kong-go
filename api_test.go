package api

import (
	"net/http"
	"testing"
)

func Test_GetApiList(t *testing.T) {

	kapi := NewKongClient(http.DefaultClient, "http://192.168.99.100:30081/ ")

	apis, err := kapi.GetAPIs()

	if err != nil {
		t.Error("ERROR getting APIs:", err)
	}

	count := len(apis.Data)

	api, err := kapi.AddAPI(KongAPI{
		Name:         "kong-test-api",
		Request_path: "/kong-test-api",
		Upstream_url: "http://www.targetprocess.com",
	})

	if err != nil {
		t.Error("ERROR adding API:", err)
	}

	apis, err = kapi.GetAPIs()

	if err != nil {
		t.Error("ERROR getting APIs:", err)
	}

	if len(apis.Data) != count+1 {
		t.Error("Can't find added api in apis list")
	}

	addedApi := KongAPI{}
	for _, a := range apis.Data {
		if a.Name == api.Name {
			addedApi = a
			break
		}
	}

	if addedApi.Id == "" ||
		addedApi.Name != api.Name ||
		addedApi.Request_path != api.Request_path ||
		addedApi.Upstream_url != api.Upstream_url {
		t.Error("Found API does not match", addedApi)
	}

	addedApi, err = kapi.GetAPI(api.Name)
	if err != nil {
		t.Error("ERROR getting APIs:", err)
	}

	if addedApi.Id == "" ||
		addedApi.Name != api.Name ||
		addedApi.Request_path != api.Request_path ||
		addedApi.Upstream_url != api.Upstream_url {
		t.Fail()
	}

	err = kapi.DeleteAPI(api.Name)

	if err != nil {
		t.Error("ERROR deleting APIs:", err)
	}

	apis, err = kapi.GetAPIs()

	if err != nil {
		t.Error("ERROR getting APIs:", err)
	}

	if len(apis.Data) != count {
		t.Error("Can't find added api in apis list")
	}
}
