# Makefile

APP_NAME := interprerer
APP_DIR := ./app
BUILD_DIR := ./build
BUILD_FILE := $(BUILD_DIR)/$(APP_NAME)

.PHONY: all build debug clean

all: build

build:
	@echo "Building..."
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_FILE) $(APP_DIR)

debug:
	@echo "Building with debug info..."
	@mkdir -p $(BUILD_DIR)
	go build -gcflags "all=-N -l" -o $(BUILD_FILE) $(APP_DIR)

clean:
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)
