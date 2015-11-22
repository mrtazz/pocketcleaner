// keep your pocket clean

package main

import (
  "fmt"
  "log"
  "github.com/docopt/docopt-go"
)

var (
  debug = false
  version = ""
  buildTime = ""
  builder = ""
  goversion = ""
  usage = `
pocket-cleaner keeps your pocket clean

Usage:
  pocket-cleaner
  pocket-cleaner -d | --debug
  pocket-cleaner -h | --help
  pocket-cleaner --version

Options:
  -h --help     Show this screen.
  -d --debug    Show debug information.
  --version     Show version.
`
)

func main() {
      version_info := fmt.Sprint("pocket-cleaner ", version,
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
}
