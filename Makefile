#
# Default Target
#
.PHONY:
all: fmt test
	go build -o build ./...

#
# Format Target
#
.PHONY:
fmt: .build.prepare
	gofmt -s -w .
	go tool goimports -w -l .

#
# Test Target
#
.PHONY:
test: .build.prepare
	go tool gotestsum --junitfile build/tests.xml --format testdox --format-hide-empty-pkg -- -count=1 -coverprofile=build/coverage.out -covermode count ./...
	go tool cover -func build/coverage.out | tail -n 1 | tr -s '[:blank:]'
	go tool cover -html=build/coverage.out -o build/coverage.html

#
# Clean Target
#
# Remove build artifacts.
.PHONY:
clean:
	rm -rf build

#
# [internal] Build Prepare Target
#
.PHONY:
.build.prepare:
	go mod download
	go mod verify
	go mod tidy
	mkdir -p build
