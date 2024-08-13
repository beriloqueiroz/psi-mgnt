# Makefile
.PHONY: build

BINARY_NAME=psi-mgnt

# build builds the tailwind css sheet, and compiles the binary into a usable thing.
build:
	go mod tidy && \
   	templ generate && \
	go generate cmd/main.go && \
	npx tailwindcss build -i static/css/style.css -o static/css/tailwind.css -m && \
	go build -ldflags="-w -s" -o ${BINARY_NAME} ./cmd/main.go

# dev runs the development server where it builds the tailwind css sheet,
# and compiles the project whenever a file is changed.
dev:
	npx tailwindcss build -i static/css/style.css -o static/css/tailwind.css -m &\
	templ generate --watch --cmd="go generate cmd/main.go" &\
	templ generate --watch --cmd="go run cmd/main.go"

clean:
	go clean