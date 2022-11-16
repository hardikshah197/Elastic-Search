package main

import (
    "bufio"
    // "bytes"
    "context"
    "encoding/json"
    "fmt"
    "os"

    "github.com/elastic/go-elasticsearch/esapi"
	"github.com/elastic/go-elasticsearch/v7"
)
var es, _ = elasticsearch.NewDefaultClient()

func Exit() {
	fmt.Println("Goodbye!")
	os.Exit(0)
}

func ReadText(reader *bufio.Scanner, prompt string) string {
	fmt.Print(prompt + ": ")
	reader.Scan()
	return reader.Text()
}

func main() {
	reader := bufio.NewScanner(os.Stdin)
	for {
        fmt.Println("0) Exit")
        fmt.Println("1) Load spacecraft")
        fmt.Println("2) Get spacecraft")
        option := ReadText(reader, "Enter option")
        if option == "0" {
            Exit()
        } else if option == "1" {
            LoadData()
        } else {
            fmt.Println("Invalid option")
        }
	}
}

func Get(reader *bufio.Scanner) {
	id := ReadText(reader, "Enter spacecraft ID")
	request := esapi.GetRequest{Index: "stsc", DocumentID: id}
	response, _ := request.Do(context.Background(), es)
	var results map[string]interface{}
	json.NewDecoder(response.Body).Decode(&results)
	Print(results["_source"].(map[string]interface{}))
}

func Print(spacecraft map[string]interface{}) {
	name := spacecraft["name"]
	status := ""
	if spacecraft["status"] != nil {

		status = "- " + spacecraft["status"].(string)
	}
	registry := ""
	if spacecraft["registry"] != nil {

		registry = "- " + spacecraft["registry"].(string)
	}
	class := ""
	if spacecraft["spacecraftClass"] != nil {

		class = "- " + spacecraft["spacecraftClass"].(map[string]interface{})["name"].(string)
	}
	fmt.Println(name, registry, class, status)
}