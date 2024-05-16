package main

import (
	`fmt`

	`github.com/davecgh/go-spew/spew`

	database `aurora/pkg/db`
	`aurora/pkg/model`
	`aurora/pkg/worker`
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

	result, _ := worker.DiscoverFeed("https://www.36kr.com/feed")
	//spew.Dump(result)

	feed, err := db.Db.Create("feed", &model.Feed{
		Title:       result.Feed.Title,
		Description: "",
		Link:        result.Feed.SiteURL,
		FeedLink:    result.FeedLink,
	})
	//
	spew.Dump(feed)
	w := worker.NewWorker(db)
	w.FindFavicons()
	for {

	}
	////feedMap := feed.([]interface{})[0].(model.Feed)
	//
	//feedMap := feed.([]interface{})[0].(map[string]interface{})
	////把feedMap转化为model.Feed
	//ff := model.Feed{
	//	Id:       feedMap["id"].(string),
	//	Title:    feedMap["title"].(string),
	//	FeedLink: feedMap["feed_link"].(string),
	//	Link:     feedMap["link"].(string),
	//}
	//items := worker.ConvertItems(result.Feed.Items, ff)
	////fmt.Println(len(items))
	//if len(items) > 0 {
	//	ff.Size = int64(len(items))
	//	//for _, item := range items {
	//	//	db.Db.Create("feed_item", item)
	//	//}
	//	fmt.Println(ff)
	//
	//	db.Db.Update(ff.Id, ff)
	//	//s.db.SyncSearch()
	//}
	//
	//spew.Dump(ff)
	//spew.Dump(feedMap)
	//feedData := feedMap["result"].([]interface{})[0].(map[string]interface{})
	//feedMap := feed.([]interface{})[0].(model.Feed)

	//w.FindFavicons()
	//for {
	//}
	//
	//result, err := db.Db.Query("SELECT * FROM feed_item where feed_id=$feed_id ", map[string]interface{}{"feed_id": "feed:i9aa5doateotfqmwsl17"})
	//
	////feeds := []model.Feed{}
	//
	//rsMap := result.([]interface{})[0].(map[string]interface{})
	//rsData := rsMap["result"].([]interface{})
	//
	////spew.Dump(rsData)
	//for _, feed := range rsData {
	//	feedMap := feed.(map[string]interface{})
	//	fmt.Println(feedMap["title"])
	//}
	//
	//	//fmt.Println(feed, i)
	//	feedMap := feed.(map[string]interface{})
	//
	//	feeds = append(feeds, model.Feed{
	//		Id:       feedMap["id"].(string),
	//		Title:    feedMap["title"].(string),
	//		FeedLink: feedMap["feed_link"].(string),
	//		Link:     feedMap["link"].(string),
	//	})
	//
	//}
	//if err != nil {
	//	fmt.Println(err)
	//}
	//db.Set("ok", "ok")
	//fmt.Println(db.Get("ok"))
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
