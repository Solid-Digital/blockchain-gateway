FROM ubuntu:18.04

ENV TZ=Etc/UTC
# ENV TZ=Europe/Amsterdam

# Here we install mongo-tools. It's actually version 3.6, so older than
# the databse, but it's compatible by adding --forceTableScan

RUN apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -y \
    mongo-tools \
    awscli 

