-include Makefile.options

data_dir?=data
file?=empty.001

jpg=data/$(file).jpg
meta=tmp/metadata.$(file).json

$(meta): | tmp
	python sign.py --hmac-key $(HMAC_KEY) $(jpg) $@
tmp:
	mkdir -p tmp
    
sign: $(meta) tmp
.phony: sign
################################################################
upload: $(meta)
	curl -i -X POST \
		-H "x-api-key: ${API_KEY}" \
    	-H "x-file-name: $(file).jpg" \
    	-H "x-disallow-duplicates: true" \
    	https://ingestion.edgeimpulse.com/api/training/data \
    	-F attachments[]="@$(meta);type=application/json" \
    	-F attachments[]="@$(jpg);type=image/jpeg"
.phony: upload
################################################################
jpg_files=$(shell ls -1 $(data_dir))
jpg_done_files=$(patsubst %,tmp/.%.done,$(jpg_files))
info:
	echo $(jpg_files)
	echo $(jpg_done_files)
upload-all: $(jpg_done_files)
.phony: upload-all

tmp/.%.jpg.done: tmp/metadata.%.json
	curl -i -X POST \
		-H "x-api-key: ${API_KEY}" \
    	-H "x-file-name: $*.jpg" \
    	-H "x-disallow-duplicates: true" \
    	https://ingestion.edgeimpulse.com/api/training/data \
    	-F attachments[]="@tmp/metadata.$*.json;type=application/json" \
    	-F attachments[]="@$(data_dir)/$*.jpg;type=image/jpeg"
	touch $@	
tmp/metadata.%.json: | tmp
	python sign.py --hmac-key $(HMAC_KEY) $(data_dir)/$*.jpg $@
################################################################
list: $(meta)
	curl -X GET \
		-H "x-api-key: ${API_KEY}" \
    	https://studio.edgeimpulse.com/v1/api/$(PROJECT_ID)/raw-data?category=training
################################################################
remove-all:
	@echo -n "WARNING! REMOVE EVERYTHING?? (remove)/N: " && read ans && [ $${ans:-N} = remove ]
	curl -X POST \
		-H "x-api-key: ${API_KEY}" \
		https://studio.edgeimpulse.com/v1/api/$(PROJECT_ID)/raw-data/delete-all
.phony: remove-all
################################################################
clean:
	rm -rf tmp
