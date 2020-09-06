# build stage
FROM golang:alpine as builder
RUN mkdir /build
ADD . /build/
WORKDIR /build
RUN go test -c -o main_test

# final stage
FROM scratch
COPY --from=builder /build/main_test /app/
WORKDIR /app
CMD ["./main_test"]