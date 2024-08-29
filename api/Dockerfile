FROM golang:1.23-alpine AS builder

WORKDIR /build

COPY ./api/ api/
COPY ../auth/ ./auth/
COPY ../database/ ./database/
COPY ../model/ ./model/
COPY ../utils/ ./utils/
COPY ../.env .
COPY ../go.mod .
COPY ../go.sum .
ARG PORT=${API_PORT}
#RUN go mod tidy
RUN echo $(ls -1)
RUN echo $(ls -1 api/)
RUN echo $(ls -1 api/run)
RUN go mod download 
RUN go build -C api/run -o reciapi
# RUN go build -C web/ -o recipouir

FROM gcr.io/distroless/base-debian12
WORKDIR /app
COPY --from=builder /build/api/run/reciapi ./reciapi 
#COPY ${LOCAL_CERT} ./certs/
#COPY ${LOCAL_KEY} ./certs/
#ARG cert=${CERT_PATH}
CMD [ "/app/reciapi" ]