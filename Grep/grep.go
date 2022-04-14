package main
import(
	"fmt"
	"os"
	"bufio"
	"log"
	"regexp"
	"strings"
	"github.com/fatih/color"
)

const usage = `
USAGE: go run grep.go [pattern] [file]
`

func readFileLinebyLine(filePath string,callback func([]byte))([]string,error){
	f,err := os.Open(filePath)

	if err!=nil {
		log.Fatal(err)
		return nil,err
	}

	defer f.Close()

	scanner := bufio.NewScanner(f)

	for scanner.Scan(){
		callback(scanner.Bytes())
	}
	
	return nil,nil
}

func getIntervals(positions [][]int) []int{
	var intervals []int
	for _,position := range positions{
		initialInterval := position[0]
		lastInterval := position[1]

		for initialInterval < lastInterval {
			intervals = append(intervals, initialInterval)
			initialInterval += 1
		}
	}
	return intervals
}

func intervalContainsPosition(interval []int, position int) bool {
	for _, elem := range interval {
		if position == elem {
			return true
		}
	}
	return false
}

func applyColor(line []byte, intervals []int) string {
	var stringifiedLine []string

	for charPosition, char := range string(line) {
		if intervalContainsPosition(intervals, charPosition) {
			stringifiedLine = append(stringifiedLine, color.RedString(string(char)))
			continue
		}
		stringifiedLine = append(stringifiedLine, string(char))
	}

	return strings.Join(stringifiedLine, "")
}

func main(){
	args := os.Args[1:]
	if len(args) < 2{
		fmt.Println("Missing args,both pattern string and target file are missing")
		fmt.Print(usage)
		os.Exit(0)
	}

	if len(args) >2 {
		fmt.Println("Only two params required")
		fmt.Print(usage)
		os.Exit(0)
	}

	pattern := args[0]
	filePath := args[1]

	if _,err := os.Stat(filePath); err!=nil {
		fmt.Printf("%s file does not exist",filePath)
		os.Exit(0)
	}

	re := regexp.MustCompile(pattern)

	_,readLineErr := readFileLinebyLine(filePath,func(line []byte){
		positions := re.FindAllIndex(line,-1)
		occurences := len(positions)
		if occurences>0 {
			intervals := getIntervals(positions)
			highlightedLine := applyColor(line, intervals)
			fmt.Println(highlightedLine)

		}
	})

	if readLineErr!=nil {
		log.Fatal(readLineErr)
		fmt.Printf("Error in reading File\n")
		os.Exit(0)
	}

}