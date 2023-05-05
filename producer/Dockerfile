FROM node:16 as builder

WORKDIR /usr/src/app

COPY package*.json ./

RUN npm ci

COPY . .

RUN npm run build

FROM node:16-bullseye-slim

ENV NODE_ENV production
USER node

WORKDIR /usr/src/app

COPY package*.json ./

RUN npm ci --production

COPY --from=builder /usr/src/app/dist ./dist

EXPOSE 8000
CMD [ "node", "dist/server.js" ]