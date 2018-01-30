GOOS=linux GOARCH=386
go build -o dist/gitmon-$GOOS-$GOARCH

GOOS=linux GOARCH=amd64
go build -o dist/gitmon-$GOOS-$GOARCH


GOOS=darwin GOARCH=386
go build -o dist/gitmon-$GOOS-$GOARCH

GOOS=darwin GOARCH=amd64
go build -o dist/gitmon-$GOOS-$GOARCH


GOOS=windows GOARCH=386
go build -o dist/gitmon-$GOOS-$GOARCH

GOOS=windows GOARCH=amd64
go build -o dist/gitmon-$GOOS-$GOARCH