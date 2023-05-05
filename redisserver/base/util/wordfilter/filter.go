// Package wordfilter .

package wordfilter

import (
	"io/ioutil"
	"regexp"
	"strings"
)

// Filter 敏感词过滤器
type Filter struct {
	trie  *Trie
	noise *regexp.Regexp
}

// New 返回一个敏感词过滤器
func New() *Filter {
	return &Filter{
		trie:  NewTrie(),
		noise: regexp.MustCompile(`[\s&%$@*]+`),
	}
}

// UpdateNoisePattern 更新去噪模式
func (filter *Filter) UpdateNoisePattern(pattern string) {
	filter.noise = regexp.MustCompile(pattern)
}

// LoadWordDict 加载敏感词字典
func (filter *Filter) LoadWordDict(path string, ignoreNoise bool, ignoreWw bool) error {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	fContent := strings.ReplaceAll(string(content), "\r\n", "\n")
	words := strings.Split(fContent, "\n")
	if ignoreNoise {
		for i, _ := range words {
			words[i] = filter.RemoveNoise(words[i])
		}
	}
	filter.trie.Add(ignoreWw, words...)
	return nil
}

// AddWord 添加敏感词
func (filter *Filter) AddWord(ignoreWw bool, words ...string) {
	filter.trie.Add(ignoreWw, words...)
}

// Replace 和谐敏感词 不会去除噪音字符
func (filter *Filter) ReplaceIgnoreWw(text string, repl rune) string {
	return filter.trie.ReplaceIgnoreWw(text, repl)
}

// Replace 和谐敏感词 不会去除噪音字符
func (filter *Filter) Replace(text string, repl rune) string {
	return filter.trie.Replace(text, repl)
}

// FindIn 检测敏感词 会去除噪音字符
func (filter *Filter) FindIn(text string) string {
	text = filter.RemoveNoise(text)
	return filter.trie.FindIn(text)
}

// RemoveNoise 去除空格等噪音
func (filter *Filter) RemoveNoise(text string) string {
	return filter.noise.ReplaceAllString(text, "")
}
