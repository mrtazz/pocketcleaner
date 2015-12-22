// keep your pocket clean

package main

import (
	"fmt"
	"github.com/Unknwon/goconfig"
	"github.com/docopt/docopt-go"
	"github.com/mrtazz/pocketcleaner"
	"log"
	"os"
	"strconv"
)

var (
	isDebug   = false
	version   = ""
	buildTime = ""
	builder   = ""
	goversion = ""
	cfgfile   = fmt.Sprint(os.Getenv("HOME"), "/.pocketcleaner.ini")
	usage     = `
pocketcleaner keeps your pocket clean

Usage:
  pocketcleaner [-d | --debug] [--keep=<keepCount>]
  pocketcleaner -h | --help | -v | --version

Options:
  -h --help          Show this screen.
  -d --debug         Show debug information.
  -v --version       Show version.
  --config=<config>  Config file to use
  --keep=<keepCount> Count of items to keep
`
)

func debugPrint(message string) {
	if isDebug == true {
		fmt.Println(message)
	}
}

func main() {
	versionInfo := fmt.Sprint("pocketcleaner ", version,
		" built at ", buildTime,
		" by ", builder,
		" with ", goversion)
	args, err := docopt.Parse(usage, nil, true, versionInfo, false)
	if err != nil {
		log.Fatal(err)
	}
	if args["--debug"] == true {
		isDebug = true
	}
	var consumerKey, accessToken string
	var keepCount int
	debugPrint(fmt.Sprint("Checking if config file '", cfgfile, "' exists."))
	if _, err := os.Stat(cfgfile); err == nil {
		if cfg, err := goconfig.LoadConfigFile(cfgfile); err == nil {
			debugPrint(fmt.Sprint("Config file '", cfgfile, "' loaded."))
			if key, err := cfg.GetValue(goconfig.DEFAULT_SECTION, "consumer_key"); err == nil {
				consumerKey = key
			} else {
				log.Fatal(err)
			}
			if token, err := cfg.GetValue(goconfig.DEFAULT_SECTION, "access_token"); err == nil {
				accessToken = token
			} else {
				log.Fatal(err)
			}
			if count, err := cfg.GetValue(goconfig.DEFAULT_SECTION, "keep_count"); err == nil {
				if keepCount, err = strconv.Atoi(count); err != nil {
					log.Fatal(err)
				}
			} else {
				log.Fatal(err)
			}
		}
	}
	if count, ok := args["--keep"].(string); ok {
		if keepCount, err = strconv.Atoi(count); err != nil {
			log.Fatal(err)
		}
	}
	pocket := pocketcleaner.PocketClientWithToken(accessToken, consumerKey, keepCount)
	pocketcleaner.Debug = isDebug
	debugPrint(fmt.Sprintf("Cleaning up items, keeping the newest %d...", keepCount))
	var archived int
	if archived, err = pocket.CleanUpItems(); err != nil {
		log.Fatal(err)
	} else {
		log.Printf("Succesfully archived %d pocket items.\n", archived)
	}
}
