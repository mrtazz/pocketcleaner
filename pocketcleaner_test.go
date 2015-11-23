package pocketcleaner_test

import (
  "testing"
  "io/ioutil"
  "github.com/mrtazz/pocketcleaner"
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
