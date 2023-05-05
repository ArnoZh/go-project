// Package bitmap .

package bitmap

var (
	// 只有第i位为1
	tA = [8]byte{1, 2, 4, 8, 16, 32, 64, 128}
	// 只有第i位为0
	// 11111110
	// 11111101
	// 11111011
	// 11110111
	// 11101111
	// 11011111
	// 10111111
	// 01111111
	tB = [8]byte{254, 253, 251, 247, 239, 223, 191, 127}
)

// Bitmap 位图定义，可以通过bitmap := Bitmap(data)直接基于已有[]byte创建
type Bitmap struct {
	Map    []byte
	Num    int64 // 位图中有多少个1
	MaxNum int64 //最多有多少个，超过就把Map设置成nil
}

// New 创建大小为l的一个位图，位图的实际大小为8的整数倍
func New(l int32) *Bitmap {
	remainder := l % 8
	if remainder != 0 { //不足1个byte也占1个
		remainder = 1
	}
	return &Bitmap{
		Map:    make([]byte, l/8+remainder),
		Num:    0,
		MaxNum: int64(l),
	}
}

// Fix 用于计算老版本的MaxNum,奖励不是8的倍数的也会重置成8的倍数 某个版本删除吧
func (b *Bitmap) Fix() {
	if b.MaxNum > 0 {
		return
	}
	b.MaxNum = int64(b.Len())
	b.CheckMax()
}

// CheckMax check是否都填满了
func (b *Bitmap) CheckMax() bool {
	if b.Num >= b.MaxNum {
		//todo Ryan 检查操作为啥给要置nil？ 而且单独也不能单独置空map，需要的时候自己加reset接口，暂时注掉
		//b.Map = nil
		return true
	}
	return false
}

// FromMaxToNormal ..目前是没有Set 0 的操作
func (b *Bitmap) FromMaxToNormal() {
	b.Map = GenAllTrueByMaxNum(int(b.MaxNum))
}

// GenAllTrueByMaxNum  按照len全部设置成为1
func GenAllTrueByMaxNum(maxNum int) []byte {
	minLen := maxNum / 8
	remainder := maxNum % 8
	if remainder != 0 {
		remainder = 1
	}
	result := make([]byte, minLen+remainder)
	for i := 0; i < minLen; i++ {
		result[i] = 255
	}
	if maxNum%8 != 0 {
		for i := 0; i < maxNum%8; i++ {
			result[minLen] += tA[i]
		}
	}
	return result
}

// Len 返回位图大小，位图的大小总是向上取到8的倍数，而不是创建的时候指定的大小
func (b *Bitmap) Len() int32 {
	return int32(len(b.Map) * 8)
}

// Get 获取位图指定位，不检查越界
func (b *Bitmap) Get(i int32) bool {
	if b.CheckMax() {
		return true
	}
	m := b.Map
	return m[i/8]&tA[i%8] != 0
}

// Set 设置位图指定位，不检查越界
func (b *Bitmap) Set(i int32, v bool) {
	m := b.Map
	index := i / 8
	bit := i % 8
	if v {
		if b.CheckMax() {
			return
		}
		if m[index]&tA[bit] == 0 {
			b.Num++
		}
		m[index] = m[index] | tA[bit]
		b.CheckMax()
	} else {
		if b.CheckMax() {
			b.FromMaxToNormal()
			m = b.Map
		}
		if m[index]&tA[bit] != 0 {
			b.Num--
		}
		m[index] = m[index] & tB[bit]
	}
}

// DeepCopy 返回位图的值，docopy表示是否需要拷贝
func (b *Bitmap) DeepCopy() *Bitmap {
	return &Bitmap{
		Map:    append(b.Map[:0:0], b.Map...),
		Num:    b.Num,
		MaxNum: b.MaxNum,
	}
}

// ResetNum 重算位图Num值
func (b *Bitmap) ResetNum() {
	b.Num = 0
	for _, elem := range b.Map {
		for {
			if elem <= 0 {
				break
			}
			b.Num++
			elem &= elem - 1
		}
	}
}
