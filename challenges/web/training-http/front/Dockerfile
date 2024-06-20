FROM    php:7.3-apache

RUN     apt-get update \
        && apt-get install -y libicu-dev \
        && docker-php-ext-configure intl \
        && docker-php-ext-install intl \
        && docker-php-ext-install gettext \
        && apt-get install -y locales-all

COPY    www/ /var/www/html/
COPY    translations/ /var/www/translations/
COPY    on-init /pwinit/