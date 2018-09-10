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
