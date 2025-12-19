FROM cgr.dev/chainguard/nginx:latest@sha256:dad6ecc27985d8f09292bb6df5778caa001299af8f486fb023b7efda3d3f3a10

COPY nginx.conf /etc/nginx/nginx.conf
COPY html/ /usr/share/nginx/html/
