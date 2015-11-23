// keep your pocket clean

package main

import (
  "fmt"
  "log"
  "os"
  "github.com/docopt/docopt-go"
  "github.com/Unknwon/goconfig"
  "github.com/mrtazz/pocketcleaner"
)

var (
  debug = false
  version = ""
  buildTime = ""
  builder = ""
  goversion = ""
  cfgfile = fmt.Sprint(os.Getenv("HOME"), "/.pocketcleaner.ini")
  usage = `
pocketcleaner keeps your pocket clean

Usage:
  pocketcleaner
  pocketcleaner [--config=<config>]
  pocketcleaner -d | --debug
  pocketcleaner -h | --help
  pocketcleaner --version

Options:
  -h --help         Show this screen.
  -d --debug        Show debug information.
  --version         Show version.
  --config=<config> Config file to use
`
)

func debugPrint(message string) {
  if debug == true {
    fmt.Println(message)
  }
}

func main() {
      version_info := fmt.Sprint("pocketcleaner ", version,
                                 " built at ", buildTime,
                                 " by ", builder,
                                 " with ", goversion  )
      args, err := docopt.Parse(usage, nil, true, version_info , false)
      if err != nil {
        log.Fatal(err)
      }
      if args["--debug"] == true {
        debug = true
      }
      if args["version"] == true {
        fmt.Println(version_info)
      }
      var consumer_key, access_token string
      debugPrint(fmt.Sprint("Checking if config file '", cfgfile, "' exists."))
      if _, err := os.Stat(cfgfile); err == nil {
        if cfg, err := goconfig.LoadConfigFile(cfgfile); err == nil {
          debugPrint(fmt.Sprint("Config file '", cfgfile, "' loaded."))
          if key, err := cfg.GetValue(goconfig.DEFAULT_SECTION, "consumer_key"); err ==nil {
            consumer_key = key
          }
          if token, err := cfg.GetValue(goconfig.DEFAULT_SECTION, "access_token"); err ==nil {
            access_token = token
          }
        }
      }
      pocketcleaner.CallPocketAPI("foo", consumer_key, access_token)
}

