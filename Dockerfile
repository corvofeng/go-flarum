FROM alpine:3.7

RUN mkdir /home/goyoubbs
WORKDIR /home/goyoubbs

COPY ./goyoubbs /home/goyoubbs/goyoubbs
COPY ./config/config.yaml /home/goyoubbs/config.yml

EXPOSE 8082
CMD ["/home/goyoubbs/goyoubbs", "-config", "/home/goyoubbs/config.yml"]