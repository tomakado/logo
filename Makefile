test:
	go test ./... -v

test-cover:
	go test ./... -v -coverprofile=cover.out &&\
	go tool cover -func=cover.out

test-cover-html:
	go test ./... -v -coverprofile=cover.out &&\
	go tool cover -html=cover.out

doc:
	godoc -http=:6060