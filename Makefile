SHA := $(shell git rev-parse --short HEAD)
VERSION := $(shell cat VERSION)
ITTERATION := $(shell date +%s)

all: build

build:
	mkdir -p packaging/output
	mkdir -p packaging/worker/root/usr/local/bin
	mkdir -p packaging/api/root/usr/local/bin
	go build -o packaging/worker/root/usr/local/bin/utisak-worker \
		-ldflags "-X main.revision $(SHA) -X main.version $(VERSION)" \
		github.com/rastasheep/utisak-worker/main/worker
	go build -o packaging/api/root/usr/local/bin/utisak-api \
		-ldflags "-X main.revision $(SHA) -X main.version $(VERSION)" \
		github.com/rastasheep/utisak-worker/main/api

test:
	@test -z "$(shell find . -name '*.go' | xargs gofmt -l)" || (echo "Need to run 'go fmt ./...'"; exit 1)
	go vet ./...
	go test -cover -short ./...

run-api:
	@go run main/api/main.go --config=main/api/config/config.json

run-worker:
	@go run main/worker/main.go --config=main/worker/config/config.json

clean:
	find . -name "*.out" -delete
	rm -rf packaging/output
	rm -f packaging/api/root/usr/local/bin/utisak-api || true
	rm -r packaging/worker/root/usr/local/bin/utisak-worker || true

packages: clean build deb rpm

deb:
	$(call build_dep,api,utisak-api)
	$(call build_dep,worker,utisak-worker)

rpm:
	$(call build_rpm,api,utisak-api)
	$(call build_rpm,worker,utisak-worker)

define build_dep
	fpm -s dir -t deb \
		--name $2 \
		--version $(VERSION) \
		--package packaging/output/$2.deb \
		--deb-priority optional \
		--category admin \
		--force \
		--iteration $(ITTERATION) \
		--deb-compression bzip2 \
		--after-install packaging/$1/scripts/postinst.deb \
		--before-remove packaging/$1/scripts/prerm.deb \
		--after-remove packaging/$1/scripts/postrm.deb \
		--url https://github.com/rastasheep/utisak-worker \
		--description "Discover news, easily" \
		--maintainer "Aleksandar Diklic <rastasheep@gmail.com>" \
		--license "MIT" \
		packaging/$1/root/=/
	cp packaging/output/$2.deb packaging/output/$2.deb.$(SHA)
endef

define build_rpm
	fpm -s dir -t rpm \
		--name $2 \
		--version $(VERSION) \
		--package packaging/output/$2.rpm \
		--rpm-compression bzip2 \
		--rpm-os linux \
		--force \
		--iteration $(ITTERATION) \
		--after-install packaging/$1/scripts/postinst.rpm \
		--before-remove packaging/$1/scripts/prerm.rpm \
		--after-remove packaging/$1/scripts/postrm.rpm \
		--url https://github.com/rastasheep/utisak-worker \
		--description "Discover news, easily" \
		--maintainer "Aleksandar Diklic <rastasheep@gmail.com>" \
		--license "MIT" \
		packaging/$1/root/=/
	cp packaging/output/$2.rpm packaging/output/$2.rpm.$(SHA)
endef
