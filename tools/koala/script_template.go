package main

var build_template = `
#!/bin/bash

export GO111MODULE=on
if [ ! -f ./go.mod ]
then
go mod init
fi

go mod edit  -require=google.golang.org/grpc@v1.24.0

mkdir -p output
cp -rf conf output/
cp -rf ./scripts/* output/

export GIT_COMMIT=$(git rev-list -1 HEAD)
cd main
go build  -ldflags "-X 'main.BUILD_TIME=` + "`date`" + `' -X 'main.GO_VERSION=` + "`go version`" + `' -X 'main.GIT_COMMIT=$GIT_COMMIT'" -o ../output/bin/{{.PackageName}}
`

var start_template = `
#!/bin/bash
./bin/{{.PackageName}}
`
