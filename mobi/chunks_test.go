package mobi

import (
	"fmt"
	"os"
	"testing"
)

func TestChunks(t *testing.T) {

	d, _ := os.ReadFile("1.html")

	list := Split(string(d))

	for i := range list {
		//		fmt.Println(len([]byte(list[i])))

		fmt.Println(list[i])
	}
}
