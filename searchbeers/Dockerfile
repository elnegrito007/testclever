FROM golang:1.16-alpine
WORKDIR /app
COPY . /app
RUN go mod tidy
EXPOSE 21002
RUN go build -o /main
CMD [ "/main" ]



