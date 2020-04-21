CGO    = 0
GOOS   = linux
GOARCH = amd64

SRC_PATH = "github.com/matthewpi/ayaya"
OUTPUT_FILE = ayaya

PKG_LIST := $(shell go list ${SRC_PATH}/... | grep -v /vendor/)

all: clean build

clean:
	@rm $(OUTPUT_FILE) -f

build:
	@CGO_ENABLED=$(CGO) GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o $(OUTPUT_FILE) main.go
