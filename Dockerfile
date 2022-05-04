FROM golang:1.18-alpine
RUN sudo apt-get install libaom-dev
WORKDIR /app
COPY . .
RUN go mod tidy