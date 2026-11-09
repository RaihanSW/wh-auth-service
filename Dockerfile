##Stage 1
## Start from the latest golang base image for builder
FROM golang:latest as builder
# ENV GO111MODULE=on
RUN mkdir /go/src/service
ADD . /go/src/service
WORKDIR /go/src/service
ENV GOPRIVATE=github.com/cinemaxxi-team
#RUN apk add git
# Copy the .netrc file, NOTE the file from Jenkins
RUN cp /go/src/service/.netrc $HOME/.netrc
RUN go mod tidy
RUN go mod vendor

#INSTALL SWAGGER
#VERSI DISAMAKAN DENGAN ISI FILE go.mod
RUN go install -v github.com/swaggo/swag/cmd/swag@v1.16.4
RUN swag init

# BUILD ENGINE
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-w -s" -o application main.go
RUN chmod 755 /go/src/service/application

##Stage 2
## Start from the latest alpine base image
FROM alpine:latest
LABEL maintainer="Muhammad Nasrul <mnasruul@gmail.com>"
RUN addgroup -S xxiadmin && adduser -S xxiadmin -G xxiadmin
RUN apk update && apk upgrade && apk add --no-cache tzdata && apk add --no-cache curl
ENV TZ=Asia/Jakarta
ARG APP_ENV
ENV env_state=$APP_ENV

RUN mkdir -p /app
RUN chown xxiadmin:xxiadmin /app
WORKDIR /app
# When using COPY with more than one source file, the destination must be a directory and end with a /
COPY --chown=xxiadmin:xxiadmin --from=builder /go/src/service/*.json /app/
COPY --chown=xxiadmin:xxiadmin --from=builder /go/src/service/*.html /app/
COPY --chown=xxiadmin:xxiadmin --from=builder  /go/src/service/application /app 
USER xxiadmin

# Expose port 2000 for application and 9090 for profiler to the outside world.
EXPOSE 4000
EXPOSE 9090

CMD /app/application -env=$env_state
