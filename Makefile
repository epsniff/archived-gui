SHELL:=/bin/bash 

OS ?= $(shell uname -s |  tr '[:upper:]' '[:lower:]')
DEPVERDEF ?= v0.3.2
CURDEPVER ?= $(shell dep version | sed '2q;d' | sed 's/.*: //')

install-dep:
	# Installs dep tool on release v0.3.2
	go get github.com/golang/dep/cmd/dep
	DEP_BUILD_PLATFORMS=$(OS) && echo $$DEP_BUILD_PLATFORMS
	# Buildiing all architectures...
	cd $(GOPATH)/src/github.com/golang/dep && git checkout $(DEPVERDEF) && ./hack/build-all.bash && cp release/dep-$(OS)-amd64 $(GOPATH)/bin/dep

vendor-ensure:
	# Used to update or delete dependencies. Just don't forget to update `Gopkg.toml` before running.
	#
	[ "$(CURDEPVER)" == "$(DEPVERDEF)" ]
	echo '-- running dep ensure --'
	dep ensure -v
	echo '-- running dep prune to rmeove unused packages --'
	dep prune -v

vendor-rebuild:
	# Used to recreate the vendor directory. 
	#
	[ "$(CURDEPVER)" == "$(DEPVERDEF)" ]
	echo "moving vendor/ to /tmp/vendor-backup-`date +%s`" 
	mv  vendor/ /tmp/vendor-backup-`date +%s` || echo "already gone"
	echo "moving Gopkg.lock to /tmp/Gopkg.lock-`date +%s`"
	mv Gopkg.lock /tmp/Gopkg.lock-`date +%s` || echo "already gone"
	echo '-- running dep ensure --'
	dep ensure -v
	echo '-- running dep prune to rmeove unused packages --'
	dep prune -v

	