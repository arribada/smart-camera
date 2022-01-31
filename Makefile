-include Makefile.options
######################################################################################
BIN_DIR=$(PWD)/bin
BIN_NAME=smart-camera
######################################################################################
$(BIN_DIR):
	mkdir $(BIN_DIR)
build: | $(BIN_DIR)
	(cd cmd/cli && CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=5 go build -o $(BIN_DIR)/$(BIN_NAME))
######################################################################################
clean:
	rm -r -f $(BIN_DIR)

######################################################################################
rpi/upload: $(BIN_DIR)/$(BIN_NAME)
	rsync -P $(BIN_DIR)/$(BIN_NAME) $(RPI_USER_URL):$(RPI_DIR)

rpi/ssh:
	ssh $(RPI_USER_URL)

rpi/dwn-images: 
	rsync -P -r $(RPI_USER_URL):$(RPI_DIR)/data ./data
######################################################################################