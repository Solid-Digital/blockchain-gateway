FROM ubuntu:18.04

ENV TZ=Etc/UTC

# Here we install mongo-tools. It's actually version 3.6, so older than
# the databse, but it's compatible by adding --forceTableScan

RUN apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -y \
    mongo-tools \
    awscli \
    wget \
    ca-certificates \
    curl \
    gnupg2 \
    lsb-release

RUN wget --quiet -O - https://www.postgresql.org/media/keys/ACCC4CF8.asc | apt-key add -
RUN sh -c 'echo "deb http://apt.postgresql.org/pub/repos/apt/ `lsb_release -cs`-pgdg main" >> /etc/apt/sources.list.d/pgdg.list'

RUN apt-get update && apt-get install -y \
    postgresql-client-11
