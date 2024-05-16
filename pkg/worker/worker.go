package worker

import (
	`fmt`
	"log"
	"sync"
	`sync/atomic`
	"time"

	`github.com/spf13/afero`

	database `aurora/pkg/db`
	`aurora/pkg/model`
)

const NUM_WORKERS = 4

type Worker struct {
	db      *database.Database
	pending *int32
	refresh *time.Ticker
	reflock sync.Mutex
	stopper chan bool
}

func NewWorker(db *database.Database) *Worker {
	pending := int32(0)
	return &Worker{db: db, pending: &pending}
}

func (w *Worker) FeedsPending() int32 {
	return *w.pending
}

func (w *Worker) StartFeedCleaner() {
	//go w.db.DeleteOldItems()
	//ticker := time.NewTicker(time.Hour * 24)
	//go func() {
	//	for {
	//		<-ticker.C
	//		w.db.DeleteOldItems()
	//	}
	//}()
}

func (w *Worker) FindFavicons() {

	//fmt.Println(len(w.db.ListFeedsMissingIcons()))
	go func() {
		for _, feed := range w.db.ListFeedsMissingIcons() {
			fmt.Println(feed)
			w.FindFeedFavicon(feed)
		}
	}()
}

func (w *Worker) FindFeedFavicon(feed model.Feed) {
	icon, err := findFavicon(feed.Link, feed.FeedLink)
	fmt.Println("icon, err")
	fmt.Println(icon, err)
	appfs := afero.NewOsFs()
	appfs.MkdirAll("data/icon", 0755)
	afero.WriteFile(appfs, fmt.Sprintf("data/icon/%s.ico", feed.Id), *icon, 0644)
	//err = ioutil.WriteFile("data/"+feed.Id+".icon", *icon, 0644)
	//fmt.Println("err", err, "data/"+feed.Id+".icon")

	if err != nil {
		log.Printf("Failed to find favicon for %s (%s): %s", feed.FeedLink, feed.Link, err)
	}
	feed.Icon = fmt.Sprintf("data/icon/%s.ico", feed.Id)
	if icon != nil {
		w.db.UpdateFeedIcon(feed)
	}
}

func (w *Worker) SetRefreshRate(minute int64) {
	if w.stopper != nil {
		w.refresh.Stop()
		w.refresh = nil
		w.stopper <- true
		w.stopper = nil
	}

	if minute == 0 {
		return
	}

	w.stopper = make(chan bool)
	w.refresh = time.NewTicker(time.Minute * time.Duration(minute))

	go func(fire <-chan time.Time, stop <-chan bool, m int64) {
		log.Printf("auto-refresh %dm: starting", m)
		for {
			select {
			case <-fire:
				log.Printf("auto-refresh %dm: firing", m)
				w.RefreshFeeds()
			case <-stop:
				log.Printf("auto-refresh %dm: stopping", m)
				return
			}
		}
	}(w.refresh.C, w.stopper, minute)
}

func (w *Worker) RefreshFeeds() {
	w.reflock.Lock()
	defer w.reflock.Unlock()

	if *w.pending > 0 {
		log.Print("Refreshing already in progress")
		return
	}

	feeds := w.db.ListFeeds()
	if len(feeds) == 0 {
		log.Print("Nothing to refresh")
		return
	}

	log.Print("Refreshing feeds")
	atomic.StoreInt32(w.pending, int32(len(feeds)))
	go w.refresher(feeds)
}

func (w *Worker) refresher(feeds []model.Feed) {
	//w.db.ResetFeedErrors()

	srcqueue := make(chan model.Feed, len(feeds))
	dstqueue := make(chan []model.Item)

	for i := 0; i < NUM_WORKERS; i++ {
		go w.worker(srcqueue, dstqueue)
	}

	for _, feed := range feeds {
		srcqueue <- feed
	}
	for i := 0; i < len(feeds); i++ {
		items := <-dstqueue

		fmt.Println("======", len(items))
		if len(items) > 0 {
			w.db.CreateItems(items)
			//w.db.SetFeedSize(items[0].FeedId, len(items))
		}
		atomic.AddInt32(w.pending, -1)
		//w.db.SyncSearch()
	}
	close(srcqueue)
	close(dstqueue)

	log.Printf("Finished refreshing %d feeds", len(feeds))
}

func (w *Worker) worker(srcqueue <-chan model.Feed, dstqueue chan<- []model.Item) {
	for feed := range srcqueue {
		items, err := listItems(feed, w.db)
		if err != nil {
			//w.db.SetFeedError(feed.Id, err)
		}
		dstqueue <- items
	}
}
