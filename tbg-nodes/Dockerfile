FROM rust:1.49.0-slim as rust


FROM elixir:1.11.3-slim AS build

# install build dependencies
RUN apt-get update
RUN apt-get install -y --no-install-recommends ca-certificates gcc libc6-dev make wget curl npm git curl \
  autoconf automake libtool 

ENV RUSTUP_HOME=/usr/local/rustup \
  CARGO_HOME=/usr/local/cargo \
  PATH=/usr/local/cargo/bin:$PATH
COPY --from=rust /usr/local/cargo /usr/local/cargo
COPY --from=rust /usr/local/rustup /usr/local/rustup

# prepare build dir
RUN mkdir /app
WORKDIR /app

# install hex + rebar
RUN mix local.hex --force && \
  mix local.rebar --force

# set build ENV
ARG MIX_ENV
ENV MIX_ENV=$MIX_ENV

# install mix dependencies
COPY mix.exs mix.lock ./

RUN mix deps.get
RUN mix deps.compile

COPY config config

# Move config settings to releases.exs, so they can be loaded at runtime.
RUN mv config/${MIX_ENV}.exs config/releases.exs

# build assets
COPY assets assets
COPY priv priv
RUN cd assets && npm rebuild node-sass && npm install && npm run deploy
RUN mix phx.digest

# build project
COPY lib lib

# build release (uncomment COPY if rel/ exists)
COPY rel rel
# Copy VERSION as late as possible in the build process since it busts the build cache.
COPY VERSION ./

RUN mix release

# prepare release image
FROM debian:stable-slim AS app
ARG MIX_ENV
RUN apt-get update
RUN apt-get install -y --no-install-recommends curl bash openssl  ca-certificates 

RUN mkdir /app
RUN mkdir /app/kubeconfig
WORKDIR /app

COPY --from=build /app/_build/$MIX_ENV/rel/tbg_nodes ./
RUN chown -R nobody: /app
USER nobody

ENV HOME=/app
