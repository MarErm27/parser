FROM golang

ENV APP_NAME myproject

WORKDIR /go/src/${APP_NAME}
COPY . /go/src/${APP_NAME}

RUN go mod download

RUN go get ./
RUN go build -o ${APP_NAME}

CMD ./${APP_NAME}
