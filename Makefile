
BIN_DIR := ./bin
CMD_DIR := ./cmd
PKG_DIR := ./pkg
INTERNAL_DIR := ./internal

CMD_PATHS := $(wildcard $(CMD_DIR)/*)
BIN_LIST := $(patsubst $(CMD_DIR)/%,$(BIN_DIR)/%.exe,$(CMD_PATHS))

SOURCES = $(wildcard $(PKG_DIR)/**/*.go) $(wildcards $(INTERNAL_DIR)/**/*.go) $(wildcard $(CMD_DIR)/**/*.go)

all: $(BIN_LIST)

$(BIN_DIR)/%.exe: $(CMD_DIR)/% $(SOURCES)
	go build -o $@ ./$<

clean:
	-if [ -d "$(BIN_DIR)" ]; then rm -r $(BIN_DIR); fi
