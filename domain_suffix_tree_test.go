package domain_suffix_trie

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDomainSuffixTrieNode_AddDomainSuffix(t *testing.T) {
	tire := NewDomainSuffixTrie[string]()
	err := tire.AddDomainSuffix("google.com", "谷歌")
	assert.Nil(t, err)
	payload := tire.FindMatchDomainSuffixPayload("asdasdasd.google.com")
	assert.Equal(t, "谷歌", payload)
}

func TestDomainSuffixTrieNode_FindMatchDomainSuffixNode(t *testing.T) {
	tire := NewDomainSuffixTrie[string]()
	err := tire.AddDomainSuffix("google.com", "谷歌")
	assert.Nil(t, err)
	node := tire.FindMatchDomainSuffixNode("asdasdasd.google.com")
	assert.NotNil(t, node)
	assert.Equal(t, "谷歌", node.GetPayload())
}

func TestDomainSuffixTrieNode_FindMatchDomainSuffixPayload(t *testing.T) {
	tire := NewDomainSuffixTrie[string]()
	err := tire.AddDomainSuffix("google.com", "谷歌")
	assert.Nil(t, err)
	payload := tire.FindMatchDomainSuffixPayload("asdasdasd.google.com")
	assert.Equal(t, "谷歌", payload)
}

func TestDomainSuffixTrieNode_GetChildNode(t *testing.T) {
	tire := NewDomainSuffixTrie[string]()
	err := tire.AddDomainSuffix("google.com", "谷歌")
	assert.Nil(t, err)
	err = tire.AddDomainSuffix("www.google.com", "谷歌web")
	assert.Nil(t, err)

	node := tire.FindMatchDomainSuffixNode("google.com")
	assert.NotNil(t, node)
	childNode, b := node.GetChildNode("www")
	assert.True(t, b)
	assert.NotNil(t, childNode)
	assert.Equal(t, "谷歌web", childNode.GetPayload())
}

func TestDomainSuffixTrieNode_GetChildrenNodeMap(t *testing.T) {
	tire := NewDomainSuffixTrie[string]()
	err := tire.AddDomainSuffix("google.com", "谷歌")
	assert.Nil(t, err)
	err = tire.AddDomainSuffix("www.google.com", "谷歌web")
	assert.Nil(t, err)

	node := tire.FindMatchDomainSuffixNode("google.com")
	assert.NotNil(t, node)
	childNodeMap := node.GetChildrenNodeMap()
	assert.NotNil(t, childNodeMap)
}

func TestDomainSuffixTrieNode_GetNodeTriePath(t *testing.T) {
	tire := NewDomainSuffixTrie[string]()
	err := tire.AddDomainSuffix("google.com", "谷歌")
	assert.Nil(t, err)
	err = tire.AddDomainSuffix("www.google.com", "谷歌web")
	assert.Nil(t, err)

	node := tire.FindMatchDomainSuffixNode("www.google.com")
	assert.NotNil(t, node)
	value := node.GetNodeTrieValue()
	assert.NotNil(t, "www.google.com", value)
}

func TestDomainSuffixTrieNode_GetNodeTrieValue(t *testing.T) {
	tire := NewDomainSuffixTrie[string]()
	err := tire.AddDomainSuffix("google.com", "谷歌")
	assert.Nil(t, err)
	err = tire.AddDomainSuffix("www.google.com", "谷歌web")
	assert.Nil(t, err)

	node := tire.FindMatchDomainSuffixNode("www.google.com")
	assert.NotNil(t, node)
	value := node.GetNodeTrieValue()
	assert.NotNil(t, "www", value)
}

func TestDomainSuffixTrieNode_GetPayload(t *testing.T) {
	tire := NewDomainSuffixTrie[string]()
	err := tire.AddDomainSuffix("google.com", "谷歌")
	assert.Nil(t, err)
	err = tire.AddDomainSuffix("www.google.com", "谷歌web")
	assert.Nil(t, err)

	node := tire.FindMatchDomainSuffixNode("www.google.com")
	assert.NotNil(t, node)
	value := node.GetPayload()
	assert.NotNil(t, "谷歌web", value)
}

func TestDomainSuffixTrieNode_SetPayload(t *testing.T) {

}

func TestDomainSuffixTrieNode_addChild(t *testing.T) {

}

func TestNewDomainSuffixTrie(t *testing.T) {

}
