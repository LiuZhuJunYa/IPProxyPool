package API

import (
	"fmt"
	"testing"
)

func TestStatistics(t *testing.T) {
	num1, num2 := Statistics()
	fmt.Println("当前有数据：", num1, "有效的：", num2)
}

func TestGetIPproxy(t *testing.T) {
	GetIPproxy(1)
}

func TestSQLinjectionIP(t *testing.T) {
	SQLinjectionIP()
}

func TestSQLinjectionASCII(t *testing.T) {
	SQLinjectionASCII()
}
