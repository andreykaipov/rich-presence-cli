usage:
	@echo "Please specify a task:"
	@awk -F: '/^[^\t#.$$]+:[^=]+?$$/ {print "-",$$1}' Makefile

name := rich-presence
bin := ./bin/$(name)
version := $(shell grep -Eo 'version = "[^"]+"' main.go | grep -Eo '[0-9]+\.[0-9]+\.[0-9]+')

ifdef WSL_DISTRO_NAME
export GOOS := windows
bin := $(bin).exe
endif

.PHONY: build

build: $(bin)

$(bin): $(shell find . -name '*.go')
	go build -o $(bin) main.go

serve: $(bin)
	$(bin) serve

clean:
	rm -rf bin

release: clean
	@for os in linux windows; do\
		if [ $$os = "windows" ]; then ext=".exe"; fi;\
		out="bin/release/$(name)-$$os-amd64$$ext";\
		echo "Building $$out";\
		GOOS=$$os GOARCH=amd64 CGO_ENABLED=0 go build -ldflags='-w -s' -o $$out main.go;\
	done
	upx -5 bin/release/*
	gh release create v$(version) bin/release/* --generate-notes
