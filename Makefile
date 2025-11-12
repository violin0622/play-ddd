VERSION = v0.1.0
PHASE = beta
FULL_VERSION := $(shell git describe --exact-match --tags HEAD 2>/dev/null || \
	echo $(VERSION)-$(PHASE)+$$(git rev-parse --short HEAD)) 

# 构建信息变量
GIT_SHA := $(shell git rev-parse HEAD 2>/dev/null || echo "unknown")
BUILD_TIME := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
GO_VERSION := $(shell go version | awk '{print $$3}')

# 探测本地操作系统和架构
OS := $(shell uname -s)
ARCH := $(shell uname -m)
BUF_VERSION = 1.59.0
GOLANGCI_LINT_VERSION = v2.5.0

.PHONY: build
build:
	@echo "Building version: $(FULL_VERSION)"; 
	go build -ldflags "\
		-X play-ddd/cmd.Version=$(FULL_VERSION) \
		-X play-ddd/cmd.GitSHA=$(GIT_SHA) \
		-X play-ddd/cmd.BuiltAt=$(BUILD_TIME) \
		-X play-ddd/cmd.GoVersion=$(GO_VERSION)" \
		-o bin/app .

.PHONY: run
run: build
	@./bin/app

.PHONY: golangci-lint
golangci-lint:
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

.PHONY: buf
buf:
	@BUF_URL="https://github.com/bufbuild/buf/releases/download/v$(BUF_VERSION)/buf-$(OS)-$(ARCH)"; \
	BUF_BIN="$(CURDIR)/bin/buf"; \
	if [ ! -f $$BUF_BIN ]; then \
		echo "未找到 buf，正在安装..."; \
		curl -sSL "$$BUF_URL" -o $$BUF_BIN; \
		chmod +x $$BUF_BIN; \
	else \
		INSTALLED=$$($$BUF_BIN --version 2>/dev/null \
			| grep -oE '[0-9]+\.[0-9]+\.[0-9]+' | head -1 || echo ""); \
		EXPECTED="$(BUF_VERSION)"; \
		if [ "$$INSTALLED" != "$$EXPECTED" ] || [ -z "$$INSTALLED" ]; then \
			echo "版本不匹配 (已安装: $$INSTALLED, 需要: $$EXPECTED)"; \
			echo "正在重新安装..."; \
			curl -sSL "$$BUF_URL" -o $$BUF_BIN; \
			chmod +x $$BUF_BIN; \
		else \
			echo "buf $(BUF_VERSION) 已安装"; \
		fi; \
	fi

.PHONY: fmt
fmt: golangci-lint buf
	@$(CURDIR)/bin/golangci-lint fmt
	@$(CURDIR)/bin/buf format -w

.PHONY: lint
lint: fmt golangci-lint
	@$(CURDIR)/bin/golangci-lint run --fix
	@$(CURDIR)/bin/buf lint

.PHONY: gen
gen:
	@buf generate
