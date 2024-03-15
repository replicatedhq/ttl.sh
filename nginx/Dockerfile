FROM nginx:latest

COPY nginx.conf /etc/nginx/nginx.conf

COPY src/favicon.ico /usr/share/nginx/html/
COPY src/index.htm /usr/share/nginx/html/
COPY src/ttlsh_files /usr/share/nginx/html/ttlsh_files/

EXPOSE 8080
