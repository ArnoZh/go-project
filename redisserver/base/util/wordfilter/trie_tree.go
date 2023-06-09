// Package wordfilter .

package wordfilter

import "strings"

// Trie 整个Trie树.
type Trie struct {
	Root *Node
}

// Node Trie树上的一个节点.
type Node struct {
	isRootNode bool
	isPathEnd  bool
	Character  rune
	Children   map[rune]*Node
}

// NewTrie 创建字典数
func NewTrie() *Trie {
	return &Trie{
		Root: NewRootNode(0),
	}
}

// Add 初始化trie
func (tree *Trie) Add(ignoreWw bool, words ...string) {
	for _, word := range words {
		if ignoreWw {
			word = strings.ToLower(word)
		}
		tree.add(word)
	}
}

func (tree *Trie) add(word string) {
	var node = tree.Root
	var runes = []rune(word)
	for position := 0; position < len(runes); position++ {
		r := runes[position]
		if next, ok := node.Children[r]; ok {
			node = next
		} else {
			newNode := NewNode(r)
			node.Children[r] = newNode
			node = newNode
		}
		if position == len(runes)-1 {
			node.isPathEnd = true
		}
	}
}

// Replace 词语替换
func (tree *Trie) ReplaceIgnoreWw(text string, character rune) string {
	var node = tree.Root
	var parent = tree.Root
	var runes = []rune(strings.ToLower(text))
	var tmprunes = []rune(text)
	var wordLength = 0

	for position := 0; position < len(runes); position++ {
		r := runes[position]
		next, ok := node.Children[r]
		parent = node

		if !ok {
			if !node.IsRootNode() {
				if wordLength > 0 {
					if parent.IsPathEnd() {
						for i := position - wordLength; i < position; i++ {
							tmprunes[i] = character
						}
					}
				}
				position -= wordLength
			}
			node = tree.Root
			wordLength = 0
			continue
		}

		if position == len(runes)-1 && next.IsPathEnd() {
			for i := position - wordLength; i <= position; i++ {
				tmprunes[i] = character
			}
		}

		wordLength++
		node = next
	}
	return string(tmprunes)
}

// Replace 词语替换
func (tree *Trie) Replace(text string, character rune) string {
	var node = tree.Root
	var parent = tree.Root
	var runes = []rune(text)
	var wordLength = 0

	for position := 0; position < len(runes); position++ {
		r := runes[position]
		next, ok := node.Children[r]
		parent = node

		if !ok {
			if !node.IsRootNode() {
				if wordLength > 0 {
					if parent.IsPathEnd() {
						for i := position - wordLength; i < position; i++ {
							runes[i] = character
						}
					}
				}
				position -= wordLength
			}
			node = tree.Root
			wordLength = 0
			continue
		}

		if position == len(runes)-1 && next.IsPathEnd() {
			for i := position - wordLength; i <= position; i++ {
				runes[i] = character
			}
		}

		wordLength++
		node = next
	}
	return string(runes)
}

// FindIn 判断text中是否含有词库中的词并返回该词，返回""为不包含敏感词
func (tree *Trie) FindIn(text string) string {
	var node = tree.Root
	var parent = tree.Root
	var runes = []rune(text)
	var wordLength int
	for position := 0; position < len(runes); position++ {
		r := runes[position]
		next, ok := node.Children[r]
		parent = node
		if !ok {
			if !node.IsRootNode() {
				// @todo lzw 这一个判断应该是永远不会走到的 2021.11.25
				if wordLength > 0 && parent.IsPathEnd() {
					return string(runes[position-wordLength : position])
				}

				node = tree.Root
				position -= wordLength
			} else {
				node = tree.Root
			}
			wordLength = 0
			continue
		}
		if next.IsPathEnd() {
			return string(runes[position-wordLength : position+1])
		}

		wordLength++
		node = next
	}
	return ""
}

// NewNode 新建子节点
func NewNode(character rune) *Node {
	return &Node{
		Character: character,
		Children:  make(map[rune]*Node),
	}
}

// NewRootNode 新建根节点
func NewRootNode(character rune) *Node {
	return &Node{
		isRootNode: true,
		Character:  character,
		Children:   make(map[rune]*Node),
	}
}

// IsLeafNode 判断是否叶子节点
func (node *Node) IsLeafNode() bool {
	return len(node.Children) == 0
}

// IsRootNode 判断是否为根节点
func (node *Node) IsRootNode() bool {
	return node.isRootNode
}

// IsPathEnd 判断是否为某个路径的结束
func (node *Node) IsPathEnd() bool {
	return node.isPathEnd
}
