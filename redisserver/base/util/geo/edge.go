// Package geo .
 
package geo

// Vertice 顶点
type Vertice struct {
	Index int32 // 顶点序号，唯一标识
	Coord Coord // 顶点坐标
}

// Edge 边
type Edge struct {
	WtCoord            Coord      // 重心点
	Vertices           [2]Vertice // 顶点数组，包含两个顶点0和1
	AdjacenctTriangles [2]int32   //相邻的两个三角形idx
	Inflects           [2]Vertice // 拐点数组
}

// CalMidCoord 计算边的终点
func (e *Edge) CalMidCoord() Coord {
	return CalMidCoord(e.Vertices[0].Coord, e.Vertices[1].Coord)
}

//// GenKey 生成边的序号
//func (e *Edge) GenKey() int32 {
//	return GenEdgeKey(e.Vertices[0].Index, e.Vertices[1].Index)
//}
//
//// GetEdgeKey 生成边的序号
//func GenEdgeKey(i, j int32) int32 {
//	if i < j {
//		return 10000*i + j
//	}
//	return 10000*j + i
//}
