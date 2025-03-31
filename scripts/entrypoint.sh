#!/bin/sh

# Create Redis config directory
mkdir -p /usr/local/etc/redis

# Allow Redis to accept connections from any IP address
echo "bind 0.0.0.0" > /usr/local/etc/redis/redis.conf

echo "REDIS_PASSWORD: $REDIS_PASSWORD"
# Set Redis password from environment variable
echo "requirepass $REDIS_PASSWORD" >> /usr/local/etc/redis/redis.conf

# Add default user with full access and no password (for testing)
echo "user default on nopass ~* +@all" > /usr/local/etc/redis/users.acl

# Add custom user with password and full access
echo "user $REDIS_USER on >$REDIS_USER_PASSWORD ~* +@all" >> /usr/local/etc/redis/users.acl

# Start Redis with custom config and ACL
exec redis-server /usr/local/etc/redis/redis.conf --aclfile /usr/local/etc/redis/users.acl

