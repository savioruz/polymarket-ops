.PHONY: build watch-reports watch-report-one scan-report scan-report-build dashboard test-dashboard setup

build:
	mkdir -p bin
	cd dashboard && go build -o ../bin/dashboard ./cmd/dashboard

run:
	./bin/dashboard

watch-reports:
	cd dashboard && go run ./cmd/watch-reports --root ..

watch-report-one:
	cd dashboard && go run ./cmd/watch-reports --root .. --wallet $(WALLET)

scan-report:
	cd dashboard && go run ./cmd/scan-report --root ..

scan-report-build:
	mkdir -p bin
	cd dashboard && go build -o ../bin/scan-report ./cmd/scan-report

dashboard:
	cd dashboard && go run ./cmd/dashboard ..

test:
	cd dashboard && go test ./...

setup:
	cd dashboard && go mod download && go mod tidy
