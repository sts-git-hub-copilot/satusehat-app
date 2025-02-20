FROM alpine:3.13.5

RUN /bin/sh -c "apk update && \
    apk add --no-cache tzdata bash libc6-compat postgresql-client"

RUN mkdir app
ADD ssbackend app/app
ADD asset /app/asset

WORKDIR /app
CMD [ "/bin/sh","-c","cd /app && ./app"]

