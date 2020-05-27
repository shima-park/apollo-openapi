package openapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

type Doer interface {
	Do(*http.Request) (*http.Response, error)
}

type LoggerFunc func(fotmat string, args ...interface{})

type ClientOptions struct {
	Doer       Doer
	Debug      bool
	LoggerFunc LoggerFunc
}

type ClientOption func(*ClientOptions)

func WithDoer(d Doer) ClientOption {
	return func(o *ClientOptions) {
		o.Doer = d
	}
}

func WithDebug(d bool) ClientOption {
	return func(o *ClientOptions) {
		o.Debug = d
	}
}

func WithLoggerFunc(lf LoggerFunc) ClientOption {
	return func(o *ClientOptions) {
		o.LoggerFunc = lf
	}
}

type client struct {
	portalAddress string
	token         string
	options       ClientOptions
}

func NewClient(portalAddress, token string, opts ...ClientOption) OpenAPI {
	var options ClientOptions
	for _, opt := range opts {
		opt(&options)
	}

	if options.Doer == nil {
		options.Doer = &http.Client{}
	}

	if options.LoggerFunc == nil {
		options.LoggerFunc = func(format string, args ...interface{}) {
			format += "\n"
			fmt.Printf(format, args...)
		}
	}

	return &client{
		portalAddress: normalizeURL(portalAddress),
		token:         token,
		options:       options,
	}
}

func (c *client) newRequest(method, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", c.token)
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	return req, nil
}

func (c *client) debug(format string, args ...interface{}) {
	if c.options.Debug {
		c.options.LoggerFunc(format, args...)
	}
}

func (c *client) do(method, url string, request, response interface{}) error {
	var (
		reqBodyReader io.Reader
		reqBody       []byte
		err           error
		req           *http.Request
		respBody      []byte
		status        int
	)
	defer func(reqBody, respBody *[]byte) {
		c.debug("Method: %s, URL: %s, \n Request body: %s,\n Response body: %s",
			method, url, *reqBody, *respBody)
	}(&reqBody, &respBody)

	if request != nil {
		reqBody, err = json.Marshal(request)
		if err != nil {
			return err
		}

		reqBodyReader = bytes.NewReader(reqBody)
	}

	req, err = c.newRequest(method, url, reqBodyReader)
	if err != nil {
		return err
	}

	status, respBody, err = parseResponseBody(c.options.Doer, req)
	if err != nil {
		return err
	}

	if status == http.StatusOK {
		if response != nil {
			return json.Unmarshal(respBody, response)
		}
		return nil
	}

	return errors.New(getErrorMessage(status))
}

func getErrorMessage(status int) string {
	switch status {
	case 400:
		return "400 - Bad Request 客户端传入参数的错误，如操作人不存在，namespace不存在等等，客户端需要根据提示信息检查对应的参数是否正确。"
	case 401:
		return "401 - Unauthorized 接口传入的token非法或者已过期，客户端需要检查token是否传入正确。"
	case 403:
		return "403 - Forbidden 接口要访问的资源未得到授权，比如只授权了对A应用下Namespace的管理权限，但是却尝试管理B应用下的配置。"
	case 404:
		return "404 - Not Found 接口要访问的资源不存在，一般是URL或URL的参数错误。"
	case 405:
		return "405 - Method Not Allowed 接口访问的Method不正确，比如应该使用POST的接口使用了GET访问等，客户端需要检查接口访问方式是否正确。"
	case 500:
		return "500 - Internal Server Error 其它类型的错误默认都会返回500，对这类错误如果应用无法根据提示信息找到原因的话，可以找Apollo研发团队一起排查问题。"
	default:
		return fmt.Sprintf("未定义错误码: %d", status)
	}
}

func (c *client) GetEnvClusters(appID string) (res []EnvWithClusters, err error) {
	url := fmt.Sprintf("%s/openapi/v1/apps/%s/envclusters", c.portalAddress, appID)
	err = c.do("GET", url, nil, &res)
	return
}

func (c *client) GetNamespaces(env, appID, clusterName string) (res []Namespace, err error) {
	url := fmt.Sprintf("%s/openapi/v1/envs/%s/apps/%s/clusters/%s/namespaces",
		c.portalAddress, env, appID, clusterName)
	err = c.do("GET", url, nil, &res)
	return
}

func (c *client) GetNamespace(env, appID, clusterName, namespaceName string) (res *Namespace, err error) {
	namespaceName = normalizeNamespace(namespaceName)
	url := fmt.Sprintf("%s/openapi/v1/envs/%s/apps/%s/clusters/%s/namespaces/%s",
		c.portalAddress, env, appID, clusterName, namespaceName)
	res = &Namespace{}
	err = c.do("GET", url, nil, &res)
	return
}

func (c *client) CreateNamespace(r CreateNamespaceRequest) (res *CreateNamespaceResponse, err error) {
	url := fmt.Sprintf("%s/openapi/v1/apps/%s/appnamespaces?appendNamespacePrefix=%v",
		c.portalAddress, r.AppID, r.AppendNamespacePrefix)
	res = &CreateNamespaceResponse{}
	err = c.do("POST", url, r, &res)
	return
}

func (c *client) GetNamespaceLock(env, appID, clusterName, namespaceName string) (res *NamespaceLock, err error) {
	namespaceName = normalizeNamespace(namespaceName)
	url := fmt.Sprintf("%s/openapi/v1/envs/%s/apps/%s/clusters/%s/namespaces/%s/lock",
		c.portalAddress, env, appID, clusterName, namespaceName)
	res = &NamespaceLock{}
	err = c.do("GET", url, nil, &res)
	return
}

func (c *client) AddItem(env, appID, clusterName, namespaceName string, r AddItemRequest) (res *Item, err error) {
	namespaceName = normalizeNamespace(namespaceName)
	url := fmt.Sprintf("%s/openapi/v1/envs/%s/apps/%s/clusters/%s/namespaces/%s/items",
		c.portalAddress, env, appID, clusterName, namespaceName)
	res = &Item{}
	err = c.do("POST", url, r, &res)
	return
}

func (c *client) UpdateItem(env, appID, clusterName, namespaceName string, r UpdateItemRequest) (err error) {
	namespaceName = normalizeNamespace(namespaceName)
	url := fmt.Sprintf("%s/openapi/v1/envs/%s/apps/%s/clusters/%s/namespaces/%s/items/%s",
		c.portalAddress, env, appID, clusterName, namespaceName, r.Key)
	err = c.do("PUT", url, r, nil)
	return
}

func (c *client) CreateOrUpdateItem(env, appID, clusterName, namespaceName string, r UpdateItemRequest) (err error) {
	namespaceName = normalizeNamespace(namespaceName)
	url := fmt.Sprintf("%s/openapi/v1/envs/%s/apps/%s/clusters/%s/namespaces/%s/items/%s?createIfNotExists=true",
		c.portalAddress, env, appID, clusterName, namespaceName, r.Key)
	err = c.do("PUT", url, r, nil)
	return
}

func (c *client) DeleteItem(env, appID, clusterName, namespaceName, key, operator string) error {
	namespaceName = normalizeNamespace(namespaceName)
	url := fmt.Sprintf("%s/openapi/v1/envs/%s/apps/%s/clusters/%s/namespaces/%s/items/%s?operator=%s",
		c.portalAddress, env, appID, clusterName, namespaceName, key, operator)
	return c.do("DELETE", url, nil, nil)
}

func (c *client) PublishRelease(env, appID, clusterName, namespaceName string, r PublishReleaseRequest) (res *Release, err error) {
	namespaceName = normalizeNamespace(namespaceName)
	url := fmt.Sprintf("%s/openapi/v1/envs/%s/apps/%s/clusters/%s/namespaces/%s/releases",
		c.portalAddress, env, appID, clusterName, namespaceName)
	res = &Release{}
	err = c.do("POST", url, r, &res)
	return
}

func (c *client) GetRelease(env, appID, clusterName, namespaceName string) (res *Release, err error) {
	namespaceName = normalizeNamespace(namespaceName)
	url := fmt.Sprintf("%s/openapi/v1/envs/%s/apps/%s/clusters/%s/namespaces/%s/releases/latest",
		c.portalAddress, env, appID, clusterName, namespaceName)
	res = &Release{}
	err = c.do("GET", url, nil, &res)
	return
}

func parseResponseBody(doer Doer, req *http.Request) (int, []byte, error) {
	resp, err := doer.Do(req)
	if err != nil {
		return 0, nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, nil, err
	}

	return resp.StatusCode, body, nil
}

func normalizeURL(url string) string {
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		url = "http://" + url
	}

	return strings.TrimSuffix(url, "/")
}

func normalizeNamespace(ns string) string {
	return strings.TrimSuffix(ns, "."+FormatProperties)
}
