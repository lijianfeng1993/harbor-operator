# harbor-operator
## Overview
Harbor Operator 提供了在kubernetes集群中快速部署高可用harbor的功能。

harbor-operator不提供数据的存储功能，在部署operator之前，需要用户自行准备高可用redis集群、pgsql集群提供数据存储。
目前Harbor-operator只提供通过对象存储保存镜像层数据的方式，后续会陆续增加其他存储对接方式。

本项目提供了测试环境部署redis、pgsql、minio等第三方存储服务的部署文件。

Table of Contents
=================

* [harbor-operator](#redis-operator)
    * [Overview](#overview)
    * [Prerequisites](#prerequisites)
    * [Quick Start](#quick-start)
        * [Deploy dependence](#deploy-dependence)
        * [Compile harbor operator](#compile-harbor-operator)
        * [Deploy harbor operator](#deploy-harbor-operator)
        * [Deploy harbor](#deploy-harbor)
        * [Cleanup](#cleanup)

## Prerequisites

* go version v1.13+
* kubernetes cluster v1.18
* redis
* pgsql
* minio
* kubebuilder v2.3.2

### Deploy dependence
部署所依赖的镜像都是harbor官方镜像

依赖部署文件位于deploy/dependency文件夹
* deploy redis
```
$ docker run -itd --name harbor-redis -p 6379:6379 -v /apps/data/harbor-redis/etc/redis.conf:/etc/redis.conf goharbor/redis-photon:v1.10.6
```
这里我们挂载了本地的配置文件redis.conf到容器中，配置文件中，我们改了默认的database 1000，因为默认的redis只提供了16个库，我们希望多个harbor复用这一个redis集群，因此需要redis支持更多的库。
* deploy pgsql
  使用docker-compose部署pgsql,镜像使用harbor中的harbor-db
```
version: '3.5'

services:
  postgres:
    container_name: harbor-pgsql
    image: goharbor/harbor-db:v1.10.6
    environment:
      POSTGRES_PASSWORD: 8cDcos11
    volumes:
       - postgres:/var/lib/postgresql/data/
    ports:
      - "5432:5432"
    networks:
      - postgres
    restart: unless-stopped

networks:
  postgres:
    driver: bridge

volumes:
    postgres:

```
```
$ docker-compose up -d
登陆pgsql，初始化用户server和signer
$ docker exec -it harbor-pgql bash
postgres [ / ]$ psql -U postgres
psql (9.6.19)
Type "help" for help.

postgres=# CREATE USER signer;
postgres=# alter user signer with encrypted password 'password';
postgres=# CREATE USER server;
postgres=# alter user server with encrypted password 'password';
postgres=# 
postgres=# select * from pg_user;
 usename  | usesysid | usecreatedb | usesuper | userepl | usebypassrls |  passwd  | valuntil | useconfig 
----------+----------+-------------+----------+---------+--------------+----------+----------+-----------
 postgres |       10 | t           | t        | t       | t            | ******** |          | 
 signer   |    16387 | f           | f        | f       | f            | ******** |          | 
 server   |    16385 | f           | f        | f       | f            | ******** |          | 
(3 rows)
 
```
* deploy minio
```
$ docker run  -d -p 9000:9000 --name minio -e "MINIO_ACCESS_KEY=admin" -e "MINIO_SECRET_KEY=8cDcos11" -v /apps/data/minio/data:/data -v /apps/conf/minio/config:/root/.minio minio/minio server /data
```

### Compile harbor operator
```
# 编译代码，可执行文件harbor-operator位于bin目录
$ make build

# 构建镜像
$ make image
```

### Deploy harbor operator
```
$ kubectl create -f deploy/crds/harbor.example.com_harborservices_crd.yaml
$ kubectl create -f deploy/configmap.yaml
$ kubectl create -f deploy/clusterrole_binding.yaml
$ kubectl create -f deploy/operator.yaml

$ kubectl get pod | grep operator
harbor-operator-5cd887c779-pcldf                1/1     Running   0          20m
```

### Deploy harbor
* 测试cr内容
```
$ cat testharbor.yaml

apiVersion: harbor.example.com/v1
kind: HarborService
metadata:
  name: testharbor
spec:  
  instanceInfo:
    instanceName: "testharbor"
    instanceType: "harbor"
    nodePortIndex: 32190
    redisDbIndex: 30
    s3Config:
      bucket: "test1"
      accesskey: "admin"
      secretkey: "8cDcos11"
```

* 创建cr
```
[root@master1 crds]# kubectl create -f deploy/crds/testharbor.yaml 
harborservice.harbor.example.com/testharbor created
```

* 查看harborservice部署情况
```
[root@master1 operator]# kubectl get harborservice
NAME         AGE
testharbor   105s

[root@master1 operator]# kubectl describe harborservice testharbor
Name:         testharbor
Namespace:    default
Labels:       <none>
Annotations:  <none>
API Version:  harbor.example.com/v1
Kind:         HarborService
......
Spec:
  Instance Info:
    Instance Name:    testharbor
    Instance Type:    harbor
    Node Port Index:  32190
    Redis Db Index:   30
    s3Config:
      Accesskey:  admin
      Bucket:     test1
      Secretkey:  8cDcos11
Status:
  Condition:
    Phase:       running
  External URL:  http://testharbor.harbor.com:32190
Events:
  Type    Reason                                Age    From             Message
  ----    ------                                ----   ----             -------
  Normal  syncDatabase                          2m13s  harbor-operator  Start sync pgsql database and table
  Normal  syncDatabaseSuccess                   2m12s  harbor-operator  Sync pgsql database success
  Normal  createS3Bucket                        2m12s  harbor-operator  Start create minio bucket for harborservice
  Normal  createS3BucketSuccess                 2m12s  harbor-operator  Sync pgsql database success
  Normal  initKubernetesResources               2m12s  harbor-operator  Start init kubernetes namespace, pv, pvc for harborservice
  Normal  initKubernetesResourcesFailedSuccess  2m12s  harbor-operator  Init kubernetes resource for harborservice success
  Normal  deployIntance                         2m12s  harbor-operator  Start deploy harbor service
  Normal  deployIntanceSuccessd                 2m8s   harbor-operator  Deploy harbor service by successd.


[root@master1 operator]# kubectl get pod -n testharbor
NAME                                               READY   STATUS    RESTARTS   AGE
testharbor-harbor-chartmuseum-69cd4987d9-c4k86     1/1     Running   0          2m52s
testharbor-harbor-clair-7bd9667c86-xxthk           2/2     Running   0          2m52s
testharbor-harbor-core-5566594f99-8g6fl            1/1     Running   0          2m52s
testharbor-harbor-jobservice-7dcf5c7bfc-7lwrw      1/1     Running   0          2m53s
testharbor-harbor-nginx-66858f8dd4-btnk5           1/1     Running   0          2m53s
testharbor-harbor-notary-server-6f64659767-r7rvv   1/1     Running   0          2m52s
testharbor-harbor-portal-7b9884bd6-rrkb9           1/1     Running   0          2m52s
testharbor-harbor-registry-6797ccdd54-z4j45        2/2     Running   0          2m53s
```

* 验证harbor功能
```
# 配置本地docker，添加--insecure-registry=0.0.0.0/0

# 配置docker客户端节点/etc/hosts文件，添加
10.142.113.231 testharbor.harbor.com

# 验证docker上传镜像
[root@node1 ~]# docker tag goharbor/nginx-photon:v2.1.3 testharbor.harbor.com:32190/library/nginx-photon:v2.1.3
[root@node1 ~]# docker push testharbor.harbor.com:32190/library/nginx-photon:v2.1.3
The push refers to repository [testharbor.harbor.com:32190/library/nginx-photon]
e7ecd452926a: Pushed 
72021dc640d8: Pushed 
v2.1.3: digest: sha256:cf7e4311220b44f6d03b093028a69a24613fac6b47bbc16c7f50857116a2f161 size: 6914
```

* 登陆页面查看镜像详情
  ![image](https://raw.githubusercontent.com/lijianfeng1993/harbor-operator/master/image/testharbor.png)

### Cleanup
* 删除cr资源
```
[root@master1 crds]# kubectl delete -f testcr.yaml 
harborservice.harbor.example.com "testharbor" deleted
```

* 查看服务
```
[root@master1 ~]# kubectl get harborservice
No resources found in default namespace.
[root@master1 ~]# 
[root@master1 ~]# kubectl get pod -n testharbor
No resources found in testharbor namespace.


[root@node1 ~]# docker login testharbor.harbor.com:32190
Authenticating with existing credentials...
Login did not succeed, error: Error response from daemon: Get http://testharbor.harbor.com:32190/v2/: dial tcp 10.142.113.231:32190: connect: connection refused
```