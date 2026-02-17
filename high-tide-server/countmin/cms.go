package countmin

import (
	"fmt"
	"math"
	"strings"
	"sync"

	"github.com/twmb/murmur3"
)

type CountMinSketch struct {
	NumberOfHashFunctions int
	certainty             float64
	Width                 uint32
	errorMargin           float64
	cmsTable              [][]int
	mu                    sync.RWMutex
}

func NewCountMinSketch(certainty float64, errorMargin float64) *CountMinSketch {
	width := math.Ceil(math.E / errorMargin)
	numberOfHashFunctions := math.Ceil(math.Log(1 / certainty))
	cms := &CountMinSketch{
		NumberOfHashFunctions: int(numberOfHashFunctions),
		certainty:             certainty,
		Width:                 uint32(width),
		errorMargin:           errorMargin,
		cmsTable:              make([][]int, int(numberOfHashFunctions)),
	}
	for i := range cms.cmsTable {
		cms.cmsTable[i] = make([]int, cms.Width)
	}
	return cms
}

// String returns a string representation of the Count-Min Sketch table,
// satisfying the fmt.Stringer interface. This is useful for debugging.
func (cms *CountMinSketch) String() string {
	cms.mu.RLock()
	defer cms.mu.RUnlock()

	var b strings.Builder
	b.WriteString("\t: ")
	for i := 0; i < int(cms.Width); i++ {
		b.WriteString(fmt.Sprintf("%d\t", i))
	}
	b.WriteRune('\n')
	for i := 0; i < cms.NumberOfHashFunctions; i++ {
		b.WriteString(fmt.Sprintf("h%d\t: ", i))
		for j := 0; j < int(cms.Width); j++ {
			b.WriteString(fmt.Sprintf("%d\t", cms.cmsTable[i][j]))
		}
		b.WriteRune('\n')
	}
	return b.String()
}

func (cms *CountMinSketch) Update(value string) {
	cms.mu.Lock()
	defer cms.mu.Unlock()
	var data []byte = []byte(value)
	for i := 0; i < cms.NumberOfHashFunctions; i++ {
		var seed uint32 = uint32(i)
		hashValue := murmur3.SeedSum32(seed, data)
		cms.cmsTable[i][hashValue%cms.Width]++
	}
}

func (cms *CountMinSketch) PointQuery(value string) int {
	cms.mu.RLock()
	defer cms.mu.RUnlock()
	var data []byte = []byte(value)
	minFreq := math.MaxInt
	for i := 0; i < cms.NumberOfHashFunctions; i++ {
		var seed uint32 = uint32(i)
		hashValue := murmur3.SeedSum32(seed, data)
		freq := cms.cmsTable[i][hashValue%cms.Width]
		if freq < minFreq {
			minFreq = freq
		}
	}
	return minFreq
}

func (cms *CountMinSketch) Reset() {
	cms.mu.Lock()
	defer cms.mu.Unlock()
	for i := range cms.cmsTable {
		for j := range cms.cmsTable[i] {
			cms.cmsTable[i][j] = 0
		}
	}
}
