// keep your pocket clean

package pocketcleaner

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"time"
)

// type for the actual item types coming from the pocket
// API. The most important bit here for the cleaner functionality is the
// TimeAdded.
type pocketItem struct {
	GivenTitle    string `json:"given_title"`
	IsArticle     string `json:"is_article"`
	ResolvedID    string `json:"resolved_id"`
	Status        string `json:"status"`
	SortID        uint   `json:"sort_id"`
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
	ResolvedURL   string `json:"resolved_url"`
	GivenURL      string `json:"given_url"`
	IsIndex       string `json:"is_index"`
	ItemID        string `json:"item_id"`
}

// PocketItemArray implement sort interface for PocketItemArray.
type pocketItemArray []pocketItem

func (list pocketItemArray) Len() int {
	return len(list)
}

func (list pocketItemArray) Less(i, j int) bool {
	return list[i].TimeAdded < list[j].TimeAdded
}

func (list pocketItemArray) Swap(i, j int) {
	list[i], list[j] = list[j], list[i]
}

// type declaration for the list of pocket articles in the API response
type pocketItemList map[string]pocketItem

// search meta type for parsing out the pocket API response
type pocketSearchMeta struct {
	SearchType string `json:"search_type"`
}

// type to hold the response of a call to the pocket API. The only thing we
// are really interested in here is the list of articles.
type pocketResponse struct {
	SearchMeta pocketSearchMeta `json:"search_meta"`
	Status     uint             `json:"status"`
	Complete   uint             `json:"complete"`
	List       pocketItemList   `json:"list"`
	Error      bool             `json:"error"`
	Since      uint             `json:"since"`
}

type pocketArchiveItem struct {
	Action string `json:"action"`
	ID     string `json:"item_id"`
	Time   string `json:"time"`
}
type pocketArchiveItemArray []pocketArchiveItem

// PocketClient struct to interact with the API. This mostly holds the API token
// and secret, but also provides a way to mock out the HTTP client library so
// the code is easier to test.
type PocketClient struct {
	BaseURL        string
	ConsumerSecret string
	APIToken       string
	HTTPClient     *http.Client
	KeepCount      int
}

// PocketClientWithToken returns a PocketClient with the provided token and
// consumer secret set as well as the provided number of articles to keep.
func PocketClientWithToken(apiToken string, consumerSecret string, toKeep int) *PocketClient {
	return &PocketClient{
		ConsumerSecret: consumerSecret,
		APIToken:       apiToken,
		KeepCount:      toKeep,
		HTTPClient:     &http.Client{},
		BaseURL:        "https://getpocket.com/v3/",
	}
}

// filters out the newest `count` items from an array of PocketItems and
// returns the resulting array. This is so the returned array can be fed
// directly into ArchiveItems.
func filterOutNewestItems(list pocketItemArray, count int) pocketItemArray {
	if len(list) < count {
		return make([]pocketItem, 0)
	}
	sort.Sort(list)

	return list[0 : len(list)-count]
}

// this parses a JSON string into a PocketResponse object. It basically only
// calls json.Unmarshal() but it's stuck into a function for usability and
// testability.
func parsePocketResponse(response string) (ret pocketResponse, err error) {
	err = json.Unmarshal([]byte(response), &ret)
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

// CleanUpItems is the main method to use this module from. After configuring
// the client with access token and consumer secret and the number of items to
// keep, just run this method and it will clean up your pocket account.
func (c *PocketClient) CleanUpItems() (err error) {
	err = errors.New("method not implemented")
	return err
}
