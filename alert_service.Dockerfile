FROM node:22.9 AS build

RUN mkdir -p /alert_service
WORKDIR /

COPY alert_service/ /alert_service/
RUN cd /alert_service && npm install

RUN cd /alert_service && npm run build

FROM node:22.9

RUN mkdir -p /alert_service

WORKDIR /

COPY --from=build /alert_service/ /alert_service/
RUN cd /alert_service && npm install --save --omit=dev

CMD ["cd", "/alert_service", "&&", "node", "/.out/main.js"]