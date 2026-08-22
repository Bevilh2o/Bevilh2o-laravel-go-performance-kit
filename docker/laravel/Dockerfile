FROM php:8.4-cli-alpine

# Install system dependencies and PostgreSQL development libraries
RUN apk add --no-cache \
    postgresql-dev \
    autoconf \
    g++ \
    make \
    linux-headers \
    git \
    unzip \
    curl \
    && docker-php-ext-install pdo pdo_pgsql pcntl \
    && pecl install redis \
    && docker-php-ext-enable redis \
    && rm -rf /tmp/pear

WORKDIR /var/www/html

EXPOSE 8000

# Start Laravel development server
CMD ["php", "artisan", "serve", "--host=0.0.0.0", "--port=8000"]