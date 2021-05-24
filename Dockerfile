FROM caddy:2.3.0-alpine as caddy
FROM scratch

COPY ./biginit/bin/prod/biginit /
COPY --from=caddy /usr/bin/caddy /usr/bin/caddy

