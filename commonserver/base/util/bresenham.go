// Package util .
 
package util

// Bresenham .
// https://slg.zh.wikipedia.org/wiki/布雷森漢姆直線演算法
// 计算两点连线经过的格子
// eg: (0,0) -> (2,3) 得到：(0,0), (1,1), (1,2), (2,3)
/*
func Bresenham(startPos int32, endPos int32) []int32 {
	relatedPos := make([]int32, 0)

	startX, startY := Pos2XY(startPos)
	endX, endY := Pos2XY(endPos)

	var steep = false
	var swapped = false

	if AbsInt32(endY-startY) > AbsInt32(endX-startX) {
		steep = true
	}

	if steep == true {
		startX, startY = startY, startX
		endX, endY = endY, endX
	}
	if startX > endX {
		startX, endX = endX, startX
		startY, endY = endY, startY
		swapped = true
	}
	deltaX := endX - startX
	deltaY := AbsInt32(endY - startY)
	error := deltaX / 2
	var stepY int32
	y := startY
	if startY < endY {
		stepY = 1
	} else {
		stepY = -1
	}
	for x := startX; x <= endX; x++ {
		var c int32
		if steep {
			c = XY2Pos(y, x)
		} else {
			c = XY2Pos(x, y)
		}
		relatedPos = append(relatedPos, c)
		error = error - deltaY
		if error < 0 {
			y = y + stepY
			error = error + deltaX
		}
	}

	if swapped {
		relatedPos = reverse(relatedPos)
	}
	return relatedPos
}

func reverse(slice []int32) []int32 {
	if len(slice) == 0 || len(slice) == 1 {
		return slice
	}
	return append(reverse(slice[1:]), slice[0])
}
*/
