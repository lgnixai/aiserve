package server

import (
	`embed`
	`io/fs`
	`log`
	`net/http`
	`strconv`
	`sync`

	`github.com/acheong08/endless`

	database `aurora/pkg/db`
	`aurora/pkg/worker`
)

//go:embed web/*
var staticFiles embed.FS

type Server struct {
	Host        string
	Addr        string
	db          *database.Database
	worker      *worker.Worker
	cache       map[string]interface{}
	cache_mutex *sync.Mutex

	BasePath string

	// auth
	Username string
	Password string
	// https
	CertFile string
	KeyFile  string
}

func NewServer(db *database.Database, addr string) *Server {
	return &Server{
		db:          db,
		Addr:        addr,
		worker:      worker.NewWorker(db),
		cache:       make(map[string]interface{}),
		cache_mutex: &sync.Mutex{},
	}
}

func (h *Server) GetAddr() string {
	proto := "http"
	if h.CertFile != "" && h.KeyFile != "" {
		proto = "https"
	}
	return proto + "://" + h.Addr + h.BasePath
}

func (s *Server) Start() {

	key, err := s.db.Get("refresh_rate")
	var refreshRate int64 = 10
	//refreshRate, _ := strconv.ParseInt("10", 10, 64)
	//key 转化成int64
	if err != nil {
		s.db.Set("refresh_rate", string(refreshRate))
	} else {
		//key string 转化成int64
		refreshRate, err = strconv.ParseInt(key, 10, 64)
		if err != nil {
			refreshRate = 10
		}

	}
	s.worker.FindFavicons()
	s.worker.StartFeedCleaner()
	s.worker.SetRefreshRate(refreshRate)
	if refreshRate > 0 {
		s.worker.RefreshFeeds()
	}

	if err != nil {
		refreshRate = 60
	}
	s.worker.FindFavicons()
	s.worker.StartFeedCleaner()
	s.worker.SetRefreshRate(refreshRate)
	if refreshRate > 0 {
		s.worker.RefreshFeeds()
	}
	router := s.RegisterRouter()
	subFS, err := fs.Sub(staticFiles, "web")
	if err != nil {
		log.Fatal(err)
	}
	router.StaticFS("/web", http.FS(subFS))

	if s.CertFile != "" && s.KeyFile != "" {
		_ = endless.ListenAndServeTLS(s.Host+":"+s.Addr, s.CertFile, s.KeyFile, router)
	} else {
		_ = endless.ListenAndServe(s.Host+":"+s.Addr, router)
	}
}
