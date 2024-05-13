package main

import (
	`fmt`

	database `aurora/pkg/db`
)

const (
	wsString  = "ws://localhost:9876/rpc"
	tableName = "user"
	namespace = "test"
	dbName    = "test"
	port      = 9876
	username  = "root"
	password  = "root"
)

func main() {
	db, err := database.NewDatabase(wsString, username, password, namespace, dbName)
	if err != nil {
		fmt.Println(err)
	}
	db.Set("ok", "ok")
	fmt.Println(db.Get("ok"))
	//feed, _ := db.Db.Create("feed", &rss.Feed{
	//	Title:       "result.Feed.Title",
	//	Description: "",
	//	Link:        "result.Feed.SiteURL",
	//})
	//spew.Dump(feed)
	//feedMap := feed.([]interface{})[0].(map[string]interface{})
	////dataData := dataMap["result"].([]interface{})
	////feed
	////feedMap := feed.(map[interface{}])
	////feedData := feedMap[0].(map[string]interface{})
	//
	//fmt.Println(feedMap)
	//fmt.Println(feedMap["id"])
	////fmt.Println(feedData)
	//
	//println("Hello, world!")
}
