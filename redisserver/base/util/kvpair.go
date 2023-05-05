// Package util .
 
package util

// StringKVpair K为string的KVPair
// +k8s:deepcopy-gen=true
type StringKVpair struct {
	K string
	V int32
}

// StringKVpairFind 根据key查找索引
func StringKVpairFind(a []*StringKVpair, k string) int {
	for i, v := range a {
		if v.K == k {
			return i
		}
	}
	return -1
}

// StringKVpairFindOrAdd 从数组根据key查找,没有则添加
func StringKVpairFindOrAdd(a *[]*StringKVpair, k string) *StringKVpair {
	for _, v := range *a {
		if v.K == k {
			return v
		}
	}
	v := &StringKVpair{k, 0}
	*a = append(*a, v)
	return v
}

// Int32KVpair K为int32的KVPair
// +k8s:deepcopy-gen=true
type Int32KVpair struct {
	K int32
	V int32
}

// Int32KVpairFind 根据key查找索引
func Int32KVpairFind(a []*Int32KVpair, k int32) int {
	for i, v := range a {
		if v.K == k {
			return i
		}
	}
	return -1
}

// Int32KVpairFindOrAdd 在Int32 KVPair List中查找或插入KVPair
func Int32KVpairFindOrAdd(a *[]*Int32KVpair, k int32) *Int32KVpair {
	for _, v := range *a {
		if v.K == k {
			return v
		}
	}
	v := &Int32KVpair{K: k, V: 0}
	*a = append(*a, v)
	return v
}

// Int32KVPairArray KVPair 数组
// +k8s:deepcopy-gen=true
type Int32KVPairArray struct {
	Array []*Int32KVpair
}

// NewInt32KVPairArray 创建KVPair数组
func NewInt32KVPairArray(size, capacity int32) *Int32KVPairArray {
	return &Int32KVPairArray{
		Array: make([]*Int32KVpair, size, capacity),
	}
}

// Insert 插入KV对
func (kvp *Int32KVPairArray) Insert(k, v int32) {
	idx := Int32KVpairFind(kvp.Array, k)
	if idx != -1 {
		kvp.Array[idx].V = v
		return
	}
	kvp.Array = append(kvp.Array, &Int32KVpair{K: k, V: v})
}

// Update 更新KV对
func (kvp *Int32KVPairArray) Update(k, v int32) bool {
	idx := Int32KVpairFind(kvp.Array, k)
	if idx == -1 {
		return false
	}
	kvp.Array[idx].V = v
	return true
}

// Get 获取KV对
func (kvp *Int32KVPairArray) Get(k int32) (v int32, find bool) {
	idx := Int32KVpairFind(kvp.Array, k)
	if idx == -1 {
		return 0, false
	}
	return kvp.Array[idx].V, true
}

// Clone 拷贝KVPair数组
func (kvp *Int32KVPairArray) Clone(other *Int32KVPairArray) {
	for _, ele := range other.Array {
		kvp.Array = append(kvp.Array, &Int32KVpair{K: ele.K, V: ele.V})
	}
}

// Int64KVpair K为int64的KVPair
// +k8s:deepcopy-gen=true
type Int64KVpair struct {
	K int64
	V int64
}

// Int64KVpairFind 根据key查找索引
func Int64KVpairFind(a []*Int64KVpair, k int64) int {
	for i, v := range a {
		if v.K == k {
			return i
		}
	}
	return -1
}

// Int64KVpairFindOrAdd Int64 KVPair List中查找或插入KVPair
func Int64KVpairFindOrAdd(a *[]*Int64KVpair, k int64) *Int64KVpair {
	for _, v := range *a {
		if v.K == k {
			return v
		}
	}
	v := &Int64KVpair{K: k, V: 0}
	*a = append(*a, v)
	return v
}
