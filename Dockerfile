# Use distroless as minimal base image to package the manager binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
#FROM registry.cn-hangzhou.aliyuncs.com/chenvista/static:nonroot
FROM frostmourner/alpine-ca:3.5

WORKDIR /

ADD bin/harbor-operator /harbor-operator

ADD deploy/harbor-helm-1.5.3 /tmp/harbor-helm-1.5.3

RUN chmod -R 777 /tmp

ENTRYPOINT ["/harbor-operator"]
