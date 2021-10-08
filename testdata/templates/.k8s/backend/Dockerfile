FROM git.test.com:5005/devops/docker/php:8.0.10

ENV COMPOSER_CACHE_DIR=/app/.composer_cache

WORKDIR /app

RUN apk add --no-cache --virtual .app-deps \
  bash \
  git \
  openssh-client \
  graphviz \
  ttf-dejavu \
  && addgroup -g 1000 -S developer \
  && adduser -u 1000 -S developer -G developer \
  && apk add --no-cache --virtual .build-deps \
    $PHPIZE_DEPS \
    libxml2-dev \
  && pecl install \
    xdebug-3.0.4 \
    apcu-5.1.18 \
  && docker-php-ext-install -j$(nproc) \
    simplexml \
    sockets \
  && docker-php-ext-enable \
    apcu \
  && curl -sS https://getcomposer.org/installer | php -- --install-dir=/usr/local/bin/ --filename=composer \
  && apk del .build-deps \
  && chmod +x /docker-entrypoint.sh

COPY ./.k8s/backend/php.ini /usr/local/etc/php/php.ini
COPY ./.k8s/backend/xdebug.ini /usr/local/etc/php/conf.d/15-xdebug.ini
COPY ./.k8s/backend/extra-php-fpm.conf /php-fpm.d/extra-php-fpm.conf
COPY ./.k8s/backend/php-fpm.conf /usr/local/etc/php-fpm.conf
