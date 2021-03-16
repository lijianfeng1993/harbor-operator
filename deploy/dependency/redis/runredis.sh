#!/bin/bash
docker run -itd --name harbor-redis -p 6379:6379 -v /apps/data/harbor-redis/etc/redis.conf:/etc/redis.conf goharbor/redis-photon:v1.10.6