package utils

import (
	"fmt"
	"strings"
)

func PrintSection(title string) {
	fmt.Printf("\n%s\n", Cyan(title))
}

func PrintKV(key string, value interface{}) {
	if value == nil {
		return
	}
	valStr := ""
	switch v := value.(type) {
	case string:
		if v == "" {
			return
		}
		valStr = v
	case []int:
		if len(v) == 0 {
			return
		}
		var s []string
		for _, i := range v {
			s = append(s, fmt.Sprintf("%d", i))
		}
		valStr = strings.Join(s, ", ")
	case []string:
		if len(v) == 0 {
			return
		}
		valStr = strings.Join(v, ", ")
	case bool:
		valStr = fmt.Sprintf("%t", v)
	default:
		s := fmt.Sprintf("%v", v)
		if s == "" {
			return
		}
		valStr = s
	}

	fmt.Printf("  %-25s %s\n", Yellow(key+":"), valStr)
}
