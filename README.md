gopack
======

## Scope

Ultra easy tool to create go-files with static content from directories. This may be usefull to embedd static content into go binaries.

## Install
```bash
> go get -v github.com/contiamo/gopack
```

## Usage
```bash
# prepare sample directory
> mkdir -p foo/bar/baz
# create content which will be embedded
> echo "example content" > foo/bar/content.txt
# create content which will be excluded
> echo "other sample content" > foo/bar/baz/other-content.txt
# call gopack
> gopack \
  -dir foo/bar \
  -package mypackage \
  -variable staticContent \
  -output output.go \
  -exclude foo/bar/baz,.git \
  -with-http-handler
# inspect output
> cat output.go
package mypackage

import "github.com/contiamo/gopack/staticserver"

var staticContent = map[string][]byte{
  "foo/bar/content.txt": {
    101,120,97,109,112,108,101,32,99,111,110,116,101,110,116,10,
  },
}

var staticContentHandler = staticserver.New(staticContent)
```

## Example
This will create a pack from the data directory, compile a little webserver using that pack, start it and call it to see all the results and mime types being returned.
```bash
> bash example/run-example.sh
~/go/src/github.com/contiamo/gopack/data ~/go/src/github.com/contiamo/gopack
~/go/src/github.com/contiamo/gopack
~/go/src/github.com/contiamo/gopack/example ~/go/src/github.com/contiamo/gopack
HTTP/1.1 200 OK
Content-Type: text/plain; charset=utf-8
Date: Tue, 11 Sep 2018 11:20:09 GMT
Content-Length: 4

HTTP/1.1 200 OK
Content-Type: application/json
Date: Tue, 11 Sep 2018 11:20:09 GMT
Content-Length: 3

HTTP/1.1 200 OK
Date: Tue, 11 Sep 2018 11:20:09 GMT
Content-Length: 3
Content-Type: text/plain; charset=utf-8

~/go/src/github.com/contiamo/gopack
```
