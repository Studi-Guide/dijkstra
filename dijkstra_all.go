package dijkstra

//ShortestAll calculates all of the shortest paths from src to dest
func (g *Graph) ShortestAll(src, dest int) (BestPaths, error) {
	return g.evaluateAll(src, dest, true)
}

//LongestAll calculates all the longest paths from src to dest
func (g *Graph) LongestAll(src, dest int) (BestPaths, error) {
	return g.evaluateAll(src, dest, false)
}

func (g *Graph) evaluateAll(src, dest int, shortest bool) (BestPaths, error) {
	//Setup graph
	context := g.setup(shortest, src, -1)
	return g.postSetupEvaluateAll(src, dest, shortest, context)
}

func (g *Graph) postSetupEvaluateAll(src, dest int, shortest bool, context Context) (BestPaths, error) {
	var current *VertexResult
	oldCurrent := -1
	for context.visiting.Len() > 0 {
		//Visit the current lowest distanced Vertex
		current = context.visiting.PopOrdered()
		if oldCurrent == current.ID {
			continue
		}
		oldCurrent = current.ID

		currentVertex := g.Verticies[current.ID]
		//If the current distance is already worse than the best try another Vertex
		if shortest && current.distance > context.best {
			continue
		}
		for v, dist := range currentVertex.arcs {
			//If the arc has better access, than the current best, update the Vertex being touched
			if (shortest && current.distance+dist < context.VertexResults[v].distance) ||
				(!shortest && current.distance+dist > context.VertexResults[v].distance) ||
				(current.distance+dist == context.VertexResults[v].distance && !g.Verticies[v].containsBest(current.ID, context)) {
				//if g.Verticies[v].bestVertex == current.ID && g.Verticies[v].ID != dest {
				if currentVertex.containsBest(v, context) && g.Verticies[v].ID != dest {
					//also only do this if we aren't checkout out the best distance again
					//This seems familiar 8^)
					return BestPaths{}, newErrLoop(current.ID, v)
				}
				if current.distance+dist == context.VertexResults[v].distance {
					//At this point we know it's not in the list due to initial check
					context.VertexResults[v].bestVerticies = append(context.VertexResults[v].bestVerticies, current.ID)
				} else {
					context.VertexResults[v].distance = current.distance + dist
					context.VertexResults[v].bestVerticies = []int{current.ID}
				}
				if v == dest {
					context.visitedDest = true
					context.best = current.distance + dist
					continue
					//If this is the destination update best, so we can stop looking at
					// useless Verticies
				}
				//Push this updated Vertex into the list to be evaluated, pushes in
				// sorted form
				context.visiting.PushOrdered(context.VertexResults[v])
			}
		}
	}
	if !context.visitedDest {
		return BestPaths{}, ErrNoPath
	}
	return context.bestPaths(src, dest), nil
}

func (ctx *Context) bestPaths(src, dest int) BestPaths {
	paths := ctx.visitPath(src, dest, dest)
	best := BestPaths{}
	for indexPaths := range paths {
		for i, j := 0, len(paths[indexPaths])-1; i < j; i, j = i+1, j-1 {
			paths[indexPaths][i], paths[indexPaths][j] = paths[indexPaths][j], paths[indexPaths][i]
		}
		best = append(best, BestPath{ctx.VertexResults[dest].distance, paths[indexPaths]})
	}

	return best
}

func (ctx *Context) visitPath(src, dest, currentNode int) [][]int {
	if currentNode == src {
		return [][]int{
			[]int{currentNode},
		}
	}
	paths := [][]int{}
	for _, vertex := range ctx.VertexResults[currentNode].bestVerticies {
		sps := ctx.visitPath(src, dest, vertex)
		for i := range sps {
			paths = append(paths, append([]int{currentNode}, sps[i]...))
		}
	}
	return paths
}
