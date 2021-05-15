FROM golang:alpine
COPY . /opt/s4g
WORKDIR /opt/s4g
EXPOSE 10001
ENTRYPOINT ["go","run","/opt/s4g/main.go"]