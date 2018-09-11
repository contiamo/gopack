#!/bin/bash


pushd ${GOPATH}/src/github.com/contiamo/gopack/data
gopack --output ../example/data.go --with-http-handler
popd

pushd ${GOPATH}/src/github.com/contiamo/gopack/example
go build -o example.out main.go data.go
./example.out &
EXAMPLE_PID=$!
sleep 0.5
curl -I localhost:8080/foo/bar/textfile.txt
curl -I localhost:8080/foo/bar/jsonfile.json
curl -I localhost:8080/foo/bar/jsonfile.wrong
kill $EXAMPLE_PID
popd

exit 0


