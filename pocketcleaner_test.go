package pocketcleaner

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
)

type pocketMock struct {
	JSON      string
	Response  pocketResponse
	ItemArray pocketItemArray
	Client    *PocketClient
	Server    *httptest.Server
}

type mockSetup struct {
	Code int
	Body string
}

// gratefully taken and adapted from
// http://keighl.com/post/mocking-http-responses-in-golang/
func expect(t *testing.T, a interface{}, b interface{}) {
	if reflect.DeepEqual(a, b) == false {
		t.Errorf("Expected: %v (type %v) - Got: %v (type %v)", b, reflect.TypeOf(b), a, reflect.TypeOf(a))
	}
}

// this is the test setup method that each test can call to get a set of mock
// objects. If arguments are passed, the first one should be the return code
// of the mock server and the second one the return body
func testSetup(ms mockSetup) pocketMock {
	input, _ := ioutil.ReadFile("fixtures/pocket_response.json")
	mockedJSON := string(input)
	mockedResponse, _ := parsePocketResponse(mockedJSON)
	mockedPocketItemArray := make(pocketItemArray, 0, len(mockedResponse.List))

	for _, v := range mockedResponse.List {
		mockedPocketItemArray = append(mockedPocketItemArray, v)
	}

	var body string
	code := ms.Code
	if ms.Body != "" {
		body = ms.Body
	} else {
		body = mockedJSON
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(code)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, body)
	}))

	transport := &http.Transport{
		Proxy: func(req *http.Request) (*url.URL, error) {
			return url.Parse(server.URL)
		},
	}

	httpClient := &http.Client{Transport: transport}
	client := &PocketClient{
		BaseURL:        server.URL,
		ConsumerSecret: "foo",
		APIToken:       "bar",
		HTTPClient:     httpClient,
		KeepCount:      5,
	}

	return pocketMock{
		mockedJSON, mockedResponse, mockedPocketItemArray, client, server,
	}
}

func TestParsePacketResponse(t *testing.T) {
	input, _ := ioutil.ReadFile("fixtures/pocket_response.json")
	ret, err := parsePocketResponse(string(input))

	expect(t, err, nil)
	expect(t, ret.Since, uint(1448244422))
	expect(t, len(ret.List), 16)
	item := ret.List["839271306"]
	expect(t, item.GivenTitle, "That's What Xu Said : Stop Blowhard Syndrome")
}

func TestFilterOutNewestItems(t *testing.T) {
	mockedPocketItemArray := testSetup(mockSetup{}).ItemArray
	ret := filterOutNewestItems(mockedPocketItemArray, 5)

	expect(t, len(ret), 11)

	var flagtests = []struct {
		id         int
		timestamp  uint64
		givenTitle string
	}{
		{0, 1385319303, "Silicon Allee » Bootstrapping Business: Grow Your Company Without VC Fundin"},
		{1, 1394546312, "Your Marriage Will Fail, by Alicia Liu | Model View Culture"},
		{2, 1400039331, ""},
		{3, 1404071597, "http://www.vox.com/2014/6/26/5837638/the-ipo-is-dying-marc-andreessen-expla"},
		{4, 1410813125, "published a fascinating writeup"},
		{5, 1414509648, "Project Managing Your Health — Medium"},
		{6, 1415733610, "What It's Like To Burn Out - Career Burnout - Elle"},
		{7, 1421113507, "Towards Better Interviews – Venkata Mahalingam"},
		{8, 1421515685, "Why Remote Engineering Is So Difficult blog.learningbyshipping.com/2014/12/"},
		{9, 1422935881, ""},
		{10, 1426045670, "Meeting with purpose derrickbradley.github.io/2015/02/20/mee…"},
	}

	for _, tt := range flagtests {
		expect(t, ret[tt.id].GivenTitle, tt.givenTitle)
		expect(t, ret[tt.id].TimeAdded, tt.timestamp)
	}

	ret2 := filterOutNewestItems(mockedPocketItemArray, 20)
	expect(t, len(ret2), 0)

}

func TestArchiveItems(t *testing.T) {
	pm := testSetup(mockSetup{200, `{"action_results":[],"status":1}`})
	mockedPocketItemArray, server, client := pm.ItemArray, pm.Server, pm.Client
	defer server.Close()
	err, _ := client.archiveItems(mockedPocketItemArray)
	expect(t, err, nil)
}

func TestArchiveItemsParseError(t *testing.T) {
	pm := testSetup(mockSetup{200, `tatus":1}`})
	mockedPocketItemArray, server, client := pm.ItemArray, pm.Server, pm.Client
	defer server.Close()
	err, _ := client.archiveItems(mockedPocketItemArray)
	expect(t, err.Error(), "invalid character 'a' in literal true (expecting 'r')")
}

func TestArchiveItemsFailed(t *testing.T) {
	pm := testSetup(mockSetup{200, `{"action_results":[],"status":0}`})
	mockedPocketItemArray, server, client := pm.ItemArray, pm.Server, pm.Client
	defer server.Close()
	err, _ := client.archiveItems(mockedPocketItemArray)
	expect(t, err.Error(), "Failed to archive some items")
}

func TestCallPocketApiDefaultNotImplemented(t *testing.T) {
	pm := testSetup(mockSetup{200, "foo"})
	server, client := pm.Server, pm.Client
	defer server.Close()
	_, err := client.callPocketAPI("foo", nil)
	expect(t, err.Error(), "unknown method: foo")
}

func TestCallPocketApiGet(t *testing.T) {
	pm := testSetup(mockSetup{200, ""})
	server, client := pm.Server, pm.Client
	defer server.Close()
	ret, err := client.callPocketAPI("get", nil)
	expect(t, err, nil)
	resp, _ := parsePocketResponse(ret)
	expect(t, resp, pm.Response)
}

func TestCallPocketApiGetFailed(t *testing.T) {
	pm := testSetup(mockSetup{404, ""})
	server, client := pm.Server, pm.Client
	defer server.Close()
	ret, err := client.callPocketAPI("get", nil)
	expect(t, err, nil)
	resp, _ := parsePocketResponse(ret)
	expect(t, resp, pm.Response)
}

func TestCallPocketApiSend(t *testing.T) {
	pm := testSetup(mockSetup{200, "foo"})
	server, client := pm.Server, pm.Client
	defer server.Close()
	_, err := client.callPocketAPI("send", nil)
	expect(t, err, nil)
}

func TestCallPocketApiGetHTTPError(t *testing.T) {
	pm := testSetup(mockSetup{200, "foo"})
	server, client := pm.Server, pm.Client
	server.Close()
	_, err := client.callPocketAPI("get", nil)
	var urlError *url.Error
	expect(t, reflect.TypeOf(err), reflect.TypeOf(urlError))
}

func TestCallPocketApiSendParseError(t *testing.T) {
	pm := testSetup(mockSetup{200, "foo"})
	server, client := pm.Server, pm.Client
	defer server.Close()
	_, err := client.callPocketAPI("send", client)
	expect(t, err.Error(), "json: unsupported type: func(*http.Request) (*url.URL, error)")
}

func TestGetAllPocketItems(t *testing.T) {
	pm := testSetup(mockSetup{200, ""})
	server, client := pm.Server, pm.Client
	defer server.Close()
	items, err := client.getAllPocketItems()
	expect(t, err, nil)
	expect(t, len(items), 16)
}

func TestGetAllPocketItemsApiError(t *testing.T) {
	pm := testSetup(mockSetup{200, ""})
	server, client := pm.Server, pm.Client
	server.Close()
	items, err := client.getAllPocketItems()
	var urlError *url.Error
	expect(t, reflect.TypeOf(err), reflect.TypeOf(urlError))
	expect(t, len(items), 0)
}

func TestGetAllPocketItemsParseError(t *testing.T) {
	pm := testSetup(mockSetup{200, "f"})
	server, client := pm.Server, pm.Client
	defer server.Close()
	items, err := client.getAllPocketItems()
	expect(t, err.Error(), "invalid character '\\n' in literal false (expecting 'a')")
	expect(t, len(items), 0)
}

func TestPocketClientWithToken(t *testing.T) {
	ret := PocketClientWithToken("bar", "foo", 5)
	expect(t, "foo", ret.ConsumerSecret)
	expect(t, "bar", ret.APIToken)
	expect(t, 5, ret.KeepCount)
}

func TestCleanUpItems(t *testing.T) {
	pm := testSetup(mockSetup{200, ""})
	server, client := pm.Server, pm.Client
	defer server.Close()
	err := client.CleanUpItems()
	expect(t, err, nil)
}

func TestPocketArchiveItemJSON(t *testing.T) {
	item := &pocketArchiveItem{"archive", "123", "123456789"}
	ret, err := json.Marshal(item)
	expect(t, err, nil)
	expect(t, string(ret), `{"action":"archive","item_id":"123","time":"123456789"}`)
}

func TestPocketArchiveResponseJSONParsing(t *testing.T) {
	item := &pocketArchiveResponse{[]bool{true, true}, 1}
	ret, err := json.Marshal(item)
	expect(t, err, nil)
	expect(t, string(ret), `{"action_results":[true,true],"status":1}`)
}
