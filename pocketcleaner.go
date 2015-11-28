// keep your pocket clean

package pocketcleaner

import (
	"encoding/json"
	"net/http"
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

// client struct to interact with the API. This mostly holds the API token
// and secret, but also provides a way to mock out the HTTP client library so
// the code is easier to test.
type PocketClient struct {
	BaseURL        string
	ConsumerSecret string
	APIToken       string
	HTTPClient     *http.Client
	KeepCount      int
}

// PocketClientWithKeys returns a PocketClient with the provided token and
// consumer secret set as well as the provided number of articles to keep.
func PocketClientWithToken(apitoken string, consumer_secret string, to_keep int) *PocketClient {
	return &PocketClient{
		ConsumerSecret: consumer_secret,
		APIToken:       apitoken,
		KeepCount:      to_keep,
		HTTPClient:     &http.Client{},
		BaseURL:        "https://getpocket.com/v3/",
	}
}

// filters out the newest `count` items from an array of PocketItems and
// returns the resulting array. This is so the returned array can be fed
// directly into ArchiveItems.
func filterOutNewestItems(list PocketItemArray, count int) PocketItemArray {
	if len(list) < count {
		return make([]PocketItem, 0)
	}
	sort.Sort(list)

	return list[0 : len(list)-count]
}

// this parses a JSON string into a PocketResponse object. It basically only
// calls json.Unmarshal() but it's stuck into a function for usability and
// testability.
func parsePocketResponse(response string) (PocketResponse, error) {
	ret := PocketResponse{}
	err := json.Unmarshal([]byte(response), &ret)
	return ret, err
}

// get all items from the configured pocket account. This is used to then
// filter out the ones to keep and archive the rest
func (c *PocketClient) getAllPocketItems() (PocketItemArray, error) {
	ret := make(PocketItemArray, 0)

	return ret, nil
}

// all items passed into this function will be archived. If one or more items
// couldn't be archived, error is != nil and the returned array contains all
// items that couldn't be archived
func (c *PocketClient) archiveItems(list PocketItemArray) (error, PocketItemArray) {
	ret := make(PocketItemArray, 0)
	return nil, ret
}

// helper method to call the pocket API via different methods
func (c *PocketClient) callPocketAPI(method string) {
}

// this is the main interface to use this from. After configuring the client
// with access token and consumer secret and the number of items to keep, just
// run this method and it will clean up your pocket account.
func (c *PocketClient) CleanUpItems() error {
	return nil
}
