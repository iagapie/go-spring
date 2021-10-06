.DEFAULT_GOAL := help

.PHONY: help
help:
	@awk 'BEGIN {FS = ":.*?## "}; /^[a-zA-Z-_]+:.*?## .*$$/ {printf "\033[32m%-15s\033[0m %s\n", $$1, $$2}' Makefile | sort

.PHONY: build_plugin_demo
build_plugin_demo: ## Build demo plugin
	go build -buildmode=plugin -o plugins/spring/demo/spring_demo.so plugins/spring/demo/spring_demo.go

.PHONY: build
build: ## Build app
	go build

.PHONY: run
run: build_plugin_demo ## Run app
	go run main.go || true
