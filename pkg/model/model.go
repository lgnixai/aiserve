package model

import (
	`encoding/json`
	"time"
)

type Feed struct {
	Id          string  `json:"id,omitempty"`
	FolderId    *int64  `json:"folder_id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Link        string  `json:"link"`
	FeedLink    string  `json:"feed_link"`
	Icon        *[]byte `json:"icon,omitempty"`
	HasIcon     bool    `json:"has_icon"`
	Size        int64   `json:"size"`
}

type ItemStatus int

const (
	UNREAD  ItemStatus = 0
	READ    ItemStatus = 1
	STARRED ItemStatus = 2
)

var StatusRepresentations = map[ItemStatus]string{
	UNREAD:  "unread",
	READ:    "read",
	STARRED: "starred",
}

var StatusValues = map[string]ItemStatus{
	"unread":  UNREAD,
	"read":    READ,
	"starred": STARRED,
}

func (s ItemStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(StatusRepresentations[s])
}

func (s *ItemStatus) UnmarshalJSON(b []byte) error {
	var str string
	if err := json.Unmarshal(b, &str); err != nil {
		return err
	}
	*s = StatusValues[str]
	return nil
}

type Item struct {
	Id string `json:"id,omitempty"`

	GUID     string     `json:"guid"`
	FeedId   string     `json:"feed_id"`
	Title    string     `json:"title"`
	Link     string     `json:"link"`
	Content  string     `json:"content,omitempty"`
	Date     time.Time  `json:"date"`
	Status   ItemStatus `json:"status"`
	ImageURL *string    `json:"image"`
	AudioURL *string    `json:"podcast_url"`
}
