usage:
	@echo "Please specify a task:"
	@awk -F: '/^[^\t#.$$]+:[^=]+?$$/ {print "-",$$1}' Makefile

bin := rich-presence

ifdef WSL_DISTRO_NAME
export GOOS := windows
bin := $(bin).exe
endif

.PHONY: build

build: $(bin)

$(bin): $(shell find . -name '*.go')
	go build -o $(bin) main.go

serve: $(bin)
	./$(bin) serve

clean:
	rm -rf $(bin)
