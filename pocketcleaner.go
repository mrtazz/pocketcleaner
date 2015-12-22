// keep your pocket clean

package pocketcleaner

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
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

// pocketArchiveItem is the type for an item to be archived. The pocket API
// expects the JSON version of this
type pocketArchiveItem struct {
	Action string `json:"action"`
	ID     string `json:"item_id"`
	Time   int64  `json:"time,string"`
}
type pocketArchiveItemArray []pocketArchiveItem

type pocketArchiveResponse struct {
	ActionResults []bool `json:"action_results"`
	Status        int    `json:"status"`
}

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

var (
	// Debug can be set to true to get debug logging from pocketcleaner
	Debug = false
)

const (
	apiBatchSize int = 50
)

// debugPrint is a helper function to make it easy to print debugging
// information if the module is used in debug mode
func debugPrint(message string) {
	if Debug == true {
		log.Println(message)
	}
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
// directly into archiveItems.
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
	if err != nil {
		debugPrint(fmt.Sprintf("Error in parsePocketResponse: %s", err.Error()))
	}
	return ret, err
}

// get all items from the configured pocket account. This is used to then
// filter out the ones to keep and archive the rest
func (c *PocketClient) getAllPocketItems() (ret pocketItemArray, err error) {
	body, apiErr := c.callPocketAPI("get", nil)
	if apiErr != nil {
		return ret, apiErr
	}

	pocketResponse, parseErr := parsePocketResponse(body)
	if parseErr != nil {
		return ret, parseErr
	}

	ret = make(pocketItemArray, 0, len(pocketResponse.List))

	for _, v := range pocketResponse.List {
		ret = append(ret, v)
	}

	return ret, err
}

// all items passed into this function will be archived. If one or more items
// couldn't be archived, error is != nil and the returned array contains all
// items that couldn't be archived
func (c *PocketClient) archiveItems(list pocketItemArray) (err error, ret pocketItemArray) {
	if len(list) == 0 {
		return err, ret
	}
	var currentTime = time.Now().Unix()
	archiveItems := make(pocketArchiveItemArray, 0, apiBatchSize)
	for k, v := range list {
		archiveItems = append(archiveItems, pocketArchiveItem{Action: "archive",
			ID: v.ItemID, Time: currentTime})

		// if we have reached the batchSize, let's send the current set of items
		// off to the API and start a new batch
		if k%apiBatchSize == 0 {

			err, _ = c.archiveBatch(archiveItems)

			// reset batch container
			archiveItems = make(pocketArchiveItemArray, 0, apiBatchSize)
		}

	}
	err, _ = c.archiveBatch(archiveItems)

	return err, ret
}

// archiveBatch has almost the same signature as archiveItems but is intended
// to be used to archive smaller batches of items. This is to encapsulate the
// actual API handling logic in here and have archiveItems be a bit more about
// the app logic.
func (c *PocketClient) archiveBatch(list pocketArchiveItemArray) (err error, ret pocketItemArray) {
	var response string
	response, err = c.callPocketAPI("send", list)
	if err != nil {
		debugPrint(fmt.Sprintf("Error in archiveBatch: %s", err.Error()))
		debugPrint(response)
		return err, ret
	}
	var archiveResponse pocketArchiveResponse
	err = json.Unmarshal([]byte(response), &archiveResponse)
	if err != nil {
		debugPrint(fmt.Sprintf("Error in archiveBatch: %s", err.Error()))
		debugPrint(response)
		return err, ret
	}

	if archiveResponse.Status == 0 {
		return fmt.Errorf("Failed to archive some items"), ret
	}
	return err, ret
}

// helper method to call the pocket API via different methods. It just returns
// the body as a string so the caller can decide what the return type is and
// parse it into a Go type as needed.
func (c *PocketClient) callPocketAPI(method string, data interface{}) (ret string, err error) {
	if method != "get" && method != "send" {
		return ret, fmt.Errorf("unknown method: %s", method)
	}
	apiURL := fmt.Sprintf("%s/%s?consumer_key=%s&access_token=%s",
		c.BaseURL, method, c.ConsumerSecret, c.APIToken)
	if data != nil {
		actionsJSON, err := json.Marshal(data)
		if err != nil {
			debugPrint(fmt.Sprintf("Error in callPocketAPI: %s", err.Error()))
			return ret, err
		}
		actions := url.QueryEscape(string(actionsJSON))
		apiURL = fmt.Sprintf("%s&actions=%s", apiURL, actions)
	}
	var response *http.Response
	response, err = c.HTTPClient.Get(apiURL)
	if err != nil {
		debugPrint(fmt.Sprintf("Error in callPocketAPI: %s", err.Error()))
		return ret, err
	}
	if response.StatusCode != 200 {
		errorMsg := fmt.Sprintf("callPocketAPI got status: %s", response.Status)
		debugPrint(errorMsg)
		return ret, errors.New(errorMsg)
	}
	defer response.Body.Close()
	var respBody []byte
	respBody, err = ioutil.ReadAll(response.Body)
	if err != nil {
		debugPrint(fmt.Sprintf("Error in callPocketAPI: %s", err.Error()))
		return ret, err
	}
	ret = string(respBody)
	return ret, err
}

// CleanUpItems is the main method to use this module from. After configuring
// the client with access token and consumer secret and the number of items to
// keep, just run this method and it will clean up your pocket account.
func (c *PocketClient) CleanUpItems() (archived int, err error) {
	var items pocketItemArray
	items, err = c.getAllPocketItems()
	if err != nil {
		debugPrint("Couldn't get pocket items.")
		return 0, err
	}
	itemsCount := len(items)
	debugPrint(fmt.Sprintf("Got %d items in the pocket list.", itemsCount))
	items = filterOutNewestItems(items, c.KeepCount)
	err, items = c.archiveItems(items)
	if err != nil {
		debugPrint("couldn't archive items")
		for _, v := range items {
			log.Println(fmt.Sprintf("failed to archive item: %s", v.GivenTitle))
		}
	}
	return (itemsCount - c.KeepCount), err
}
