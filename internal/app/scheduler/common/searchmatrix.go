package common

func SearchMatrix(matrix [][]string, status string) (int, int) {
	for row := 0; row < len(matrix); row++ {
		nums := matrix[row]
		for col, v := range nums {
			if v == status {
				return row, col
			}
		}
	}
	return -1, -1
}
