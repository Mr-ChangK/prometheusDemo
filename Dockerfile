FROM golang

WORKDIR /home

COPY * /home

RUN go build -o demo

CMD ./demo