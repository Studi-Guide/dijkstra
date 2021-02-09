package dijkstra

type VertexResult struct {
	bestVerticies []int
	distance      int64
	ID            int
}

type Context struct {
	visiting      dijkstraList
	VertexResults map[int]*VertexResult
	visitedDest   bool
	best          int64
}

//SetDefaults sets the distance and best node to that specified
func (ctx *Context) setDefaults(Distance int64, BestNode int, g *Graph) {
	for i := range g.Verticies {
		ctx.VertexResults[i] = &VertexResult{
			bestVerticies: []int{BestNode},
			distance:      Distance,
			ID:            i,
		}
	}
}
