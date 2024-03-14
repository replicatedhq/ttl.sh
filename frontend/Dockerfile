FROM nginx:latest

RUN rm /etc/nginx/conf.d/default.conf

COPY nginx.conf /etc/nginx/conf.d/

COPY src/favicon.ico /usr/share/nginx/html/
COPY src/index.htm /usr/share/nginx/html/
COPY src/ttlsh_files /usr/share/nginx/html/ttlsh_files/

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]