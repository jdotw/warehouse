FROM golang AS build

WORKDIR /app

COPY . . 

ARG GITHUB_USER
ARG GITHUB_PAT
RUN echo "machine github.com login ${GITHUB_USER} password ${GITHUB_PAT}" > ~/.netrc
RUN chmod 0600 ~/.netrc 
RUN go mod tidy -compat=1.17
RUN go mod download 
RUN CGO_ENABLED=0 go build -o app .

##
## Deploy
##
FROM alpine:latest  

WORKDIR /root/

COPY --from=build /app/app ./

EXPOSE 8081
EXPOSE 8082

CMD ["./app"]
