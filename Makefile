usage:
	@echo "Please specify a task:"
	@awk -F: '/^[^\t#.$$]+:[^=]+?$$/ {print "-",$$1}' Makefile

export DISCORD_APP_ID := 942604338927374438

bin := rich-presence

ifdef WSL_DISTRO_NAME
export GOOS := windows
bin := $(bin).exe
WSLENV := DISCORD_APP_ID/w
endif

.PHONY: build

build: $(bin)

$(bin): main.go
	go build -o $(bin) $<

run: $(bin)
	./$(bin) -verbose

clean:
	rm -rf $(bin)
