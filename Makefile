APP = my_forward_proxy
BIN_DIR = bin

build:
	@mkdir -p $(BIN_DIR)
	go build -o $(BIN_DIR)/$(APP) .

run: build
	./$(BIN_DIR)/$(APP)

config: build
	@./$(BIN_DIR)/$(APP) --init

delete-cache: build
	./$(BIN_DIR)/$(APP) --delete-cache

clean:
	@rm -rf $(BIN_DIR) cache.db  config.toml
	@echo "Cleaned up completed."
