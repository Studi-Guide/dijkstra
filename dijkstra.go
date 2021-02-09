package dijkstra

import (
	"math"
)

//Shortest calculates the shortest path from src to dest
func (g *Graph) Shortest(src, dest int) (BestPath, error) {
	return g.evaluate(src, dest, true)
}

//Longest calculates the longest path from src to dest
func (g *Graph) Longest(src, dest int) (BestPath, error) {
	return g.evaluate(src, dest, false)
}

func (g *Graph) setup(shortest bool, src int, list int) Context {
	//-1 auto list
	//Get a new list regardless

	var context = Context{
		visiting:      nil,
		VertexResults: map[int]*VertexResult{},
		visitedDest:   false,
		best:          0,
	}

	if list >= 0 {
		context.visiting = g.forceList(list)
	} else if shortest {
		context.visiting = g.forceList(-1)
	} else {
		context.visiting = g.forceList(-2)
	}
	//Reset state
	context.visitedDest = false
	//Reset the best current value (worst so it gets overwritten)
	// and set the defaults *almost* as bad
	// set all best verticies to -1 (unused)
	if shortest {
		context.setDefaults(int64(math.MaxInt64)-2, -1, g)
		context.best = int64(math.MaxInt64)
	} else {
		context.setDefaults(int64(math.MinInt64)+2, -1, g)
		context.best = int64(math.MinInt64)
	}

	//Set the distance of initial vertex 0
	context.VertexResults[src].distance = 0

	//Add the source vertex to the list
	context.visiting.PushOrdered(context.VertexResults[src])
	return context
}

func (g *Graph) forceList(i int) dijkstraList {
	//-2 long auto
	//-1 short auto
	//0 short pq
	//1 long pq
	//2 short ll
	//3 long ll
	switch i {
	case -2:
		if len(g.Verticies) < 800 {
			return g.forceList(2)
		} else {
			return g.forceList(0)
		}
		break
	case -1:
		if len(g.Verticies) < 800 {
			return g.forceList(3)
		} else {
			return g.forceList(1)
		}
		break
	case 0:
		return priorityQueueNewShort()
		break
	case 1:
		return priorityQueueNewLong()
		break
	case 2:
		return linkedListNewShort()
		break
	case 3:
		return linkedListNewLong()
		break
	default:
		panic(i)
	}

	return nil
}

func (g *Graph) bestPath(src, dest int, ctx Context) BestPath {
	var path []int
	for c := g.Verticies[dest]; c.ID != src; c = g.Verticies[ctx.VertexResults[c.ID].bestVerticies[0]] {
		path = append(path, c.ID)
	}
	path = append(path, src)
	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}
	return BestPath{ctx.VertexResults[dest].distance, path}
}

func (g *Graph) evaluate(src, dest int, shortest bool) (BestPath, error) {
	//Setup graph
	graphContext := g.setup(shortest, src, -1)
	return g.postSetupEvaluate(src, dest, shortest, graphContext)
}

func (g *Graph) postSetupEvaluate(src, dest int, shortest bool, graphContext Context) (BestPath, error) {
	var currentVertexResult *VertexResult
	oldCurrent := -1
	for graphContext.visiting.Len() > 0 {
		//Visit the currentVertexResult lowest distanced Vertex
		currentVertexResult = graphContext.visiting.PopOrdered()
		if oldCurrent == currentVertexResult.ID {
			continue
		}

		oldCurrent = currentVertexResult.ID
		currentVertex := g.Verticies[currentVertexResult.ID]
		//If the currentVertexResult distance is already worse than the best try another Vertex
		if shortest && graphContext.VertexResults[oldCurrent].distance >= graphContext.best {
			continue
		}
		for v, dist := range currentVertex.arcs {
			//If the arc has better access, than the currentVertexResult best, update the Vertex being touched
			if (shortest && currentVertexResult.distance+dist < graphContext.VertexResults[v].distance) ||
				(!shortest && currentVertexResult.distance+dist > graphContext.VertexResults[v].distance) {
				if currentVertexResult.bestVerticies[0] == v && g.Verticies[v].ID != dest {
					//also only do this if we aren't checkout out the best distance again
					//This seems familiar 8^)
					return BestPath{}, newErrLoop(currentVertex.ID, v)
				}
				graphContext.VertexResults[v].distance = currentVertexResult.distance + dist
				graphContext.VertexResults[v].bestVerticies[0] = currentVertex.ID
				if v == dest {
					//If this is the destination update best, so we can stop looking at
					// useless Verticies
					graphContext.best = currentVertexResult.distance + dist
					graphContext.visitedDest = true
					continue // Do not push if dest
				}
				//Push this updated Vertex into the list to be evaluated, pushes in
				// sorted form
				graphContext.visiting.PushOrdered(graphContext.VertexResults[v])
			}
		}
	}
	return g.finally(src, dest, graphContext)
}

func (g *Graph) finally(src, dest int, ctx Context) (BestPath, error) {
	if !ctx.visitedDest {
		return BestPath{}, ErrNoPath
	}
	return g.bestPath(src, dest, ctx), nil
}

//BestPath contains the solution of the most optimal path
type BestPath struct {
	Distance int64
	Path     []int
}

//BestPaths contains the list of best solutions
type BestPaths []BestPath
