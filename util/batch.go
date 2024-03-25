package util

func BatchItems(itemIds []int, batchSize int) [][]int {
	var batches [][]int
	for batchSize < len(itemIds) {
		itemIds, batches = itemIds[batchSize:], append(batches, itemIds[0:batchSize:batchSize])
	}
	batches = append(batches, itemIds)
	return batches
}
