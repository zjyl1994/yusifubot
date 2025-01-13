TARGET=yusifubot

UPX := $(shell command -v upx 2>/dev/null)

all: build compress

build:
	go build -o $(TARGET) .

compress: $(TARGET)
ifdef UPX
	$(UPX) $(TARGET)
else
	@echo "UPX not found, skipping compression."
endif

clean:
	rm -f $(TARGET)