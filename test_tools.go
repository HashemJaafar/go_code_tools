
//this file created automatically
package main
import (
	"fmt"
	"log"
	"os"
	"reflect"
	"text/tabwriter"
)
func TEST[t any](should_equal bool, actual, expected t) {
	if !reflect.DeepEqual(actual, expected) == should_equal {
		p := tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0)
		fmt.Fprintln(p, "\033[35m", "should_equal\t:", should_equal) //purple
		fmt.Fprintln(p, "\033[34m", "actual\t:", actual)             //blue
		fmt.Fprintln(p, "\033[33m", "expected\t:", expected)         //yellow
		p.Flush()
		fmt.Println("\033[31m") //red
		log.Panic()
	}
}