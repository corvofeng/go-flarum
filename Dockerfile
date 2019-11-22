FROM alpine:3.7

RUN mkdir /home/goyoubbs
WORKDIR /home/goyoubbs
ENV PROJDIR /home/goyoubbs

COPY ./goyoubbs $PROJDIR/goyoubbs
COPY ./config/config.yaml $PROJDIR/config.yml
COPY ./static $PROJDIR/static
COPY ./view $PROJDIR/view

EXPOSE 8082
CMD ["/home/goyoubbs/goyoubbs", "-config", "/home/goyoubbs/config.yml"]
