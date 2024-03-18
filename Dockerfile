FROM golang:1.22.1

RUN go version
ENV GOPATH=/

COPY ./go.* ./
RUN go mod download

COPY ./ ./

RUN go build -o vk_task ./cmd/main.go

CMD ["./vk_task"]