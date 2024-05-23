# Makefile
.PHONY: build

BINARY_NAME=trello-app

# build builds the tailwind css sheet, and compiles the binary into a usable thing.
build:
	go mod tidy && \
   	templ generate && \
	go generate ./cmd/main.go && \
	go build -ldflags="-w -s" -o ${BINARY_NAME} ./cmd/main.go

# dev runs the development server where it builds the tailwind css sheet,
# and compiles the project whenever a file is changed.
dev:
	templ generate --watch --cmd="go generate" &\
	templ generate --watch --cmd="go run ."

clean:
	go clean