FROM node:alpine

WORKDIR /app/
COPY package.json /app/package.json
RUN yarn install

COPY . /app/

ENV HOST=0.0.0.0
ENV PORT=3000

EXPOSE 3000
CMD node app.js