From elixir:1.12.1-slim AS build
RUN apt-get update
RUN apt-get install -y --no-install-recommends ca-certificates gcc libc6-dev \
  make wget curl npm git curl autoconf automake libtool

RUN mkdir /app
WORKDIR /app

RUN mix local.hex --force && mix local.rebar --force

ARG MIX_ENV
ENV MIX_ENV=$MIX_ENV
ARG DATABASE_URL

COPY mix.exs ./
COPY mix.lock ./

RUN mix deps.get
RUN mix deps.compile

COPY config config

COPY assets assets
COPY priv priv
RUN cd assets && npm rebuild node-sass && npm install && npm run deploy

RUN mix phx.digest

COPY lib lib
RUN mix release
COPY start.sh ./

FROM debian:stable-slim AS app
ARG MIX_ENV
RUN apt-get update
RUN apt-get install -y --no-install-recommends curl bash openssl ca-certificates

RUN mkdir /app
WORKDIR /app

COPY --from=build /app/_build/$MIX_ENV/rel/blyver ./
COPY --from=build /app/start.sh ./

RUN chmod +x ./start.sh
RUN chown -R nobody: /app
USER nobody

ENV HOME=/app

CMD ["./start.sh"]
