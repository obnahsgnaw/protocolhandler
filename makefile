export PACKAGE=github.com/obnahsgnaw/protocolhandler
export INPUT=cmd/main.go
export OUT=out
export APP=protocolhandler
export REPO=xxx.com:5000

.PHONY: help
help:base_help build_help

.PHONY: base_help
base_help:
	@echo "usage: make <option> <params>"
	@echo "options and effects:"
	@echo "    help   : Show help"

include ./build/build/makefile
include ./build/docker/makefile
include ./build/test/makefile
include ./build/version/makefile