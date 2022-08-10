package openapi

const (
	FormatProperties Format = "properties"
	FormatXML        Format = "xml"
	FormatYML        Format = "yml"
	FormatYAML       Format = "yaml"
	FormatJSON       Format = "json"
)

type Format string

func (f Format) String() string {
	return string(f)
}

type OpenAPI interface {
	GetApps() ([]Application, error)
	GetEnvClusters(appID string) ([]EnvWithClusters, error)
	GetNamespaces(env, appID, clusterName string) ([]Namespace, error)
	GetNamespace(env, appID, clusterName, namespaceName string) (*Namespace, error)
	CreateNamespace(r CreateNamespaceRequest) (*CreateNamespaceResponse, error)
	GetNamespaceLock(env, appID, clusterName, namespaceName string) (*NamespaceLock, error)
	AddItem(env, appID, clusterName, namespaceName string, r AddItemRequest) (*Item, error)
	UpdateItem(env, appID, clusterName, namespaceName string, r UpdateItemRequest) error
	CreateOrUpdateItem(env, appID, clusterName, namespaceName string, r UpdateItemRequest) error
	DeleteItem(env, appID, clusterName, namespaceName, key, operator string) error
	PublishRelease(env, appID, clusterName, namespaceName string, r PublishReleaseRequest) (*Release, error)
	GetRelease(env, appID, clusterName, namespaceName string) (*Release, error)
}

type Application struct {
	Name                       string `json:"name"`
	AppID                      string `json:"appId"`
	OrgID                      string `json:"orgId"`
	OrgName                    string `json:"orgName"`
	OwnerName                  string `json:"ownerName"`
	OwnerEmail                 string `json:"ownerEmail"`
	DataChangeCreatedBy        string `json:"dataChangeCreatedBy"`
	DataChangeLastModifiedBy   string `json:"dataChangeLastModifiedBy"`
	DataChangeCreatedTime      string `json:"dataChangeCreatedTime"`
	DataChangeLastModifiedTime string `json:"dataChangeLastModifiedTime"`
}

type EnvWithClusters struct {
	Env      string   `json:"env"`
	Clusters []string `json:"clusters"`
}

type Item struct {
	Key                        string `json:"key"`
	Value                      string `json:"value"`
	Comment                    string `json:"comment"`
	DataChangeCreatedBy        string `json:"dataChangeCreatedBy"`
	DataChangeLastModifiedBy   string `json:"dataChangeLastModifiedBy"`
	DataChangeCreatedTime      string `json:"dataChangeCreatedTime"`
	DataChangeLastModifiedTime string `json:"dataChangeLastModifiedTime"`
}

type Namespace struct {
	AppID                      string `json:"appId"`
	ClusterName                string `json:"clusterName"`
	NamespaceName              string `json:"namespaceName"`
	Comment                    string `json:"comment"`
	Format                     string `json:"format"`
	IsPublic                   bool   `json:"isPublic"`
	Items                      []Item `json:"items"`
	DataChangeCreatedBy        string `json:"dataChangeCreatedBy"`
	DataChangeLastModifiedBy   string `json:"dataChangeLastModifiedBy"`
	DataChangeCreatedTime      string `json:"dataChangeCreatedTime"`
	DataChangeLastModifiedTime string `json:"dataChangeLastModifiedTime"`
}

//{"name":"TEST1.sms.yml","appId":"conn.mysql","format":"yml","isPublic":true,"appendNamespacePrefix":true,"comment":"migrate namespace from consul","dataChangeCreatedBy":"apollo","dataChangeLastModifiedBy":"apollo","dataChangeCreatedTime":"2019-07-19T19:52:03.923+0800","dataChangeLastModifiedTime":"2019-07-19T19:52:03.923+0800"}
type CreateNamespaceResponse struct {
	Name                       string `json:"name"`
	AppID                      string `json:"appId"`
	Format                     string `json:"format"`
	IsPublic                   bool   `json:"isPublic"`
	AppendNamespacePrefix      bool   `json:"appendNamespacePrefix"`
	Comment                    string `json:"comment"`
	DataChangeCreatedBy        string `json:"dataChangeCreatedBy"`
	DataChangeLastModifiedBy   string `json:"dataChangeLastModifiedBy"`
	DataChangeCreatedTime      string `json:"dataChangeCreatedTime"`
	DataChangeLastModifiedTime string `json:"dataChangeLastModifiedTime"`
}

type CreateNamespaceRequest struct {
	Name                  string `json:"name"`                // Namespace的名字
	AppID                 string `json:"appId"`               // Namespace所属的AppId
	Format                Format `json:"format"`              // Namespace的格式，只能是以下类型： properties、xml、json、yml、yaml
	IsPublic              bool   `json:"isPublic"`            // 是否是公共文件
	DataChangeCreatedBy   string `json:"dataChangeCreatedBy"` // namespace的创建人，格式为域账号，也就是sso系统的User ID
	Comment               string `json:"comment"`             // Namespace说明(非必选)
	AppendNamespacePrefix bool   `json:"appendNamespacePrefix"`
}

type NamespaceLock struct {
	NamespaceName string `json:"namespaceName"`
	IsLocked      bool   `json:"isLocked"`
	LockedBy      string `json:"lockedBy"` //锁owner
}

type AddItemRequest struct {
	Key                 string `json:"key"`                 // 配置的key，长度不能超过128个字符。非properties格式，key固定为content
	Value               string `json:"value"`               // 配置的value，长度不能超过20000个字符，非properties格式，value为文件全部内容
	Comment             string `json:"comment"`             // 配置的备注,长度不能超过1024个字符(非必选)
	DataChangeCreatedBy string `json:"dataChangeCreatedBy"` // item的创建人，格式为域账号，也就是sso系统的User ID
}

type UpdateItemRequest struct {
	Key                      string `json:"key"`                      // 配置的key，长度不能超过128个字符。非properties格式，key固定为content
	Value                    string `json:"value"`                    // 配置的value，长度不能超过20000个字符，非properties格式，value为文件全部内容
	Comment                  string `json:"comment"`                  // 配置的备注,长度不能超过1024个字符(非必选)
	DataChangeLastModifiedBy string `json:"dataChangeLastModifiedBy"` // item的修改人，格式为域账号，也就是sso系统的User ID
	DataChangeCreatedBy      string `json:"dataChangeCreatedBy"`      // 当createIfNotExists为true时必选。item的创建人，格式为域账号，也就是sso系统的User ID
}

type PublishReleaseRequest struct {
	ReleaseTitle   string `json:"releaseTitle"`   // 此次发布的标题，长度不能超过64个字符
	ReleaseComment string `json:"releaseComment"` // 发布的备注，长度不能超过256个字符(非必选)
	ReleasedBy     string `json:"releasedBy"`     // 发布人，域账号，注意：如果ApolloConfigDB.ServerConfig中的namespace.lock.switch设置为true的话（默认是false），那么该环境不允许发布人和编辑人为同一人。所以如果编辑人是zhanglea，发布人就不能再是zhanglea。
}

type Release struct {
	AppID                      string            `json:"appId"`
	ClusterName                string            `json:"clusterName"`
	NamespaceName              string            `json:"namespaceName"`
	Name                       string            `json:"name"`
	Configurations             map[string]string `json:"configurations"`
	Comment                    string            `json:"comment"`
	DataChangeCreatedBy        string            `json:"dataChangeCreatedBy"`
	DataChangeLastModifiedBy   string            `json:"dataChangeLastModifiedBy"`
	DataChangeCreatedTime      string            `json:"dataChangeCreatedTime"`
	DataChangeLastModifiedTime string            `json:"dataChangeLastModifiedTime"`
}
