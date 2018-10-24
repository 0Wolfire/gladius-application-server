# Gladius Appication Server

# if we are running on a windows machine
# we need to append a .exe to the
BINARY_SUFFIX=
ifeq ($(OS),Windows_NT)
	BINARY_SUFFIX=.exe
endif

ifeq ($(GOOS),windows)
	BINARY_SUFFIX=.exe
endif

# code source and build directories
SRC_DIR=./cmd
DST_DIR=./build

AS_SRC=$(SRC_DIR)/gladius-application-server/main.go
AS_DEST=$(DST_DIR)/gladius-application-server$(BINARY_SUFFIX)

# commands for go
GOBUILD=vgo build
GOTEST=vgo test

# define cleanup target for windows and *nix
ifeq ($(OS),Windows_NT)
clean:
	del /Q /F .\\build\\*
	vgo clean $(AS_SRC)

else
clean:
	rm -rf ./build/*
	vgo clean $(AS_SRC)
endif

# build steps
build:
	make clean
	$(GOBUILD) -o $(AS_DEST) $(AS_SRC)

install:
	make build
	sudo $(AS_DEST) install
	$(AS_DEST) start

update:
	make build
	$(AS_DEST) stop
	$(AS_DEST) uninstall
	sudo $(AS_DEST) install
	$(AS_DEST) start