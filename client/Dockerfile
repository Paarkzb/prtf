# syntax=docker/dockerfile:1

# build stage
FROM node:latest as build-stage

WORKDIR /app

COPY package*.json ./

RUN npm install

COPY ./ ./

RUN npm run build

# production stage

FROM nginx:1.26-alpine as production-stage

WORKDIR /app
# RUN mkdir /app

COPY --from=build-stage /app/dist /prtf

EXPOSE 8080

COPY nginx.conf /etc/nginx/nginx.conf

CMD ["nginx", "-g", "daemon off;"]