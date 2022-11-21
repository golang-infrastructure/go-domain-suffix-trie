package domain_suffix_trie

import "errors"

// DomainSuffixIsEmptyError 错误：域名后缀是空的
var DomainSuffixIsEmptyError = errors.New("域名后缀是空的")
