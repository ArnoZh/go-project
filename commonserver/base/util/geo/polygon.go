// Package geo .
 
package geo

// Polygon 多边形接口
type Polygon interface {
	IsCoordInside(p Coord) bool
	GetVectors() []Vector
	ToRect() (minX, minZ, maxX, maxZ int32)
	GetIndex() int32
	GetEdgeIDs() []int32
	GetEdgeMidCoords() []Coord
	GetVertices() []Vertice
}

// CrossProduct 计算叉积
func CrossProduct(p1, p2, p3 Vertice) int64 {
	ax := int64(p2.Coord.X - p1.Coord.X)
	ay := int64(p2.Coord.Z - p1.Coord.Z)
	bx := int64(p3.Coord.X - p2.Coord.X)
	by := int64(p3.Coord.Z - p2.Coord.Z)
	cp := ax*by - ay*bx

	return cp
}

// IsConvex 查看是否为凸多边形
func IsConvex(vertices []Vertice) bool {
	numPoints := len(vertices)
	negativeFlag := false
	positiveFlag := false
	for i := 0; i < numPoints; i++ {
		curvec := NewVector(vertices[i].Coord, vertices[(i+1)%numPoints].Coord)
		vec2next := NewVector(vertices[i].Coord, vertices[(i+2)%numPoints].Coord)
		if curvec.Cross(&vec2next) > 0 {
			positiveFlag = true
		} else {
			negativeFlag = true
		}
	}
	return positiveFlag != negativeFlag
}
