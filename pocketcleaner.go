// keep your pocket clean

package pocketcleaner

import (
	"encoding/json"
	"sort"
)

// type for the actual item types coming from the pocket API. The most
// important bit here for the cleaner functionality is the TimeAdded.
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

// implement sort interface for PocketItemArray.
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

// type declaration for the list of pocket articles in the API response
type PocketItemList map[string]PocketItem

// search meta type for parsing out the pocket API response
type PocketSearchMeta struct {
	SearchType string `json:"search_type"`
}

// type to hold the response of a call to the pocket API. The only thing we
// are really interested in here is the list of articles.
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

// filters out the newest `count` items from an array of PocketItems and
// returns the resulting array. This is so the returned array can be fed
// directly into ArchiveItems.
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

// this parses a JSON string into a PocketResponse object. It basically only
// calls json.Unmarshal() but it's stuck into a function for usability and
// testability.
func ParsePocketResponse(response string) (PocketResponse, error) {
	ret := PocketResponse{}
	err := json.Unmarshal([]byte(response), &ret)
	return ret, err
}
