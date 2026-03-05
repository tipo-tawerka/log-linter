BINARY_NAME  := log-linter
BUILD_DIR    := bin
CMD_PATH     := ./cmd/log-linter

# Пакеты для покрытия — исключаем точки входа
TEST_PACKAGES := $(shell go list ./... | grep -v '/cmd/')

.PHONY: all build test coverage clean

all: build

## build: собрать бинарник в ./bin/
build:
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) $(CMD_PATH)
	@echo "Built: $(BUILD_DIR)/$(BINARY_NAME)"

## test: запустить unit-тесты
test:
	go test $(TEST_PACKAGES)

## coverage: запустить тесты и показать процент покрытия
coverage:
	go test -coverprofile=coverage.out $(TEST_PACKAGES)
	@echo ""
	@echo "=== Coverage ==="
	go tool cover -func=coverage.out

## clean: удалить артефакты сборки
clean:
	@rm -rf $(BUILD_DIR) coverage.out
	@echo "Cleaned"

## lint-examples: прогнать линтер по тестовым примерам в testdata/ и показать нарушения
lint-examples: build
	@echo ""
	@echo "package: all_correct"
	@cd testdata && ../$(BUILD_DIR)/$(BINARY_NAME) ./all_correct/ 2>&1 && echo "  OK" || true
	@echo ""
	@echo "package: uppercase_message"
	@cd testdata && ../$(BUILD_DIR)/$(BINARY_NAME) ./uppercase_message/ 2>&1; true
	@echo ""
	@echo "package: russian_message"
	@cd testdata && ../$(BUILD_DIR)/$(BINARY_NAME) ./russian_message/ 2>&1; true
	@echo ""
	@echo "package: specchar_message"
	@cd testdata && ../$(BUILD_DIR)/$(BINARY_NAME) ./specchar_message/ 2>&1; true
	@echo ""
	@echo "package: sensitive_args"
	@cd testdata && ../$(BUILD_DIR)/$(BINARY_NAME) ./sensitive_args/ 2>&1; true
	@echo ""
	@echo "package: skip_tests"
	@cd testdata && ../$(BUILD_DIR)/$(BINARY_NAME) ./skip_tests/ 2>&1; true
	@echo ""
	@echo " с -skiptests: тест-файлы не проверяются, нарушений нет"
	@cd testdata && ../$(BUILD_DIR)/$(BINARY_NAME) -skiptests ./skip_tests/ 2>&1 && echo "  OK" || true
	@echo ""
	@echo "package: mixed_violations"
	@cd testdata && ../$(BUILD_DIR)/$(BINARY_NAME) ./mixed_violations/ 2>&1; true
	@echo ""
