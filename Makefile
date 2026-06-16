APP=Tunneler
EXECUTABLE=tunneler
IDENTIFIER=app.tunneler
APP_DIR=$(APP).app
BINARY=$(APP_DIR)/Contents/MacOS/$(EXECUTABLE)
PLIST=$(APP_DIR)/Contents/Info.plist

GO_SOURCES := $(shell find . -name '*.go' -not -path './third_party/menuet/cmd/*')

.PHONY: all run install clean

all: $(BINARY) $(PLIST)

run: all
	open $(APP_DIR)

install: all
	cp -R $(APP_DIR) /Applications/

$(BINARY): $(GO_SOURCES) go.mod go.sum
	mkdir -p $(dir $@)
	go build -o $(BINARY) ./cmd/tunneler

$(PLIST): Makefile
	mkdir -p $(dir $@)
	@echo '<?xml version="1.0" encoding="UTF-8"?>' > $(PLIST)
	@echo '<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">' >> $(PLIST)
	@echo '<plist version="1.0">' >> $(PLIST)
	@echo '<dict>' >> $(PLIST)
	@echo '  <key>CFBundleExecutable</key>' >> $(PLIST)
	@echo '  <string>$(EXECUTABLE)</string>' >> $(PLIST)
	@echo '  <key>CFBundleIdentifier</key>' >> $(PLIST)
	@echo '  <string>$(IDENTIFIER)</string>' >> $(PLIST)
	@echo '  <key>CFBundleName</key>' >> $(PLIST)
	@echo '  <string>$(APP)</string>' >> $(PLIST)
	@echo '  <key>CFBundlePackageType</key>' >> $(PLIST)
	@echo '  <string>APPL</string>' >> $(PLIST)
	@echo '  <key>CFBundleShortVersionString</key>' >> $(PLIST)
	@echo '  <string>0.1</string>' >> $(PLIST)
	@echo '  <key>LSUIElement</key><true/>' >> $(PLIST)
	@echo '  <key>NSHighResolutionCapable</key><true/>' >> $(PLIST)
	@echo '</dict>' >> $(PLIST)
	@echo '</plist>' >> $(PLIST)

clean:
	rm -rf $(APP_DIR)
