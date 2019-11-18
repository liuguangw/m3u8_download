package tools

import "fmt"

func calcPercent(value, max int) string {
	if value == 0 {
		return "0%"
	} else if value == max {
		return "100%"
	}
	return fmt.Sprintf("%.2f", float32(value)*float32(100)/float32(max)) + "%"
}
