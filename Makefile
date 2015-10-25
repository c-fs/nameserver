ORG_PATH:=github.com/c-fs
REPO_PATH:=$(ORG_PATH)/nameserver

.PHONY: docker
docker:
	# Static binary built here may fail to call os/user/lookup functions due to library
	# conflict. (http://stackoverflow.com/questions/8140439/why-would-it-be-impossible-to-fully-statically-link-an-application)
	# Because nameserver doesn't use these functions, it is ok to ignore the error.
	go build -a -tags netgo -installsuffix netgo --ldflags '-extldflags "-static"' -o nameserver ${REPO_PATH}/server
	docker build -t yunxing/nameserver .
	rm nameserver
