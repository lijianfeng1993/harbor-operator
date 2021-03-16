#!/bin/bash
docker run  -d -p 9000:9000 --name minio -e "MINIO_ACCESS_KEY=admin" -e "MINIO_SECRET_KEY=8cDcos11" -v /apps/data/minio/data:/data -v /apps/conf/minio/config:/root/.minio minio/minio server /data