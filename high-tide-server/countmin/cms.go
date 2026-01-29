package countmin

import (
	"fmt"
	"math"

	"github.com/twmb/murmur3"
)

type CountMinSketch struct {
	numberOfHashFunctions int
	width                 uint32
	cmsTable              [][]int
}

func NewCountMinSketch(numberOfHashFunctions int, width uint32) *CountMinSketch {
	cms := &CountMinSketch{
		numberOfHashFunctions: numberOfHashFunctions,
		width:                 width,
		cmsTable:              make([][]int, numberOfHashFunctions),
	}
	for i := range cms.cmsTable {
		cms.cmsTable[i] = make([]int, width)
	}
	return cms
}

func (cms *CountMinSketch) DisplayCMS() {
	fmt.Printf("\t: ")
	for i := 0; i < int(cms.width); i++ {
		fmt.Printf("%d\t", i)
	}
	fmt.Println()
	for i := 0; i < cms.numberOfHashFunctions; i++ {
		fmt.Printf("h%d\t: ", i)
		for j := 0; j < int(cms.width); j++ {
			fmt.Printf("%d\t", cms.cmsTable[i][j])
		}
		fmt.Println()
	}
}

func (cms *CountMinSketch) Update(value string) {
	var data []byte = []byte(value)
	for i := 0; i < cms.numberOfHashFunctions; i++ {
		var seed uint32 = uint32(i)
		hashValue := murmur3.SeedSum32(seed, data)
		cms.cmsTable[i][hashValue%cms.width]++
	}
}

func (cms *CountMinSketch) PointQuery(value string) int {
	var data []byte = []byte(value)
	minFreq := math.MaxInt
	for i := 0; i < cms.numberOfHashFunctions; i++ {
		var seed uint32 = uint32(i)
		hashValue := murmur3.SeedSum32(seed, data)
		freq := cms.cmsTable[i][hashValue%cms.width]
		if freq < minFreq {
			minFreq = freq
		}
	}
	return minFreq
}
