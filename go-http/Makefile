version = $(shell cat VERSION)
date := $(shell date -u '+%Y%m%d%H%M%S')
registry := registry.digitalocean.com
do-app-id-prod := $(shell doctl apps list | grep go-http-prod | awk '{print $1}')
do-app-id-staging := $(shell doctl apps list | grep go-http-staging | awk '{print $1}')
do-app-id-latest := $(shell doctl apps list | grep go-http-latest | awk '{print $1}')

build:
	CGO_ENABLE=0 go build -o bin/go-http cmd/go-http/...

install-tools:
	# install yq
	sudo wget -O /usr/local/bin/yq https://github.com/mikefarah/yq/releases/download/v4.6.3/yq_linux_amd64
	sudo chmod +x /usr/local/bin/yq

image:
	docker build -t ${registry}/unchain/go-http:${version} .

push:
	docker push ${registry}/unchain/go-http:${version}

update-deployment-%:
	# get deployed app spec
	doctl apps spec get $(do-app-id-$*) > spec.yml
	# update image tag
	yq eval '.services.[0].image.tag = "$(version)"' -i spec.yml
	# update DO app deployment
	doctl apps update $(do-app-id-$*) --spec spec.yml

echo-version:
	echo ${version}

write-version:
	# Setting export TAG={value} will overwrite the version.
	echo $(shell if [ ! -z "${TAG}" ]; then echo "${TAG}"; else git describe --tags --always --dirty=-${date}; fi) > VERSION