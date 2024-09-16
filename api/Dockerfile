FROM golang:1.23-alpine AS builder

WORKDIR /build 

COPY . /build
#ARG PORT=${API_PORT}
#RUN go mod tidy
RUN echo $(ls -1)
RUN echo $(pwd)
RUN echo $(ls -1 /certs)
RUN echo $(ls -1 /)

RUN go mod download -C /build
RUN go build -C api/run -o /go/bin/reciapi 
#RUN go build -C web/ -o /go/bin/recipouir

FROM gcr.io/distroless/base-debian12
WORKDIR /app
COPY --from=builder /go/bin/reciapi ./reciapi 
COPY --from=builder /build/.env ./
#COPY --from=builder /certs /certs

CMD [ "/app/reciapi" ]