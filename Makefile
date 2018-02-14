SHELL:=/bin/bash 

OS ?= $(shell uname -s |  tr '[:upper:]' '[:lower:]')
DEPVER ?= v0.4.1
CURDEPVER ?= $(shell dep version | sed '2q;d' | sed 's/.*: //')

install-dep:
	# Installs dep tool on release
	go get -u github.com/golang/dep/cmd/dep
	cd $(GOPATH)/src/github.com/golang/dep && \
		git checkout $(DEPVER) && \
		export DEP_BUILD_PLATFORMS=$(OS) DEP_BUILD_ARCHS=$(ARCH) && \
		./hack/build-all.bash && \
		cp release/dep-$(OS)-$(ARCH) $(GOPATH)/bin/dep

vendor-ensure:
	# Used to update or delete dependencies. Just don't forget to update `Gopkg.toml` before running.
	#
	[ "$(CURDEPVER)" == "$(DEPVERDEF)" ]
	echo '-- running dep ensure --'
	dep ensure -v
	echo '-- running dep prune to rmeove unused packages --'
	dep prune -v



