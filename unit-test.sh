go test -v ./... --coverprofile cover.out
go tool cover -func cover.out
