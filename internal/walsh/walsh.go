package walsh

import some_stuff "CDMA-telecom-lab1/internal/some-stuff"

var H1 = [][]int64{{1}}

var H2 = [][]int64{{1, 1}, {1, -1}}

var H4 = [][]int64{{1, 1, 1, 1}, {1, -1, 1, -1}, {1, 1, -1, -1}, {1, -1, -1, 1}}

func NewWalshTable(n int64) [][]int64 {
	if n == 1 {
		return [][]int64{{1}}
	}

	walshSize := some_stuff.NextPowOf2(n)

	walshTable := make([][]int64, walshSize)
	for i := range walshTable {
		walshTable[i] = make([]int64, walshSize)
	}

	generateWalshTable(walshTable, walshSize, 0, walshSize-1, 0, walshSize-1, false)

	return walshTable[0:n]
}

func generateWalshTable(walshTable [][]int64, length, i1, i2, j1, j2 int64, isComplement bool) {
	if length == 2 {
		if !isComplement {
			walshTable[i1][j1] = 1
			walshTable[i1][j2] = 1
			walshTable[i2][j1] = 1
			walshTable[i2][j2] = -1
		} else {
			walshTable[i1][j1] = -1
			walshTable[i1][j2] = -1
			walshTable[i2][j1] = -1
			walshTable[i2][j2] = 1
		}
		return
	}

	midi := (i1 + i2) / 2
	midj := (j1 + j2) / 2

	generateWalshTable(walshTable, length/2, i1, midi, j1, midj, isComplement)
	generateWalshTable(walshTable, length/2, i1, midi, midj+1, j2, isComplement)
	generateWalshTable(walshTable, length/2, midi+1, i2, j1, midj, isComplement)
	generateWalshTable(walshTable, length/2, midi+1, i2, midj+1, j2, !isComplement)
}
