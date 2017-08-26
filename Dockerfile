FROM alpine:3.6

RUN apk update

WORKDIR /opt/services/dash/1.0

COPY ./bin /opt/services/dash/1.0/bin
COPY ./settings.toml /opt/services/dash/1.0/etc/settings.toml

CMD cd /opt/services/dash/1.0/bin && ./dash-bin --config=../etc/settings.toml
