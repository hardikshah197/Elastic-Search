package elasticsearch

import (
    "bufio"
    "context"
    "encoding/json"

    "github.com/elastic/go-elasticsearch/esapi"
)

func Get(reader *bufio.Scanner) {
	id := ReadText(reader, "Enter spacecraft ID")
	request := esapi.GetRequest{Index: "stsc", DocumentID: id}
	response, _ := request.Do(context.Background(), es)
	var results map[string]interface{}
	json.NewDecoder(response.Body).Decode(&results)
	Print(results["_source"].(map[string]interface{}))
}