# FROM golang AS builder

# ENV ARES=$GOPATH/src/bitbucket.org/ares
# COPY . $ARES
# WORKDIR $ARES
# RUN echo $ARES; pwd; ls; 
# RUN scripts/build.sh
# # RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-X 'main.metaTOML=${metaTOML}'" -a -installsuffix cgo -o ${OUTPUT} $1



FROM unchain/alpine-glibc-dockerize
RUN apk update && apk add postgresql-client

COPY ./bin /app/
RUN ln -s /app/ares* /sbin/
COPY migrations migrations
COPY fixtures fixtures
COPY scripts/database /usr/local/bin/

WORKDIR /app/

ENTRYPOINT /app/ares
EXPOSE 8080
