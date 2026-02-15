package countmin

import (
	"fmt"
	"math"
	"strings"
	"sync"
	"github.com/twmb/murmur3"
)

type CountMinSketch struct {
	numberOfHashFunctions int
	width                 uint32
	cmsTable              [][]int
	mu                    sync.RWMutex
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

// String returns a string representation of the Count-Min Sketch table,
// satisfying the fmt.Stringer interface. This is useful for debugging.
func (cms *CountMinSketch) String() string {
	cms.mu.RLock()
	defer cms.mu.RUnlock()

	var b strings.Builder

	b.WriteString("\t: ")
	for i := 0; i < int(cms.width); i++ {
		b.WriteString(fmt.Sprintf("%d\t", i))
	}
	b.WriteRune('\n')
	for i := 0; i < cms.numberOfHashFunctions; i++ {
		b.WriteString(fmt.Sprintf("h%d\t: ", i))
		for j := 0; j < int(cms.width); j++ {
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
	for i := 0; i < cms.numberOfHashFunctions; i++ {
		var seed uint32 = uint32(i)
		hashValue := murmur3.SeedSum32(seed, data)
		cms.cmsTable[i][hashValue%cms.width]++
	}
}

func (cms *CountMinSketch) PointQuery(value string) int {
	cms.mu.RLock()
	defer cms.mu.RUnlock()
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

func (cms *CountMinSketch) Reset() {
	cms.mu.Lock()
	defer cms.mu.Unlock()
	for i := range cms.cmsTable {
		for j := range cms.cmsTable[i] {
			cms.cmsTable[i][j] = 0
		}
	}
}
