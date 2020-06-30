FROM golang:1.14.2-alpine as build

RUN mkdir /app
ADD . /app/
WORKDIR /app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o safira .

FROM docker:19.03.8-git

RUN apk update && \
    apk add --no-cache libc6-compat

COPY --from=build /app/safira /usr/local/bin/safira
CMD ["safira"]
