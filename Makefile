VERSION = v0.1.0
PHASE = beta
FULL_VERSION := $(shell git describe --exact-match --tags HEAD 2>/dev/null || \
	echo $(VERSION)-$(PHASE)+$$(git rev-parse --short HEAD)) 

GOLANGCI_LINT_VERSION = v2.5.0

.PHONY: build
build:
	@echo "Building version: $(FULL_VERSION)"; 
	@go build -ldflags "-X main.version=$(FULL_VERSION)" -o bin/app .

.PHONY: run
run: build
	@./bin/app

.PHONY: linter
linter:
	@CMD="curl -sSfL \
		https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh \
		| sh -s -- -b $(CURDIR)/bin $(GOLANGCI_LINT_VERSION)"; \
	if [ ! -f $(CURDIR)/bin/golangci-lint ]; then \
		echo "未找到 golangci-lint，正在安装..."; \
		eval $$CMD; \
	else \
		INSTALLED=$$($(CURDIR)/bin/golangci-lint version 2>/dev/null \
			| grep -oE 'version [0-9.]+' | grep -oE '[0-9.]+' || echo ""); \
		EXPECTED=$$(echo $(GOLANGCI_LINT_VERSION) | sed 's/v//'); \
		if [ "$$INSTALLED" != "$$EXPECTED" ] || [ -z "$$INSTALLED" ]; then \
			echo "版本不匹配 (已安装: $$INSTALLED, 需要: $$EXPECTED)"; \
			echo "正在重新安装..."; \
			eval $$CMD; \
		else \
			echo "golangci-lint $(GOLANGCI_LINT_VERSION) 已安装"; \
		fi; \
	fi 

.PHONY: lint
lint: linter
	@$(CURDIR)/bin/golangci-lint run --fix