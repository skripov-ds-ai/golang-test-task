FROM nginx:latest

RUN rm /etc/nginx/nginx.conf /etc/nginx/conf.d/default.conf
ADD nginx.conf /etc/nginx/nginx.conf
COPY conf.conf /etc/nginx/conf.d