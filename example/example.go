package main

import (
	"fmt"
	"github.com/CC11001100/go-domain-suffix-tree/pkg/domain_suffix_trie"
)

func main() {

	// 调用 #NewDomainSuffixTrie 创建一颗后缀树
	tree := domain_suffix_trie.NewDomainSuffixTrie()

	// 将需要匹配的域名后缀依次调用 #AddDomainSuffix 添加到树上，添加的时候可以为后缀指定一个payload（使用集合A构建树）
	_ = tree.AddDomainSuffix("google.com", "谷歌主站子域名")
	_ = tree.AddDomainSuffix("map.google.com", "谷歌地图子域名")
	_ = tree.AddDomainSuffix("baidu.com", "百度主站子域名")
	_ = tree.AddDomainSuffix("jd.com", "京东子域名")

	// 需要查询的时候调用 #FindMatchDomainSuffixPayload 或者 #FindMatchDomainSuffixNode 查询，
	// 参数是一个完整的域名，会返回此域名匹配到的后缀在之前指定的payload（将集合B的每个元素依次在树上查询）
	fmt.Println(tree.FindMatchDomainSuffixPayload("test.google.com"))           // output: 谷歌主站子域名
	fmt.Println(tree.FindMatchDomainSuffixPayload("test.map.google.com"))       // output: 谷歌地图子域名
	fmt.Println(tree.FindMatchDomainSuffixNode("test.baidu.com").GetNodePath()) // output: baidu.com
	fmt.Println(tree.FindMatchDomainSuffixNode("test.jd.com").GetNodeValue())   // output: jd

}
