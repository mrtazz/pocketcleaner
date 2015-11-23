package pocketcleaner_test

import (
  "testing"
  "github.com/mrtazz/pocketcleaner"
)

func TestParsePacketResponse(t *testing.T) {
  var input = `
{
   "search_meta" : {
      "search_type" : "normal"
   },
   "status" : 1,
   "complete" : 1,
   "list" : {
      "485722939" : {
         "word_count" : "711",
         "item_id" : "485722939",
         "is_article" : "1",
         "time_added" : "1385319303",
         "resolved_title" : "Bootstrapping Business: Grow Your Company Without VC Funding",
         "status" : "0",
         "favorite" : "0",
         "time_updated" : "1385319303",
         "given_title" : "Silicon Allee » Bootstrapping Business: Grow Your Company Without VC Fundin",
         "sort_id" : 2048,
         "excerpt" : "In my 30 plus years of business, I’ve been an executive, advisor or investor in over two dozen companies.",
         "resolved_id" : "485722939",
         "has_image" : "1",
         "is_index" : "0",
         "given_url" : "http://siliconallee.com/silicon-allee/editorial/2013/11/22/bootstrapping-business-grow-your-company-without-vc-funding",
         "has_video" : "0",
         "resolved_url" : "http://siliconallee.com/silicon-allee/editorial/2013/11/22/bootstrapping-business-grow-your-company-without-vc-funding",
         "time_read" : "0",
         "time_favorited" : "0"
      },
      "1029402404" : {
         "is_article" : "1",
         "word_count" : "2470",
         "item_id" : "1029402404",
         "time_added" : "1441297590",
         "resolved_title" : "One of the biggest mistakes I’ve made in my career",
         "status" : "0",
         "time_updated" : "1441417608",
         "favorite" : "0",
         "sort_id" : 214,
         "given_title" : "One of the biggest mistakes I’ve made in my career — Twenty Years in the Va",
         "excerpt" : "Should designers in high-tech learn how to code?  It’s a question that has cropped up consistently over the past two decades, creating scores of heated debate at design conferences, as well as on every design website or blog you can find.",
         "resolved_id" : "1029402404",
         "has_image" : "0",
         "is_index" : "0",
         "resolved_url" : "https://medium.com/twenty-years-in-the-valley/one-of-the-biggest-mistakes-i-ve-made-in-my-career-72bf27c538b4",
         "has_video" : "0",
         "given_url" : "https://medium.com/twenty-years-in-the-valley/one-of-the-biggest-mistakes-i-ve-made-in-my-career-72bf27c538b4",
         "time_read" : "0",
         "time_favorited" : "0"
      },
      "839271306" : {
         "status" : "0",
         "resolved_title" : "Stop Blowhard Syndrome",
         "item_id" : "839271306",
         "word_count" : "319",
         "is_article" : "1",
         "time_added" : "1426698487",
         "given_title" : "That's What Xu Said : Stop Blowhard Syndrome",
         "sort_id" : 894,
         "favorite" : "0",
         "time_updated" : "1426698487",
         "has_image" : "0",
         "is_index" : "0",
         "excerpt" : "When I express any shred of doubt about whether I deserve or am qualified for something, people often try to reassure me that I am just experiencing impostor syndrome. About 10% of the time, it’s true.",
         "resolved_id" : "839271306",
         "time_read" : "0",
         "time_favorited" : "0",
         "has_video" : "0",
         "given_url" : "http://xuhulk.tumblr.com/post/110549967516/stop-blowhard-syndrome",
         "resolved_url" : "http://xuhulk.tumblr.com/post/110549967516/stop-blowhard-syndrome"
      }
   },
   "error" : null,
   "since" : 1448244422
}
  `
  ret, err := pocketcleaner.ParsePocketResponse(input)
  if err != nil {
    t.Errorf("ParsePocketResponse: parse failed: %s", err.Error())
  }
  if ret.Since != 1448244422 {
    t.Errorf("ParsePocketResponse: expected %d, actual %d", 1448244422, ret.Since)
  }
  if len(ret.List) != 3 {
    t.Errorf("ParsePocketResponse: expected %d, actual %d", 3, len(ret.List))
  }
  item := ret.List["839271306"]
  if item.GivenTitle != "That's What Xu Said : Stop Blowhard Syndrome" {
    t.Errorf("ParsePocketResponse: expected %s, actual %s", "That's What Xu Said : Stop Blowhard Syndrome", item.GivenTitle)
  }
}
