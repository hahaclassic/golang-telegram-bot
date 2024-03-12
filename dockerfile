FROM golang:1.20-alpine3.18 AS builder

WORKDIR /usr/local/src

RUN apk --no-cache add bash make nano gcc musl-dev  

COPY ["go.mod", "go.sum", "./"]

# RUN go get -d -v ./...  
RUN go mod download

COPY . .
# build
RUN go build -v -o main.exe ./cmd/main.go

FROM alpine AS runner

COPY --from=builder /usr/local/src/main.exe /
COPY /configs/config.yaml /configs/config.yaml

CMD ["/main.exe"]



