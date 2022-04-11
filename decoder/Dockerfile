FROM node:17-alpine3.14

COPY package*.json .

RUN yarn install

COPY . .

CMD [ "yarn", "start" ]
