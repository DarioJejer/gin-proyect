# =========================
# Configuración
# =========================

APP_NAME := gin-app
PKGS := $(shell go list ./...)
COVER_PKGS := $(shell go list ./... | grep -v migrations | grep -v mocks | grep -v docs | grep -v routes | grep -v initializers | tr '\n' ',')
COVER_FILE := coverage.out
COVER_FILTERED := coverage_filtered.out

# =========================
# Comandos principales
# =========================

.PHONY: all
all: test

# Run all tests
.PHONY: test
test:
	go test ./...

# =========================
# Coverage
# =========================

# Generate coverage file
.PHONY: coverage
coverage:
	go test ./... -coverpkg=$(COVER_PKGS) -coverprofile=$(COVER_FILE)

# Filter unwanted files (main.go, migrations, docs, mocks)
.PHONY: coverage-filter
coverage-filter: coverage
	grep -v -E "main.go|helpers/swaggerConfig.go" $(COVER_FILE) > $(COVER_FILTERED)

# Show coverage in terminal
.PHONY: coverage-report
coverage-report: coverage-filter
	go tool cover -func=$(COVER_FILTERED)

# Open coverage in browser
.PHONY: coverage-html
coverage-html: coverage-filter
	go tool cover -html=$(COVER_FILTERED)

# =========================
# Extras útiles
# =========================

# Run tests with verbose output
.PHONY: test-verbose
test-verbose:
	go test -v ./...

# Run tests for a specific package
.PHONY: test-pkg
test-pkg:
	go test -v $(PKG)

# Clean coverage files
.PHONY: clean
clean:
	rm -f $(COVER_FILE) $(COVER_FILTERED)

# Tidy dependencies
.PHONY: tidy
tidy:
	go mod tidy

# Format code
.PHONY: fmt
fmt:
	go fmt ./...

# Lint (si tenés golangci-lint)
.PHONY: lint
lint:
	golangci-lint run
