TARGET=yusifubot

UPX := $(shell command -v upx 2>/dev/null)

all: build compress

build:
	go build -ldflags "-s -w" -o $(TARGET) .

compress: $(TARGET)
ifdef UPX
	$(UPX) -9 $(TARGET)
else
	@echo "UPX not found, skipping compression."
endif

clean:
	rm -f $(TARGET)