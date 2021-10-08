FROM git.test.com:5005/devops/docker/php:8.0.10

# need to /usr/local/etc/php-fpm.conf 
ENV PHP_FCGI_CHILDREN=10
ENV PHP_FCGI_MAX_REQUESTS=100
ENV PHP_FCGI_TERMINATE_TIMEOUT=60

# phpconfig
COPY ./.gitlab-ci/backend/php-fpm.conf /usr/local/etc/php-fpm.conf
COPY ./.gitlab-ci/backend/extra-php-fpm.conf /php-fpm.d/extra-php-fpm.conf

WORKDIR /app
USER www-data

COPY . /app

CMD php-fpm
