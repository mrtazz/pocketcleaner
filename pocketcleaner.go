// keep your pocket clean

package pocketcleaner

import (
	"encoding/json"
	"sort"
)

type PocketItem struct {
	GivenTitle    string `json:"given_title"`
	IsArticle     string `json:"is_article"`
	ResolvedID    string `json:"resolved_id"`
	Status        string `json:"status"`
	SortId        uint   `json:"sort_id"`
	HasImage      string `json:"has_image"`
	Excerpt       string `json:"excerpt"`
	TimeFavorited string `json:"time_favorited"`
	WordCount     string `json:"word_count"`
	TimeRead      string `json:"time_read"`
	ResolvedTitle string `json:"resolved_title"`
	TimeUpdated   string `json:"time_updated"`
	Favorite      string `json:"favorite"`
	HasVideo      string `json:"has_video"`
	TimeAdded     uint64 `json:"time_added,string"`
	ResolvedUrl   string `json:"resolved_url"`
	GivenUrl      string `json:"given_url"`
	IsIndex       string `json:"is_index"`
	ItemID        string `json:"item_id"`
}

// // implement sort interface for PocketItemArray
type PocketItemArray []PocketItem

func (list PocketItemArray) Len() int {
	return len(list)
}

func (list PocketItemArray) Less(i, j int) bool {
	return list[i].TimeAdded < list[j].TimeAdded
}

func (list PocketItemArray) Swap(i, j int) {
	list[i], list[j] = list[j], list[i]
}

type PocketItemList map[string]PocketItem

type PocketSearchMeta struct {
	SearchType string `json:"search_type"`
}

type PocketResponse struct {
	SearchMeta PocketSearchMeta `json:"search_meta"`
	Status     uint             `json:"status"`
	Complete   uint             `json:"complete"`
	List       PocketItemList   `json:"list"`
	Error      bool             `json:"error"`
	Since      uint             `json:"since"`
}

func GetAllPocketItems() {
}

func FilterOutNewestItems(list PocketItemArray, count int) PocketItemArray {
	if len(list) < count {
		return make([]PocketItem, 0)
	}
	sort.Sort(list)

	return list[0 : len(list)-count]
}

func ArchiveItems() {
}

func CallPocketAPI(method string, consumer_key string, access_token string) {
}

func ParsePocketResponse(response string) (PocketResponse, error) {
	ret := PocketResponse{}
	err := json.Unmarshal([]byte(response), &ret)
	return ret, err
}
