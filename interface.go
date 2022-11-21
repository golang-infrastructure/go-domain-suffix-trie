package domain_suffix_trie

// DomainSuffixTrieInterface 域名后缀树的接口定义
type DomainSuffixTrieInterface[T any] interface {

	// FindMatchDomainSuffixPayload
	//
	//	@Description: 根据域名查询所匹配的后缀所对应的payload，会遵循最长匹配原则，比如如果可以同时匹配api.google.com和google.com，
	//	              则最终会匹配到api.google.com
	//	@receiver x:
	//	@param domain: 要匹配的域名，比如 www.google.com
	//	@return interface{}: 匹配到的后缀所绑定的payload，如果没有匹配到的话则返回nil
	FindMatchDomainSuffixPayload(domain string) T

	// FindMatchDomainSuffixNode
	//
	//	@Description: 根据域名查询所匹配的后缀所对应的节点，会遵循最长匹配原则，比如如果可以同时匹配api.google.com和google.com，
	//	              则最终会匹配到api.google.com
	//	@receiver x:
	//	@param domain: 要匹配的域名，比如 www.google.com
	//	@return *SyncDomainSuffixTrieNode: 匹配到的后缀所对应的TreeNode，如果没有匹配到的话则返回nil
	FindMatchDomainSuffixNode(domain string) *DomainSuffixTrieNode[T]

	// AddDomainSuffix
	//
	//	@Description: 添加域名后缀追到字典树上，如果已经存在的话则会更新之前的值
	//	@receiver x:
	//	@param domainSuffix: 要添加的域名后缀
	//	@param Payload: 可以为这个后缀绑定一些payload，在后面拿域名匹配到这个后缀的时候可以一起获取到这个payload
	//	@return error: 如果添加后缀到树上时发生错误则返回error，否则返回nil
	AddDomainSuffix(domainSuffix string, payload T) error

	// GetPayload
	//
	//	@Description: 获取当前节点绑定的payload
	//	@receiver x:
	//	@return interface{}:
	GetPayload() T

	// SetPayload
	//
	//	@Description: 修改节点所绑定的payload，允许在节点创建之后修改其绑定的payload保存一些上下文信息之类的
	//	@receiver x:
	//	@param Payload:
	//	@return *SyncDomainSuffixTrieNode:
	SetPayload(payload T) T

	// GetChildNode
	//
	//	@Description: 按照字典值获取当前节点的孩子节点
	//	@receiver x:
	//	@param childValue:
	//	@return *SyncDomainSuffixTrieNode:
	//	@return bool:
	GetChildNode(childTrieValue string) (*DomainSuffixTrieNode[T], bool)

	// GetChildrenNodeMap
	//
	//	@Description: 返回当前节点的所有孩子节点，注意返回的是一个拷贝，树是不允许直接修改的
	//	@receiver x:
	//	@return map[string]SyncDomainSuffixTrieNode:
	GetChildrenNodeMap() map[string]*DomainSuffixTrieNode[T]

	// GetNodeTriePath
	//
	//	@Description: 获取当前节点对应的字典后缀路径，比如 com --> google --> api，如果当前节点是在api这个节点上，则此方法返回 "api.google.com"
	//	@receiver x:
	//	@return string:
	GetNodeTriePath() string

	// GetNodeTrieValue
	//
	//	@Description: 获取当前节点对应的字典值，比如 com --> google --> api，如果当前节点是在api这个节点上，则此方法返回 "api"
	//	@receiver x:
	//	@return string:
	GetNodeTrieValue() string
}
