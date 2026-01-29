package main

import (
	"fmt"

	"github.com/eventuallyconsistentwrites/high-tide-server/countmin" // Use your module path from go.mod
)

func main() {
	// Initialize
	cms := countmin.NewCountMinSketch(3, 4)

	// Display
	cms.DisplayCMS()

	// Insert Values
	cms.Update("Iphone")
	fmt.Println("Inserted Iphone")
	cms.DisplayCMS()

	cms.Update("Android")
	fmt.Println("Inserted Android")
	cms.DisplayCMS()

	cms.Update("Windows")
	fmt.Println("Inserted Windows")
	cms.DisplayCMS()

	//Query Value
	val := "Android"
	count := cms.PointQuery(val)
	fmt.Printf("Estimated count of %s = %d\n", val, count)
}
