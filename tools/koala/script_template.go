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

var window_build_template = `

set GO111MODULE=on
IF NOT EXIST ./go.mod (
go mod init
)

go mod edit  -require=google.golang.org/grpc@v1.24.0

IF NOT EXIST output (
md  output
md conf
)
xcopy  conf output\conf /y /e /i /q
xcopy  scripts output /y /e /i /q

for /F %%i in ('git rev-list -1 HEAD') do ( set GIT_COMMIT=%%i)
for /F %%j in ('go version') do ( set GO_VERSION=%%%j)
set DATE=%date% %time%
cd main
go build  -ldflags "-X 'main.BUILD_TIME=%DATE%' -X 'main.GO_VERSION=%GO_VERSION%' -X 'main.GIT_COMMIT=%GIT_COMMIT%'" -o ../output/bin/{{.PackageName}}.exe

`

var window_start_template = `
bin/{{.PackageName}}.exe
`


