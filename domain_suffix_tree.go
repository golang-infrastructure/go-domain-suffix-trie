package domain_suffix_trie

import (
	"strings"
)

// ------------------------------------------------- DomainSuffixTrieNode ----------------------------------------------

// DomainSuffixTrieNode
//
//	@Description: 域名后缀树，用来做域名后缀匹配查询，这个结构是线程安全的
//	@thread-safe: 是线程安全的
type DomainSuffixTrieNode[T any] struct {

	// TrieValue
	//  @Description: 此节点的值，用来存储域名后缀使用.分割后的一段，每一段是一个节点
	//  @thread-safe: 这个节点在创建的时候确定，再之后就不会被改变，所以这个字段是线程安全的
	TrieValue string

	// ParentNode
	//  @Description: 此节点的父节点
	//  @thread-safe: 创建Node的时候就会初始化parent字段，同时后边不会再改变，因此这个字段是线程安全的
	ParentNode *DomainSuffixTrieNode[T]

	// ChildrenNodeMap
	//  @Description: 此节点的孩子的值
	//  @thread-safe: 因为是可以动态的往树上添加后缀的，因此孩子节点也是会动态改变的，
	ChildrenNodeMap map[string]*DomainSuffixTrieNode[T]

	// Payload
	//  @Description: 关联到从根路径到子节点的这条后缀路径上的一些额外信息，
	//                可以给某个域名后缀指定一些附加信息，当匹配的时候就能把它取回来
	Payload T
}

var _ DomainSuffixTrieInterface[any] = &DomainSuffixTrieNode[any]{}

// NewDomainSuffixTrie
//
//	@Description: 创建一颗新的域名后缀树，将这颗树的根节点返回
//	@return *SyncDomainSuffixTrieNode
func NewDomainSuffixTrie[T any]() *DomainSuffixTrieNode[T] {
	return &DomainSuffixTrieNode[T]{
		// 根节点为空
		TrieValue: "",
		// 根节点没有父节点
		ChildrenNodeMap: make(map[string]*DomainSuffixTrieNode[T]),
	}
}

// GetNodeTrieValue
//
//	@Description: 获取当前节点对应的值，比如 com --> google --> api，如果当前节点是在api这个节点上，则此方法返回 "api"
//	@receiver x:
//	@return string:
func (x *DomainSuffixTrieNode[T]) GetNodeTrieValue() string {
	return x.TrieValue
}

// GetNodeTriePath
//
//	@Description: 获取当前节点对应的后缀路径，比如 com --> google --> api，如果当前节点是在api这个节点上，则此方法返回 "api.google.com"
//	@receiver x:
//	@return string:
func (x *DomainSuffixTrieNode[T]) GetNodeTriePath() string {
	valueSlice := make([]string, 0)
	currentNode := x
	for currentNode != nil && currentNode.TrieValue != "" {
		valueSlice = append(valueSlice, currentNode.TrieValue)
		currentNode = currentNode.ParentNode
	}
	return strings.Join(valueSlice, ".")
}

// GetChildrenNodeMap
//
//	@Description: 返回当前节点的所有孩子节点，注意返回的是一个拷贝，树是不允许直接修改的
//	@receiver x:
//	@return map[string]SyncDomainSuffixTrieNode:
func (x *DomainSuffixTrieNode[T]) GetChildrenNodeMap() map[string]*DomainSuffixTrieNode[T] {
	copyChildrenNodeMap := make(map[string]*DomainSuffixTrieNode[T])
	for key, value := range x.ChildrenNodeMap {
		copyChildrenNodeMap[key] = value
	}
	return copyChildrenNodeMap
}

// GetChildNode
//
//	@Description: 获取当前节点的孩子节点
//	@receiver x:
//	@param childValue:
//	@return *SyncDomainSuffixTrieNode:
//	@return bool:
func (x *DomainSuffixTrieNode[T]) GetChildNode(childTrieValue string) (*DomainSuffixTrieNode[T], bool) {
	childNode, exists := x.ChildrenNodeMap[childTrieValue]
	return childNode, exists
}

// addChild
//
//	@Description: 为当前节点添加孩子节点
//	@receiver x:
//	@param childNode:
//	@return *SyncDomainSuffixTrieNode: 如果要设置的key已经存在的话，会返回原来的key
func (x *DomainSuffixTrieNode[T]) addChild(childNode *DomainSuffixTrieNode[T]) *DomainSuffixTrieNode[T] {
	if childNode == nil {
		return nil
	}
	lastChildNode := x.ChildrenNodeMap[childNode.TrieValue]
	x.ChildrenNodeMap[childNode.TrieValue] = childNode
	return lastChildNode
}

// SetPayload
//
//	@Description: 修改节点所绑定的payload，允许在节点创建之后修改其绑定的payload
//	@receiver x:
//	@param Payload:
//	@return *SyncDomainSuffixTrieNode:
func (x *DomainSuffixTrieNode[T]) SetPayload(payload T) T {
	lastPayload := x.Payload
	x.Payload = payload
	return lastPayload
}

// GetPayload
//
//	@Description: 获取当前节点绑定的payload
//	@receiver x:
//	@return interface{}:
func (x *DomainSuffixTrieNode[T]) GetPayload() T {
	return x.Payload
}

// AddDomainSuffix
//
//	@Description: 添加域名后缀追到字典树上
//	@receiver x:
//	@param domainSuffix: 要添加的域名后缀
//	@param Payload: 可以为这个后缀绑定一些payload，在后面拿域名匹配到这个后缀的时候可以一起获取到这个payload
//	@return error: 如果添加后缀到树上时发生错误则返回error，否则返回nil
func (x *DomainSuffixTrieNode[T]) AddDomainSuffix(domainSuffix string, payload T) error {

	// 必须是合法的后缀域名
	if domainSuffix == "" {
		return DomainSuffixIsEmptyError
	}

	// 然后就是将每个级别对应上往树上放就可以了，放的时候是倒序放的
	domainLevelValueSlice := strings.Split(domainSuffix, ".")
	currentNode := x
	for index := len(domainLevelValueSlice) - 1; index >= 0; index-- {
		v := domainLevelValueSlice[index]

		// 要把v插入到currentNode的孩子节点上，先看看之前是不是已经存在过
		if node, exists := currentNode.ChildrenNodeMap[v]; exists {
			currentNode = node
		} else {
			node := &DomainSuffixTrieNode[T]{
				ChildrenNodeMap: make(map[string]*DomainSuffixTrieNode[T]),
				TrieValue:       v,
				ParentNode:      currentNode,
			}
			currentNode.addChild(node)
			currentNode = node
		}
	}
	// 都放完了把对应的信息放在叶子节点上
	//if currentNode.Payload != nil {
	//	return DomainSuffixRepetitionError
	//}
	// 允许后来的payload把之前的payload给覆盖掉
	currentNode.SetPayload(payload)

	return nil
}

// FindMatchDomainSuffixNode
//
//	@Description: 根据域名查询所匹配的后缀所对应的节点，会遵循最长匹配原则，比如如果可以同时匹配api.google.com和google.com，
//	              则最终会匹配到api.google.com
//	@receiver x:
//	@param domain: 要匹配的域名，比如 www.google.com
//	@return *SyncDomainSuffixTrieNode: 匹配到的后缀所对应的TreeNode，如果没有匹配到的话则返回nil
func (x *DomainSuffixTrieNode[T]) FindMatchDomainSuffixNode(domain string) *DomainSuffixTrieNode[T] {
	// 对输入的域名切割为不同的级别
	domainLevelValueSlice := strings.Split(domain, ".")
	// 然后倒着去字典树中匹配，采用最长匹配策略
	currentNode := x // x is root
	for level := len(domainLevelValueSlice) - 1; level >= 0; level-- {
		v := domainLevelValueSlice[level]
		node, exists := currentNode.ChildrenNodeMap[v]
		if exists {
			currentNode = node
		} else {
			return currentNode
		}
	}
	return currentNode
}

// FindMatchDomainSuffixPayload
//
//	@Description: 根据域名查询所匹配的后缀所对应的payload，会遵循最长匹配原则，比如如果可以同时匹配api.google.com和google.com，
//	              则最终会匹配到api.google.com
//	@receiver x:
//	@param domain: 要匹配的域名，比如 www.google.com
//	@return interface{}: 匹配到的后缀所绑定的payload，如果没有匹配到的话则返回nil
func (x *DomainSuffixTrieNode[T]) FindMatchDomainSuffixPayload(domain string) T {
	node := x.FindMatchDomainSuffixNode(domain)
	if node != nil {
		return node.GetPayload()
	} else {
		var zero T
		return zero
	}
}
