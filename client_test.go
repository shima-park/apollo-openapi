package openapi

import (
	"fmt"
	"testing"
)

var (
	portalAddress = "localhost:8070"
	token         = "5573adf2107c1e34d9b726e1570f64a62c9f1ca5"
	c             = NewClient(portalAddress, token, WithDebug(true))
)

func TestClientGet(t *testing.T) {
	var (
		env           = "DEV"
		appID         = "SampleApp"
		clusterName   = "default"
		namespaceName = "express.yml"
	)
	t.Log(c.GetEnvClusters(appID))
	t.Log(c.GetNamespaces(env, appID, clusterName))
	t.Log(c.GetNamespace(env, appID, clusterName, namespaceName))
	t.Log(c.GetNamespaceLock(env, appID, clusterName, namespaceName))
	t.Log(c.GetRelease(env, appID, clusterName, namespaceName))
}

func TestCreateNamespace(t *testing.T) {
	var (
		env         = "DEV"
		appID       = "SampleApp"
		clusterName = "default"
		operator    = "apollo"
		formats     = []Format{
			FormatProperties,
			FormatXML,
			FormatYML,
			FormatYAML,
			FormatJSON,
		}
	)

	t.Log("Before")
	t.Log(c.GetNamespaces(env, appID, clusterName))

	for _, visibleLevel := range []bool{true, false} {
		for _, format := range formats {
			ns, err := c.CreateNamespace(CreateNamespaceRequest{
				Name:                string(format),
				AppID:               appID,
				Format:              format,
				IsPublic:            visibleLevel,
				DataChangeCreatedBy: operator,
				Comment:             fmt.Sprintf("balabala format:%s visible:%v", format, visibleLevel),
			})
			if err != nil {
				t.Fatal(err)
			}
			t.Logf("Namespace: %+v", ns)
		}
	}

	t.Log("After")
	t.Log(c.GetNamespaces(env, appID, clusterName))
}

func TestAddItem(t *testing.T) {
	var (
		env         = "DEV"
		appID       = "SampleApp"
		clusterName = "default"
		operator    = "apollo"
		namespace   = "application"
	)

	i, err := c.AddItem(env, appID, clusterName, namespace, AddItemRequest{
		Key:                 "name",
		Value:               "foo",
		DataChangeCreatedBy: operator,
		Comment:             "add item",
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Item: %+v", i)
}

func TestUpdateItem(t *testing.T) {
	var (
		env         = "DEV"
		appID       = "SampleApp"
		clusterName = "default"
		namespace   = "application"
		operator    = "apollo"
	)

	err := c.UpdateItem(env, appID, clusterName, namespace, UpdateItemRequest{
		Key:                      "name",
		Value:                    "bar",
		DataChangeLastModifiedBy: operator,
		Comment:                  "update item",
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestDeleteItem(t *testing.T) {
	var (
		env         = "DEV"
		appID       = "SampleApp"
		clusterName = "default"
		namespace   = "application"
		key         = "name"
		operator    = "apollo"
	)

	err := c.DeleteItem(env, appID, clusterName, namespace, key, operator)
	if err != nil {
		t.Fatal(err)
	}
}

func TestPublishRelease(t *testing.T) {
	var (
		env            = "DEV"
		appID          = "SampleApp"
		clusterName    = "default"
		namespace      = "application"
		releaseTitle   = "test"
		releaseComment = " comment "
		releasedBy     = "apollo"
	)

	r, err := c.PublishRelease(env, appID, clusterName, namespace, PublishReleaseRequest{
		ReleaseTitle:   releaseTitle,
		ReleaseComment: releaseComment,
		ReleasedBy:     releasedBy,
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(r)
}
