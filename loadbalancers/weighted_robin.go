package loadbalancers

type Tuple struct {
	a, b int64
}

type WeightedBalancer struct {
	currIndex            *int64
	nodeIndexWithWeights map[int64]Tuple
}

func NewWeightedBalancer(indexsWithWeight map[int64]int64) *WeightedBalancer {
	index := int64(0)

	wb := &WeightedBalancer{
		currIndex:            &index,
		nodeIndexWithWeights: make(map[int64]Tuple), // its a tuple of current weight & effective weight
	}

	for k, v := range indexsWithWeight {
		wb.nodeIndexWithWeights[k] = Tuple{v, v}
	}

	return wb
}

func (r *WeightedBalancer) Get() int64 {
	// nodes := []int64{}
	// weights := []int64{}

	// totalEffectiveWeight := int64(0)
	// maxWeightNode := int64(0)

	// for n, tw := range r.nodeIndexWithWeights {
	// 	nodes = append(nodes, n)

	// 	totalEffectiveWeight += tw.b
	// 	tw.a += tw.b

	// }

	return 0
}
