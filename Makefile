graphql_files = pkg/graph/exec.go pkg/graph/model.go

BIN := ghreviews-server
OUTBIN = bin/$(BIN)

default: run

$(graphql_files): gqlgen.yml
	go run github.com/99designs/gqlgen --verbose

graphql:
	go run github.com/99designs/gqlgen --verbose

build: graphql
	go build -o $(OUTBIN) github.com/sebastianvera/ghreviews/cmd/server

run: $(graphql_files)
	go run cmd/server/server.go

clean:
	-rm -f $(graphql_files)
	-rm -f $(OUTBIN)

.PHONY: run build clean graphql
