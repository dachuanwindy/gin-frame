package es

import (
	"encoding/json"
	"fmt"
	"git.ymt360.com/go/gocommons/logging"
	"github.com/olivere/elastic"
	"reflect"
	"testing"
	modules "zhuayu-search/module"
	//"zhuayu-search/util"
)

const (
	Index = "zhuayu_user_alias"
	Typ   = "user_info"
)

//func init() {
//	util.InitConfig("", "default")
//	util.InitLog()
//}
func TestInitES(t *testing.T) {
	InitES("zhuayu-feature-es", "http://dev-uc-host003.ymt.io:9200")
}
func TestCreateIndex(t *testing.T) {
	es := InitES("zhuayu_feature_es", "http://dev-uc-host003.ymt.io:9200")
	mapping := `
	{
	"settings":{
		"number_of_shards":1,
		"number_of_replicas":0
	},
	"mappings":{
			"user_info":{
				"properties":{
					"customer_id":{
						"type":"keyword"
					},
					"nickname":{
						"type":"keyword"
					},
					"user_type":{
						"type":"keyword"
					},
					"location":{
						"type":"keyword"
					},
					"main_product":{
						"type":"text"
					},
					"introduce":{
						"type":"text"
					}
				}
			}
		}
	}
	`
	t.Log(CreateIndex(es.Client, "zhuayu_user_test", mapping, logging.NewLogHeader()))
}

func TestDelIndex(t *testing.T) {
	es := InitES("zhuayu_feature_es", "http://dev-uc-host003.ymt.io:9200")
	t.Log(DelIndex(es.Client, Index, logging.NewLogHeader()))
}

func TestPut(t *testing.T) {
	es := InitES("zhuayu-feature-es", "http://dev-uc-host003.ymt.io:9200")
	user := modules.ZhuayuUser{
		CustomerId:     12345,
		Nickname:       "我是中文",
		LocationDetail: "北京",
		UserType:       "卖家",
		MainProduct:    "苹果",
		InAWord:        "换个签名",
	}
	str, _ := json.Marshal(user)
	t.Log(Put(es.Client, Index, Typ, string(user.CustomerId), string(str), logging.NewLogHeader()))
}

func TestUpdate(t *testing.T) {
	es := InitES("zhuayu-feature-es", "http://dev-uc-host003.ymt.io:9200")
	updateMap := make(map[string]interface{})
	updateMap["nickname"] = "张大三"
	t.Log(Update(es.Client, Index, Typ, "12345", updateMap, logging.NewLogHeader()))
}

func TestTermQuery(t *testing.T) {
	es := InitES("zhuayu_feature_es", "http://dev-uc-host003.ymt.io:9200")
	searchResult := TermQuery(es.Client, Index, Typ, elastic.NewTermQuery("customer_id", 8086763), 0, 10, logging.NewLogHeader())
	var typ modules.ZhuayuUser
	for _, item := range searchResult.Each(reflect.TypeOf(typ)) {
		t := item.(modules.ZhuayuUser)
		fmt.Printf("ZhuayuUser by %d: %s\n", t.CustomerId, t.Nickname)
	}
	// TotalHits is another convenience function that works even when something goes wrong.
	fmt.Printf("Found a total of %d ZhuayuUser\n", searchResult.TotalHits())

	// Here's how you iterate through results with full control over each step.
	if searchResult.TotalHits() > 0 {
		fmt.Printf("Found a total of %d ZhuayuUser\n", searchResult.TotalHits())

		// Iterate through results
		for _, hit := range searchResult.Hits.Hits {

			// Deserialize hit.Source into a ZhuayuUser (could also be just a map[string]interface{}).
			var t modules.ZhuayuUser
			err := json.Unmarshal(*hit.Source, &t)
			if err != nil {
				// Deserialization failed
			}

			fmt.Printf("ZhuayuUser by %d: %s\n", t.CustomerId, t.Nickname)
		}
	} else {
		// No hits
		fmt.Print("Found no ZhuayuUser\n")
	}
}

func TestMultiMatchQueryBestFields(t *testing.T) {
	es := InitES("zhuayu_feature_es", "http://dev-uc-host003.ymt.io:9200")

	searchResult := MultiMatchQueryBestFields(es.Client, Index, Typ, "柚子", logging.NewLogHeader(), "nickname", "location_detail", "main_product", "in_a_word")
	if searchResult == nil {
		t.Log("search nothing")
		return
	}
	for _, hit := range searchResult.Hits.Hits {
		user := modules.ZhuayuUser{}
		_ = json.Unmarshal(*hit.Source, &user)
		fmt.Println(user)
	}

}
func TestQueryString(t *testing.T) {
	es := InitES("zhuayu_feature_es", "http://dev-uc-host003.ymt.io:9200")
	query := "nickname:胡琰川 AND location_detail:江苏 AND main_product:种子"
	searchResult := QueryString(es.Client, Index, Typ, query, 200, logging.NewLogHeader())
	if searchResult == nil {
		t.Log("search nothing")
		return
	}
	for _, hit := range searchResult.Hits.Hits {
		user := modules.ZhuayuUser{}
		_ = json.Unmarshal(*hit.Source, &user)
		fmt.Println(user)
	}
}

func TestQueryStringRandomSearch(t *testing.T) {
	query := "nickname:?* AND location_detail:?* AND main_product:?* AND user_type:?*"
	es := InitES("zhuayu_feature_es", "http://dev-uc-host003.ymt.io:9200")
	searchResult := QueryStringRandomSearch(es.Client, Index, Typ, query, 200, logging.NewLogHeader())
	if searchResult == nil {
		t.Log("search nothing")
		return
	}
	users := make([]modules.ZhuayuUser, 0)
	for _, hit := range searchResult.Hits.Hits {
		user := modules.ZhuayuUser{}
		_ = json.Unmarshal(*hit.Source, &user)
		fmt.Println(user)
		users = append(users, user)
	}
	t.Log(users)
}

func TestTermQuery2(t *testing.T) {
	es := InitES("zhuayu_feature_es", "http://dev-uc-host003.ymt.io:9200")

	searchResult := MultiMatchQueryBestFields(es.Client, Index, Typ, "柚子", logging.NewLogHeader(), "nickname", "location_detail", "main_product", "in_a_word")
	if searchResult == nil {
		t.Log("search nothing")
		return
	}
	user := modules.ZhuayuUser{}
	json.Unmarshal(*searchResult.Hits.Hits[1].Source, &user)
	fmt.Printf("%+v", user)
}

func TestRangeQueryLoginDate(t *testing.T) {
	es := InitES("zhuayu_feature_es", "http://dev-uc-host003.ymt.io:9200")
	searchResult := RangeQueryLoginDate(es.Client, Index, Typ, logging.NewLogHeader())
	//t.Log(searchResult)
	user := modules.ZhuayuUser{}
	json.Unmarshal(*searchResult.Hits.Hits[0].Source, &user)
	fmt.Printf("%+v", user)
}
