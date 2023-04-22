package Check

import (
	"fmt"
	"testing"
)

func TestTestIP(t *testing.T) {
	falg := TestIP("http://123.163.52.24:9091")
	fmt.Println(falg)
}

func TestCheckPool(t *testing.T) {
	CheckPool()
}
