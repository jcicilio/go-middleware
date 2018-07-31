#===============================================================================
#
#                    ##        .
#              ## ## ##       ==
#           ## ## ## ##      ===
#       /""""""""""""""""\___/ ===
#  ~~~ {~~ ~~~~ ~~~ ~~~~ ~~ ~ /  ===- ~~~
#       \______ o          __/
#         \    \        __/
#          \____\______/
#
# VERSION:        0.1
# DESCRIPTION:    Dockerfile for go-middleware demonstration
# AUTHOR:         John Cicilio <praeduco@gmail.com>
# COMMENTS:
#   This file describes how to build a container for go-middleware
#
# USAGE:
#   # Build image:
#   docker build --tag="tag name here" .
#
#
#   # Run:
#
#
#===============================================================================

FROM golang:latest
WORKDIR /go/src/go-middleware
RUN mkdir -p /go/src/go-middleware/handlers && mkdir -p /go/src/go-middleware/middleware

COPY go-middleware.go .
COPY handlers/*.go ./handlers/
COPY middleware/*.go ./middleware/

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o go-middleware .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /app/
COPY --from=0 /go/src/go-middleware/go-middleware .
CMD ["./go-middleware"]