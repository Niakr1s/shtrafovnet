FROM alpine:latest
COPY server server
CMD [ "./server" ]
EXPOSE 8081 9000