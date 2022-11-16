package elasticsearch

import (
    "bufio"
    "fmt"
    "os"

	"github.com/elastic/go-elasticsearch/v7"
)
var es, _ = elasticsearch.NewDefaultClient()

func Initiator() {
	reader := bufio.NewScanner(os.Stdin)
	for {
        fmt.Println("0) Exit")
		fmt.Println("1) Load spacecraft")
		fmt.Println("2) Get spacecraft")
		fmt.Println("3) Search spacecraft by key and value")
		fmt.Println("4) Search spacecraft by key and prefix")
		option := ReadText(reader, "Enter option")
		if option == "0" {
			Exit()
		} else if option == "1" {
			Load()
		} else if option == "2" {
			Get(reader)
		} else if option == "3" {
			Search(reader, "match")
		} else if option == "4" {
			Search(reader, "prefix")
		} else {
			fmt.Println("Invalid option")
		}
	}
}