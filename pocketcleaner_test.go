package pocketcleaner_test

import (
	"github.com/mrtazz/pocketcleaner"
	"io/ioutil"
	"testing"
)

func TestParsePacketResponse(t *testing.T) {
	input, _ := ioutil.ReadFile("fixtures/pocket_response.json")
	ret, err := pocketcleaner.ParsePocketResponse(string(input))
	if err != nil {
		t.Errorf("ParsePocketResponse: parse failed: %s", err.Error())
	}
	if ret.Since != 1448244422 {
		t.Errorf("ParsePocketResponse: expected %d, actual %d", 1448244422, ret.Since)
	}
	if len(ret.List) != 16 {
		t.Errorf("ParsePocketResponse: expected %d, actual %d", 16, len(ret.List))
	}
	item := ret.List["839271306"]
	if item.GivenTitle != "That's What Xu Said : Stop Blowhard Syndrome" {
		t.Errorf("ParsePocketResponse: expected %s, actual %s", "That's What Xu Said : Stop Blowhard Syndrome", item.GivenTitle)
	}
}

func TestFilterOutNewestItems(t *testing.T) {
	input, _ := ioutil.ReadFile("fixtures/pocket_response.json")
	items, err := pocketcleaner.ParsePocketResponse(string(input))
	if err != nil {
		t.Errorf("ParsePocketResponse: parse failed: %s", err.Error())
	}
	arr := make(pocketcleaner.PocketItemArray, 0)

	for _, v := range items.List {
		arr = append(arr, v)
	}

	ret := pocketcleaner.FilterOutNewestItems(arr, 5)

	if len(ret) != 11 {
		t.Errorf("FilterOutNewestItems: expected: %d, actual: %d", 11, len(ret))
	}

	// TODO: test that the right ones got filtered
	var flagtests = []struct {
		id          int
		timestamp   string
		given_title string
	}{
		{0, "1385319303", "Silicon Allee » Bootstrapping Business: Grow Your Company Without VC Fundin"},
		{1, "1394546312", "Your Marriage Will Fail, by Alicia Liu | Model View Culture"},
		{2, "1400039331", ""},
		{3, "1404071597", "http://www.vox.com/2014/6/26/5837638/the-ipo-is-dying-marc-andreessen-expla"},
		{4, "1410813125", "published a fascinating writeup"},
		{5, "1414509648", "Project Managing Your Health — Medium"},
		{6, "1415733610", "What It's Like To Burn Out - Career Burnout - Elle"},
		{7, "1421113507", "Towards Better Interviews – Venkata Mahalingam"},
		{8, "1421515685", "Why Remote Engineering Is So Difficult blog.learningbyshipping.com/2014/12/"},
		{9, "1422935881", ""},
		{10, "1426045670", "Meeting with purpose derrickbradley.github.io/2015/02/20/mee…"},
	}

	for _, tt := range flagtests {
		if ret[tt.id].GivenTitle != tt.given_title {
			t.Errorf("FilterOutNewestItems Title: expected: %s, actual: %s, at index: %d", tt.given_title, ret[tt.id].GivenTitle, tt.id)
			t.Errorf("FilterOutNewestItems Timestamp: expected: %s, actual: %d, at index: %d", tt.timestamp, ret[tt.id].TimeAdded, tt.id)
		}
	}

}
