CIRCLE_BUILD_NUM ?= 0

APP=straw-pod
PKG=/go/src/github.com/TimWoolford/${APP}
TAG=timwoolford/${APP}:0.1.$(CIRCLE_BUILD_NUM)

BIN=$(firstword $(subst :, ,${GOPATH}))/bin
GODEP = $(BIN)/dep
M = $(shell printf "\033[34;1m▶\033[0m")

.PHONY: gobuild
gobuild: vendor ; $(info $(M) building…)
	GOOS=linux go build -v -o bin/${APP} ./cmd/${APP}

.PHONY: gotest
gotest: gobuild ; $(info $(M) running tests…)
	@go test ./...

.PHONY: build
build:
	docker run --rm \
	 -v "${PWD}":${PKG} \
	 -w ${PKG} \
	 golang:1.9 \
	 make gobuild

.PHONY: build-image
build-image:
	docker build -t ${TAG} .

.PHONY: push-image
push-image:
	docker push ${TAG}

.PHONY: clean
clean: ; $(info $(M) cleaning…)
	@docker images -q ${APP} | xargs docker rmi -f
	@rm -rf bin/*

.PHONY: vendor
vendor: .vendor

.vendor: Gopkg.toml Gopkg.lock
	command -v dep >/dev/null 2>&1 || go get github.com/golang/dep/cmd/dep
	$(GODEP) ensure -v
	@touch $@

clean-minikube:
	helm delete ${APP} --purge

.PHONY: deploy-minikube
deploy-minikube:
	helm upgrade --install ${APP} charts/minikube --namespace test