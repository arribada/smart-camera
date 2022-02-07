-include Makefile.options
######################################################################################
BIN_DIR=$(PWD)/bin
BIN_NAME=smart-camera
######################################################################################
$(BIN_DIR):
	mkdir -p $(BIN_DIR)
$(BIN_DIR)/rpi:
	mkdir -p $(BIN_DIR)	
rpi/build: | $(BIN_DIR)/rpi
	(cd cmd/cli && CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=5 go build -o $(BIN_DIR)/rpi/$(BIN_NAME))

build: | $(BIN_DIR)
	(cd cmd/cli && CGO_ENABLED=0 go build -o $(BIN_DIR)/$(BIN_NAME))
test:
	go test -v -race -count=1 ./...	
######################################################################################
clean:
	rm -r -f $(BIN_DIR)
######################################################################################
rpi/upload: $(BIN_DIR)/rpi/$(BIN_NAME)
	rsync -P $(BIN_DIR)/rpi/$(BIN_NAME) $(RPI_USER_URL):$(RPI_DIR)

rpi/ssh:
	ssh $(RPI_USER_URL)

rpi/dwn-images: 
	rsync -P -r $(RPI_USER_URL):$(RPI_DIR)/data ./data
######################################################################################
.PHONY: rpi/build rpi/upload rpi/ss rpi/dwn-mages build clean
