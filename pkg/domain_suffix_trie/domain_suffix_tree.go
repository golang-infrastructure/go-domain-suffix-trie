package domain_suffix_trie

import (
	"errors"
	"strings"
	"sync"
)

// DomainSuffixTrieNode
//  @Description: 域名后缀树，用来做域名后缀匹配查询，这个结构是线程安全的
//  @thread-safe: 是线程安全的
type DomainSuffixTrieNode struct {

	// 用于保证线程安全操作，每个节点会有一个锁，对接点线程不安全的字段操作、访问之前都要先获取锁
	lock sync.RWMutex

	// value
	//  @Description: 此节点的值，用来存储域名后缀使用.分割后的一段，每一段是一个节点
	//  @thread-safe: 这个节点在创建的时候确定，再之后就不会被改变，所以这个字段是线程安全的
	value string

	// parent
	//  @Description: 此节点的父节点
	//  @thread-safe: 创建Node的时候就会初始化parent字段，同时后边不会再改变，因此这个字段是线程安全的
	parent *DomainSuffixTrieNode

	// childrenNodeMap
	//  @Description: 此节点的孩子的值
	//  @thread-unsafe: 因为是可以动态的往树上添加后缀的，因此孩子节点也是会动态改变的，
	//                  因此此字段不是线程安全的，要加锁
	childrenNodeMap map[string]*DomainSuffixTrieNode

	// payload
	//  @Description: 关联到从根路径到子节点的这条后缀路径上的一些额外信息，
	//                可以给某个域名后缀指定一些附加信息，当匹配的时候就能把它取回来
	//  @thread-unsafe: 在节点创建之后payload允许被修改，因此是线程不安全的，对此字段的操作要加锁
	payload interface{}
}

// NewDomainSuffixTrie
//  @Description: 创建一颗新的域名后缀树，将这颗树的根节点返回
//  @return *DomainSuffixTrieNode
func NewDomainSuffixTrie() *DomainSuffixTrieNode {
	return &DomainSuffixTrieNode{
		lock: sync.RWMutex{},
		// 根节点为空
		value: "",
		// 根节点没有父节点
		parent:          nil,
		childrenNodeMap: make(map[string]*DomainSuffixTrieNode),
		payload:         nil,
	}
}

// GetNodeValue
//  @Description: 获取当前节点对应的值，比如 com --> google --> api，如果当前节点是在api这个节点上，则此方法返回 "api"
//  @receiver x:
//  @return string:
func (x *DomainSuffixTrieNode) GetNodeValue() string {
	return x.value
}

// GetNodePath
//  @Description: 获取当前节点对应的后缀路径，比如 com --> google --> api，如果当前节点是在api这个节点上，则此方法返回 "api.google.com"
//  @receiver x:
//  @return string:
func (x *DomainSuffixTrieNode) GetNodePath() string {
	valueSlice := make([]string, 0)
	currentNode := x
	for currentNode != nil {
		if currentNode.value != "" {
			valueSlice = append(valueSlice, currentNode.value)
		}
		currentNode = currentNode.parent
	}
	return strings.Join(valueSlice, ".")
}

// GetParentNode
//  @Description: 获取当前节点的父节点
//  @receiver x:
//  @return *DomainSuffixTrieNode:
func (x *DomainSuffixTrieNode) GetParentNode() *DomainSuffixTrieNode {
	return x.parent
}

// GetChildrenNodeMap
//  @Description: 返回当前节点的所有孩子节点，注意返回的是一个拷贝，树是不允许直接修改的
//  @receiver x:
//  @return map[string]*DomainSuffixTrieNode:
func (x *DomainSuffixTrieNode) GetChildrenNodeMap() map[string]*DomainSuffixTrieNode {
	x.lock.Lock()
	defer x.lock.Unlock()

	childrenNodeMap := make(map[string]*DomainSuffixTrieNode)
	for value, node := range x.childrenNodeMap {
		childrenNodeMap[value] = node
	}
	return childrenNodeMap
}

// GetChild
//  @Description: 获取当前节点的孩子节点
//  @receiver x:
//  @param childValue:
//  @return *DomainSuffixTrieNode:
//  @return bool:
func (x *DomainSuffixTrieNode) GetChild(childValue string) (*DomainSuffixTrieNode, bool) {
	x.lock.Lock()
	defer x.lock.Unlock()

	childNode, exists := x.childrenNodeMap[childValue]
	return childNode, exists
}

// addChild
//  @Description: 为当前节点添加孩子节点
//  @receiver x:
//  @param childNode:
//  @return *DomainSuffixTrieNode:
func (x *DomainSuffixTrieNode) addChild(childNode *DomainSuffixTrieNode) *DomainSuffixTrieNode {
	x.lock.Lock()
	defer x.lock.Unlock()

	x.childrenNodeMap[childNode.value] = childNode
	return x
}

// SetPayload
//  @Description: 修改节点所绑定的payload，允许在节点创建之后修改其绑定的payload
//  @receiver x:
//  @param payload:
//  @return *DomainSuffixTrieNode:
func (x *DomainSuffixTrieNode) SetPayload(payload interface{}) *DomainSuffixTrieNode {
	x.lock.Lock()
	defer x.lock.Unlock()

	x.payload = payload
	return x
}

// GetPayload
//  @Description: 获取当前节点绑定的payload
//  @receiver x:
//  @return interface{}:
func (x *DomainSuffixTrieNode) GetPayload() interface{} {
	x.lock.Lock()
	defer x.lock.Unlock()

	return x.payload
}

// setValue
//  @Description: 设置节点的值
//  @receiver x:
//  @param value:
//  @return *DomainSuffixTrieNode:
func (x *DomainSuffixTrieNode) setValue(value string) *DomainSuffixTrieNode {
	x.value = value
	return x
}

// setParent
//  @Description: 设置父节点的值
//  @receiver x:
//  @param parent:
//  @return *DomainSuffixTrieNode:
func (x *DomainSuffixTrieNode) setParent(parent *DomainSuffixTrieNode) *DomainSuffixTrieNode {
	x.parent = parent
	return x
}

// --------------------------------------------------------------------------------------------------------------------

// DomainSuffixIsEmptyError 错误：域名后缀是空的
var DomainSuffixIsEmptyError = errors.New("域名后缀是空的")

// AddDomainSuffix
//  @Description: 添加域名后缀追到字典树上
//  @receiver x:
//  @param domainSuffix: 要添加的域名后缀
//  @param payload: 可以为这个后缀绑定一些payload，在后面拿域名匹配到这个后缀的时候可以一起获取到这个payload
//  @return error: 如果添加后缀到树上时发生错误则返回error，否则返回nil
func (x *DomainSuffixTrieNode) AddDomainSuffix(domainSuffix string, payload interface{}) error {

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
		if node, exists := currentNode.GetChild(v); exists {
			currentNode = node
		} else {
			node := NewDomainSuffixTrie().setValue(v).setParent(currentNode)
			currentNode.addChild(node)
			currentNode = node
		}
	}
	// 都放完了把对应的信息放在叶子节点上
	//if currentNode.payload != nil {
	//	return DomainSuffixRepetitionError
	//}
	// 允许后来的payload把之前的payload给覆盖掉
	currentNode.SetPayload(payload)

	return nil
}

// FindMatchDomainSuffixNode
//  @Description: 根据域名查询所匹配的后缀所对应的节点，会遵循最长匹配原则，比如如果可以同时匹配api.google.com和google.com，
//                则最终会匹配到api.google.com
//  @receiver x:
//  @param domain: 要匹配的域名，比如 www.google.com
//  @return *DomainSuffixTrieNode: 匹配到的后缀所对应的TreeNode，如果没有匹配到的话则返回nil
func (x *DomainSuffixTrieNode) FindMatchDomainSuffixNode(domain string) *DomainSuffixTrieNode {
	// 对输入的域名切割为不同的级别
	domainLevelValueSlice := strings.Split(domain, ".")
	// 然后倒着去字典树中匹配，采用最长匹配策略
	currentNode := x // x is root
	for level := len(domainLevelValueSlice) - 1; level >= 0; level-- {
		v := domainLevelValueSlice[level]
		node, exists := currentNode.GetChild(v)
		if exists {
			currentNode = node
		} else {
			return currentNode
		}
	}
	return currentNode
}

// FindMatchDomainSuffixPayload
//  @Description: 根据域名查询所匹配的后缀所对应的payload，会遵循最长匹配原则，比如如果可以同时匹配api.google.com和google.com，
//                则最终会匹配到api.google.com
//  @receiver x:
//  @param domain: 要匹配的域名，比如 www.google.com
//  @return interface{}: 匹配到的后缀所绑定的payload，如果没有匹配到的话则返回nil
func (x *DomainSuffixTrieNode) FindMatchDomainSuffixPayload(domain string) interface{} {
	node := x.FindMatchDomainSuffixNode(domain)
	if node != nil {
		return node.GetPayload()
	} else {
		return nil
	}
}