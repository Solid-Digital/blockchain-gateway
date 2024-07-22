version = $(shell cat VERSION)
date := $(shell date -u '+%Y%m%d%H%M%S')
registry := registry.digitalocean.com
do-app-id-staging := $(shell doctl apps list | grep blyver-staging | awk '{print $1}')

install-tools:
	# install yq
	sudo wget -O /usr/local/bin/yq https://github.com/mikefarah/yq/releases/download/v4.6.3/yq_linux_amd64
	sudo chmod +x /usr/local/bin/yq

write-version:
	echo $(shell if [ ! -z "${TAG}" ]; then echo "${TAG}"; else git describe --tags --always --dirty=-${date}; fi) > VERSION

build-image-%:
	docker build -t ${registry}/unchain/blyver-$*:${version} -f Dockerfile --build-arg MIX_ENV=$* .

push-image-%:
	docker push ${registry}/unchain/blyver-$*:${version}

update-deployment-%:
	doctl apps spec get $(do-app-id-$*) > spec.yml
	yq eval '.services.[] |= (select(.image.repository == "blyver-*") | .image.tag = "$(version)")' -i spec.yml
	doctl apps update $(do-app-id-$*) --spec spec.yml
