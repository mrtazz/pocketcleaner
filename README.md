# pocket-cleaner

[![Build Status](https://travis-ci.org/mrtazz/pocketcleaner.svg?branch=master)](https://travis-ci.org/mrtazz/pocket-cleaner)
[![Coverage Status](https://coveralls.io/repos/mrtazz/pocketcleaner/badge.svg?branch=master&service=github)](https://coveralls.io/github/mrtazz/pocket-cleaner?branch=master)
[![Packagecloud](https://img.shields.io/badge/packagecloud-available-green.svg)](https://packagecloud.io/mrtazz/pocketcleaner)

## Overview

This is a utility to keep your [Pocket](https://getpocket.com) list small and
manageable. It will archive all items in your list except for the newest `n`
items.

**WARNING**: This will *archive* items. Something you can't just undo.

## Usage
```
pocketcleaner
pocketcleaner [--config=<config>]
pocketcleaner [--keep=<keepCount>]
pocketcleaner -d | --debug
pocketcleaner -h | --help
pocketcleaner --version

Options:
  -h --help          Show this screen.
  -d --debug         Show debug information.
  --version          Show version.
  --config=<config>  Config file to use
  --keep=<keepCount> Count of items to keep
```

## Installation

### Get auth tokens to use with pocketcleaner

Pocketcleaner doesn't come with any auth keys, so you will have to get your
own. The following steps are taken from
[here](http://www.jamesfmackenzie.com/getting-started-with-the-pocket-developer-api/).

1. Create an app in the [developer portal](http://getpocket.com/developer/)
Make sure it has at least the `retrieve` and `modify` permissions.

2. Get a request token
```
curl https://getpocket.com/v3/oauth/request -X POST \
-H "Content-Type: application/json" -H "X-Accept: application/json" \
-d '{"consumer_key":"your-key-here","redirect_uri":"http://www.google.com"}'
```
3. Authorize the app
Visit the following URL in your browser and authorize the app.
```
https://getpocket.com/auth/authorize?request_token=request-token-here&redirect_uri=http://www.google.com
```
4. Convert request token to access token
```
curl https://getpocket.com/v3/oauth/authorize -X POST \
-H "Content-Type: application/json" -H "X-Accept: application/json" \
-d '{"consumer_key":"your-key-here","code":"request-token"}'
```
5. Make sure the token works by getting a list of your saved items
```
curl http://getpocket.com/v3/get -X POST -H "Content-Type: application/json" \
-H "X-Accept: application/json" \
-d '{"consumer_key":"your-key-here", "access_token":"access-token"}'
```

### Install pocketcleaner

There are rpms for linux up [on
packagecloud.io][https://packagecloud.io/mrtazz/pocketcleaner] but you can
also just

```
go install github.com/mrtazz/pocketcleaner
```

or clone the repo and build the binary yourself

```
go get github.com/mrtazz/pocketcleaner
cd $GOPATH/src/github.com/mrtazz/pocketcleaner
make
```
