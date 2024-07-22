!#/bin/bash

## Install docker
curl -L https://get.docker.com/ |sh
sudo usermod -aG docker ghuser

## initialize swarm
sudo docker swarm init

## install traefik
wget https://raw.githubusercontent.com/unchainio/install/master/azure/docker-compose.traefik.yml
wget https://raw.githubusercontent.com/unchainio/install/master/azure/traefik.toml

sudo docker network create -d overlay proxy
sudo docker volume create traefik-acme-data
sudo docker stack deploy -c docker-compose.traefik.yml traefik

## install initial adapters
wget https://raw.githubusercontent.com/unchainio/install/master/azure/docker-compose.adapters.yml
sudo docker stack deploy -c docker-compose.adapters.yml adapters

echo "all done, you should now be able to access the adapter under <host>:80/adapter"
