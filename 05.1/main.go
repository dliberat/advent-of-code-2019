package intcode

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		retval := Run(scanner.Text())
		fmt.Println("Run completed. Output: ", retval)
	}
}
