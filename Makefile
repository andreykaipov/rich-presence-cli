usage:
	@echo "Please specify a task:"
	@awk -F: '/^[^\t#.$$]+:[^=]+?$$/ {print "-",$$1}' Makefile

bin := rich-presence

ifdef WSL_DISTRO_NAME
export GOOS := windows
bin := $(bin).exe
WSLENV := DISCORD_APP_ID/w
endif

.PHONY: build

build: $(bin)

$(bin): $(shell find . -name '*.go')
	go build -o $(bin) main.go

serve: $(bin)
	./$(bin) serve --verbose

clean:
	rm -rf $(bin)
