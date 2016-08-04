FROM gliderlabs/alpine
COPY auth-server /usr/bin
COPY swagger.json /
RUN mkdir /data
EXPOSE 3000
ENTRYPOINT ["auth-server"]
