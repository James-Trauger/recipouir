FROM golang:1.23-alpine AS builder

#WORKDIR /build
WORKDIR /src
#COPY ./api/ api/
#COPY ../auth/ ./auth/
#COPY ../database/ ./database/
#COPY ../model/ ./model/
#COPY ../utils/ ./utils/
#COPY ../.env .
#COPY ../go.mod .
#COPY ../go.sum .
COPY ./ /src/
ARG PORT=${API_PORT}
#RUN go mod tidy
#RUN echo $(ls -1)

RUN go mod download 
RUN go build -C api/run -o /go/bin/reciapi
#RUN go build -C web/ -o /go/bin/recipouir

FROM gcr.io/distroless/base-debian12
WORKDIR /app
COPY --from=builder /go/bin/reciapi ./reciapi 
COPY --from=builder /src/.env .

##COPY ${LOCAL_CERT} ./certs/
##COPY ${LOCAL_KEY} ./certs/
##ARG cert=${CERT_PATH}

CMD [ "/app/reciapi" ]