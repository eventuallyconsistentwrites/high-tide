package main

import (
	"fmt"

	"github.com/eventuallyconsistentwrites/high-tide-server/countmin" // Use your module path from go.mod
)

func main() {
	// Initialize
	cms := countmin.NewCountMinSketch(3, 4)

	// Display
	fmt.Println(cms)

	// Insert Values
	cms.Update("Iphone")
	fmt.Println("Inserted Iphone")
	fmt.Println(cms)

	cms.Update("Android")
	fmt.Println("Inserted Android")
	fmt.Println(cms)

	cms.Update("Windows")
	fmt.Println("Inserted Windows")
	fmt.Println(cms)

	//Query Value
	val := "Android"
	count := cms.PointQuery(val)
	fmt.Printf("Estimated count of %s = %d\n", val, count)
}
