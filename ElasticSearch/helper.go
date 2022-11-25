package elasticsearch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/elastic/go-elasticsearch/esapi"
	"github.com/elastic/go-elasticsearch/v7"
	osrmodels "github.com/osr/pkg/models"
	"gorm.io/gorm"
)
var es, _ = elasticsearch.NewDefaultClient()

type SearchRequest struct {
	db *gorm.DB
}

func Initiator(dbClient *gorm.DB) ElasticSearch {
	return &SearchRequest{
		db: dbClient,
	}
}

func (svc SearchRequest) Get(queryString string) (map[string]interface{}, error) {
	// If application needs to fetch then, queryString should be like: "application,{application_id}""
	// and in case of blog, the queryString should be like: "blog,{blog_id}"
	i := strings.Split(queryString, ",")[0]

	var id string
	if i == "application" {
		id = strings.Split(queryString, ",")[1]
	} else {
		id = strings.Split(queryString, ",")[1]
	} 
	
	request := esapi.GetRequest{Index: i, DocumentID: id}
	response, err := request.Do(context.Background(), es)

	if err != nil {
		panic(err)
	}
	
	var results map[string]interface{}
	json.NewDecoder(response.Body).Decode(&results)
	result := results["_source"].(map[string]interface{})
	
	fmt.Printf("record found: %v\n", result)

	return result, nil
}


func (svc SearchRequest) LoadApplication(data osrmodels.Application) error {
	appid := data.ID
	jsonString, _ := json.Marshal(data)
	request := esapi.IndexRequest{Index: "application", DocumentID: appid, Body: strings.NewReader(string(jsonString))}
	_, err := request.Do(context.Background(), es)
	if err != nil {
		panic(err)
	}

	print("application loaded successfully\n")

	return nil
}

func (svc SearchRequest) LoadBlog(data osrmodels.BlogData) error {

	blogid := data.ID
	jsonString, _ := json.Marshal(data)
	request := esapi.IndexRequest{Index: "blog", DocumentID: blogid, Body: strings.NewReader(string(jsonString))}
	_, err := request.Do(context.Background(), es)
	if err != nil {
		panic(err)
	}

	print("blog loaded successfully\n")

	return nil
}

func (svc SearchRequest) LoadAllApplications() error {
	var applications []osrmodels.Application
	
	appModel := osrmodels.Application{}
	app_data, err := appModel.List(svc.db)
	
	if err != nil {
		panic(err)
	}
	
	for _, app := range app_data {
		applications = append(applications, app)
	}

	for _, data := range applications {
		appid := data.ID
		jsonString, _ := json.Marshal(data)
		request := esapi.IndexRequest{Index: "application", DocumentID: appid, Body: strings.NewReader(string(jsonString))}
		_, err := request.Do(context.Background(), es)
		if err != nil {
			panic(err)
		}
	}

	print(len(applications)," records loaded\n")
	return nil
}

func (svc SearchRequest) LoadAllBlogs() error {

	var blogs [] osrmodels.BlogData

	
	blogModel := osrmodels.BlogData{}
	blog_data, errr := blogModel.ListAll(svc.db)

	if errr != nil {
		panic(errr)
	}

	for _, blog := range blog_data {
		blogs = append(blogs, blog)
	}

	for _, data := range blogs {
		blogid := data.ID
		jsonString, _ := json.Marshal(data)
		request := esapi.IndexRequest{Index: "blog", DocumentID: blogid, Body: strings.NewReader(string(jsonString))}
		_, err := request.Do(context.Background(), es)
		if err != nil {
			panic(err)
		}
	}

	print(len(blogs), " records loaded\n")
	return nil
}

func (svc SearchRequest) Search(value, index string) ([]map[string]interface{}, error) {
	output := make([]map[string]interface{}, 0)
	fields := make([]string, 0)

	if index == "application" {
		fields = append(fields, "name^5")
		fields = append(fields, "category^4")
		fields = append(fields, "tags^3")
		fields = append(fields, "short_description^2")
		fields = append(fields, "long_description^1")
	} else {
		fields = append(fields, "title^3")
		fields = append(fields, "content^2")
		fields = append(fields, "tags^1")
	}

	var buffer bytes.Buffer
	queryMatch := map[string]interface{}{
		"query": map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query": value,
				"fields": fields,
			},
		},
	}
	json.NewEncoder(&buffer).Encode(queryMatch)
	matchResponse, err := es.Search(es.Search.WithIndex(index), es.Search.WithBody(&buffer))
	
	if err != nil {
		panic(err)
	}
	
	var result map[string]interface{}
	json.NewDecoder(matchResponse.Body).Decode(&result)
	for _, hit := range result["hits"].(map[string]interface{})["hits"].([]interface{}) {
		craft := hit.(map[string]interface{})["_source"].(map[string]interface{})
		output = append(output, craft)
	}

	if len(output) == 0 {
		queryPrefix := map[string]interface{}{
			"query": map[string]interface{}{
				"prefix": map[string]interface{}{
					"name": value,
				},
			},
		}
		json.NewEncoder(&buffer).Encode(queryPrefix)
		prefixResponse, err := es.Search(es.Search.WithIndex(index), es.Search.WithBody(&buffer))

		if err != nil {
			panic(err)
		}

		json.NewDecoder(prefixResponse.Body).Decode(&result)
		for _, hit := range result["hits"].(map[string]interface{})["hits"].([]interface{}) {
			craft := hit.(map[string]interface{})["_source"].(map[string]interface{})
			output = append(output, craft)
		}
	}

	return output, nil
}
