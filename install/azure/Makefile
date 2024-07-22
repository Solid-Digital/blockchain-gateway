
define mklockdir
		mkdir -p .lock
		touch .lock/$1
endef

define machine_name
$(shell cat .lock/machine_name)
endef

define machine_ip
$(shell cat .lock/machine_ip)
endef

define hostname
$(shell cat .lock/hostname)
endef


tt:
	echo $(machine_name)


.lock:
	mkdir -p .lock

##
# Intiialize docker-machine
##
.lock/machine_name: .lock
	@read -p "name your machine:" machine_name; \
	echo $$machine_name > .lock/machine_name

docker-machine: .lock/machine_name
	docker-machine create \
	--driver amazonec2 \
	--amazonec2-instance-type "t2.small" \
	--amazonec2-region "eu-central-1" \
	--amazonec2-open-port 80 \
	--amazonec2-open-port 443 \
	--amazonec2-open-port 8200 \
	--amazonec2-open-port 8585 \
	--amazonec2-open-port 8787 \
	$(machine_name)

	@echo OK. Please run: docker-machine env $(machine_name)

	docker-machine ip $(machine_name) > .lock/machine_ip
	$(call mklockdir,docker-machine)

azure-docker-machine: .lock/machine_name
	docker-machine create \
	--driver azure \
	--azure-location "northeurope" \
	--azure-size "Standard_A2" \
	--azure-ssh-user "docker-user" \
  --azure-image "canonical:UbuntuServer:16.04.0-LTS:latest" \
	--azure-subscription-id "ca1abd6a-e1cf-492b-9340-a288e1abaa06" \
	$(machine_name)

	@echo OK. Please run: docker-machine env $(machine_name)

	docker-machine ip $(machine_name) > .lock/machine_ip
	$(call mklockdir,docker-machine)

###
# Install remote
###

.lock/make:
	docker-machine ssh $(machine_name) sudo apt-get install -y make
	$(call mklockdir,make)
install-make: .lock/make

upload-makefile: install-make
	docker-machine scp Makefile $(machine_name):Makefile

do-make:
	docker-machine ssh $(machine_name) sudo make swarm


.lock/machine_ip:
	@echo "No docker-machine ip is set. Is it configured?"
	exit 1

##
# Intiialize the swarm host(s)
##
.lock/swarm: export DOCKER_ORCHESTRATOR=swarm
# .lock/swarm: .lock/machine_ip
.lock/swarm:
	docker swarm init # --advertise-addr $(machine_ip)
	$(call mklockdir,swarm)
swarm: .lock/swarm


##
# Preconditions for traefik
##
.lock/network:
	docker network create -d overlay proxy
	$(call mklockdir,network)

.lock/volume:
	docker volume create traefik-acme-data
	$(call mklockdir,volume)

.lock/hostname:
	@read -p "To which domain should the adapter(s) listen? :" hostname; \
	echo $$hostname > .lock/hostname

##
## Traefik, includes letsencrypt
##
.lock/traefik: export DOCKER_ORCHESTRATOR = swarm
.lock/traefik: .lock/network .lock/volume .lock/hostname
	docker stack deploy -c docker-compose.traefik.yml traefik
traefik: .lock/traefik
	$(call mklockdir,traefik)

destroy-treafik: export DOCKER_ORCHESTRATOR = swarm
destroy-traefik:
	docker stack rm traefik

##
## Traefik, includes letsencrypt
##
.lock/adapters: export DOCKER_ORCHESTRATOR = swarm
.lock/adapters:
	HOSTNAME=$(hostname) docker stack deploy -c docker-compose.adapters.yml adapters
adapters: .lock/adapters


tls:
	mkdir -p tls
	./setup-cert.sh

all: .lock/network .lock/traefik .lock/adapters
