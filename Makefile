.DEFAULT_GOAL := help

DASHBOARD_DIR := dashboard
BIN_DIR := bin
GO := go

.PHONY: help build dashboard-build run dashboard dashboard-run watch-reports watch-report-one scan-report scan-report-build test setup

help:
	@printf "Available targets:\n"
	@printf "  make setup             Download/tidy dashboard Go modules\n"
	@printf "  make build             Build dashboard binary to bin/dashboard\n"
	@printf "  make run               Run built dashboard binary\n"
	@printf "  make dashboard         Run dashboard directly with go run\n"
	@printf "  make watch-reports     Generate watch reports for all wallets\n"
	@printf "  make watch-report-one  Generate watch report for one wallet (set WALLET=0x...)\n"
	@printf "  make scan-report       Generate leaderboard scan report\n"
	@printf "  make scan-report-build Build scan-report binary to bin/scan-report\n"
	@printf "  make test              Run dashboard Go tests\n"

build: dashboard-build

dashboard-build:
	mkdir -p $(BIN_DIR)
	$(GO) -C $(DASHBOARD_DIR) build -o ../$(BIN_DIR)/dashboard ./cmd/dashboard

run:
	./$(BIN_DIR)/dashboard

dashboard: dashboard-run

dashboard-run:
	$(GO) -C $(DASHBOARD_DIR) run ./cmd/dashboard ..

watch-reports:
	$(GO) -C $(DASHBOARD_DIR) run ./cmd/watch-reports --root ..

watch-report-one:
	$(GO) -C $(DASHBOARD_DIR) run ./cmd/watch-reports --root .. --wallet $(WALLET)

scan-report:
	$(GO) -C $(DASHBOARD_DIR) run ./cmd/scan-report --root ..

scan-report-build:
	mkdir -p $(BIN_DIR)
	$(GO) -C $(DASHBOARD_DIR) build -o ../$(BIN_DIR)/scan-report ./cmd/scan-report

test:
	$(GO) -C $(DASHBOARD_DIR) test ./...

setup:
	$(GO) -C $(DASHBOARD_DIR) mod download
	$(GO) -C $(DASHBOARD_DIR) mod tidy
