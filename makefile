.PHONY: build install clean

BINARY_NAME := go-typer
BUILD_DIR := ./bin
BUILD_PATH := $(BUILD_DIR)/$(BINARY_NAME)

build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_PATH)

install: build
	@echo "Starting installation process..."
	@if [ ! -f "$(BUILD_PATH)" ]; then \
		echo "❌ Error: Binary not found at $(BUILD_PATH)"; \
		exit 1; \
	fi
	@echo "Please select your preferred installation method:"
	@echo "1) Local user install (~/.local/bin) - Recommended for single user"
	@echo "2) System-wide install (/usr/local/bin) - Requires sudo, for all users"
	@read -p "Enter your choice [1-2]: " install_choice; \
	if [[ ! "$$install_choice" =~ ^[12]$$ ]]; then \
		echo "❌ Invalid selection. Please choose 1 or 2."; \
		exit 1; \
	fi; \
	case "$$(uname -s)" in \
		Linux*) OS=Linux ;; \
		Darwin*) OS=macOS ;; \
		*) OS="UNKNOWN" ;; \
	esac; \
	if [[ "$$OS" == "UNKNOWN" ]]; then \
		echo "⚠️ Unsupported operating system detected. Proceeding with basic installation."; \
	fi; \
	case $$install_choice in \
		1) \
			TARGET_DIR="$$HOME/.local/bin"; \
			mkdir -p "$$TARGET_DIR"; \
			if cp "$(BUILD_PATH)" "$$TARGET_DIR/"; then \
				echo "✅ Successfully installed to $$TARGET_DIR/$(BINARY_NAME)"; \
			else \
				echo "❌ Local installation failed. Please check permissions."; \
				exit 1; \
			fi; \
			;; \
		2) \
			TARGET_DIR="/usr/local/bin"; \
			echo "Installing system-wide (requires sudo privileges)..."; \
			if sudo cp "$(BUILD_PATH)" "$$TARGET_DIR/"; then \
				echo "✅ Successfully installed to $$TARGET_DIR/$(BINARY_NAME)"; \
			else \
				echo "❌ System-wide installation failed. Please check sudo permissions."; \
				exit 1; \
			fi; \
			;; \
	esac; \

# ---  Clean Target ---
clean:
	@echo "Cleaning up build artifacts..."
	@rm -rf $(BUILD_DIR)
	@echo "✅ Cleanup complete."
# --------------------------

.DEFAULT_GOAL := install
