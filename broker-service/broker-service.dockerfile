FROM alpine:3.14

RUN mkdir /app

COPY brokerApp /app

CMD [ "/app/brokerApp"]