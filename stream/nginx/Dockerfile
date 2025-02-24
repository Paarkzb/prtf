FROM alpine:latest AS builder

ENV NGINX_VERSION nginx-1.27.4
ENV NGINX_HTTP_FLV_MODULE_VERSION 1.2.12

# Установка зависимостей для сборки RTMP модуля
RUN apk update &&  apk add --no-cache --virtual .build-deps gcc libc-dev make zlib-dev openssl-dev pcre-dev

# Сборка nginx с RTMP модулем
RUN cd /tmp && \
    wget -O nginx-http-flv-module-${NGINX_HTTP_FLV_MODULE_VERSION}.tar.gz https://github.com/winshining/nginx-http-flv-module/archive/v${NGINX_HTTP_FLV_MODULE_VERSION}.tar.gz && \
    tar xzf nginx-http-flv-module-${NGINX_HTTP_FLV_MODULE_VERSION}.tar.gz && \
    wget -O ${NGINX_VERSION}.tar.gz http://nginx.org/download/${NGINX_VERSION}.tar.gz && \
    tar xzf ${NGINX_VERSION}.tar.gz

RUN cd /tmp/${NGINX_VERSION} && \
    ./configure \
        --with-compat \
        --sbin-path=/usr/sbin/nginx \
        --conf-path=/etc/nginx/nginx.conf \
        --http-log-path=/var/log/nginx/access.log \
        --error-log-path=/var/log/nginx/error.log \
        --with-pcre \
        --lock-path=/var/lock/nginx.lock \
        --pid-path=/var/run/nginx.pid \
        --with-http_ssl_module \
        --with-http_v2_module \
        --with-http_stub_status_module \
        --add-dynamic-module=/tmp/nginx-http-flv-module-${NGINX_HTTP_FLV_MODULE_VERSION} \
        --with-debug && \
    make && \
    make install

# Финальный образ
FROM nginx:1.27.4-alpine

RUN apk update && apk add ffmpeg

COPY --from=builder /usr/local/nginx/modules/ngx_http_flv_live_module.so /etc/nginx/modules/

RUN mkdir -p /var/www/hls /var/rec /var/vod && \
    chown -R nginx:nginx /var/www /var/rec /var/vod
    
COPY nginx.conf /etc/nginx/nginx.conf
COPY generate_master_playlist.sh /var/rec/

