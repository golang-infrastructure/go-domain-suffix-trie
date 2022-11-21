package domain_suffix_trie

import (
	"strings"
	"sync"
)

// SyncDomainSuffixTrieNode 线程安全的实现
type SyncDomainSuffixTrieNode[T any] struct {
	lock sync.RWMutex
	node *DomainSuffixTrieNode[T]
}

var _ DomainSuffixTrieInterface[any] = &SyncDomainSuffixTrieNode[any]{}

func NewSyncDomainSuffixTrie[T any]() *SyncDomainSuffixTrieNode[T] {
	return &SyncDomainSuffixTrieNode[T]{
		lock: sync.RWMutex{},
		node: NewDomainSuffixTrie[T](),
	}
}

func (x *SyncDomainSuffixTrieNode[T]) FindMatchDomainSuffixPayload(domain string) T {
	x.lock.RLock()
	defer x.lock.RUnlock()
	return x.node.FindMatchDomainSuffixPayload(domain)
}

func (x *SyncDomainSuffixTrieNode[T]) FindMatchDomainSuffixNode(domain string) *DomainSuffixTrieNode[T] {
	x.lock.RLock()
	defer x.lock.RUnlock()
	return x.node.FindMatchDomainSuffixNode(domain)
}

func (x *SyncDomainSuffixTrieNode[T]) AddDomainSuffix(domainSuffix string, payload T) error {
	x.lock.Lock()
	defer x.lock.Unlock()
	return x.node.AddDomainSuffix(domainSuffix, payload)
}

func (x *SyncDomainSuffixTrieNode[T]) GetPayload() T {
	x.lock.RLock()
	defer x.lock.RUnlock()
	return x.node.GetPayload()
}

func (x *SyncDomainSuffixTrieNode[T]) SetPayload(payload T) T {
	x.lock.Lock()
	defer x.lock.Unlock()
	return x.node.SetPayload(payload)
}

func (x *SyncDomainSuffixTrieNode[T]) GetChildNode(childTrieValue string) (*DomainSuffixTrieNode[T], bool) {
	x.lock.RLock()
	defer x.lock.RUnlock()
	return x.node.GetChildNode(childTrieValue)
}

func (x *SyncDomainSuffixTrieNode[T]) GetChildrenNodeMap() map[string]*DomainSuffixTrieNode[T] {
	x.lock.RLock()
	defer x.lock.RUnlock()
	return x.node.ChildrenNodeMap
}

func (x *SyncDomainSuffixTrieNode[T]) GetNodeTriePath() string {
	valueSlice := make([]string, 0)
	currentNode := x.node
	for currentNode != nil && currentNode.TrieValue != "" {
		valueSlice = append(valueSlice, currentNode.GetNodeTrieValue())
		currentNode = currentNode.ParentNode
	}
	return strings.Join(valueSlice, ".")
}

func (x *SyncDomainSuffixTrieNode[T]) GetNodeTrieValue() string {
	x.lock.RLock()
	defer x.lock.RUnlock()
	return x.node.TrieValue
}
