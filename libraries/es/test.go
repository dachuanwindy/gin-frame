package es

import (
	"fmt"
	"github.com/olivere/elastic"
	"gin-frame/libraries/es"
	"gin-frame/libraries/log"
)

type Res struct {
	name 	string
	age 	int
	gender 	string
}

func TermMap(esConn *es.Elastic, index, typ string, logFormat *log.LogFormat){
	searchResult := esConn.TermQueryMap(index, typ, elastic.NewTermQuery("name", "why"), 0,10, logFormat)
	fmt.Println(searchResult)
}

func StringMap(esConn *es.Elastic, index, typ string, logFormat *log.LogFormat){
	query := "name:why OR age:19"
	searchResult := esConn.QueryStringMap(index, typ, query, 0,10, logFormat)
	fmt.Println(searchResult)
}

func MultiMatchQueryBestFieldsMap(esConn *es.Elastic, index, typ string, logFormat *log.LogFormat){
	query := "why"
	searchResult := esConn.MultiMatchQueryBestFieldsMap(index, typ, query, 0,10, logFormat, "name", "desc")
	fmt.Println(searchResult)
}

func JsonMap(esConn *es.Elastic, index,typ string, logFormat *log.LogFormat){
	terms := make(map[string]interface{})
	terms["type"] = []interface{}{
		"hangqing_analyze",
		"hangqing",
		"hangqing_desc",
		"scdm",
	}

	var mustNot []map[string]interface{}
	mustNot = append(mustNot, map[string]interface{}{
		"term": map[string]interface{}{
			"img_count": 0,
		},
	})

	var filter []map[string]interface{}
	filter = append(filter, map[string]interface{}{
		"range": map[string]interface{}{
			"video_count": map[string]interface{}{
				"lte":1,
			},
		},
	})

	var sort []map[string]interface{}
	sort = append(sort, map[string]interface{}{
		"_score": map[string]interface{}{
			"order": "desc",
		},
	})

	res := esConn.JsonMap(index, typ, "小龙虾批发", []string{"content", "tags"},0, 3, terms, mustNot, filter, sort, logFormat )
	fmt.Println(res)
}

func main(){
	index := "why_index"
	typ   := "why_type"
	logFormat := new(log.LogFormat)
	esConn := es.InitES("default")

	TermMap(esConn, index, typ, logFormat)
	StringMap(esConn, index, typ, logFormat)
	MultiMatchQueryBestFieldsMap(esConn, index, typ, logFormat)

	esConn = es.InitES("ymt")
	index = "dynamic"
	typ = "detail"
	JsonMap(esConn, index, typ, logFormat)
}