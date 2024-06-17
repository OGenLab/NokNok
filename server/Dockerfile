########################################
## 编译源码 ############################
########################################
FROM golang:1.19.2 as builder
MAINTAINER ""

ENV     APPROOT /usr/local/service/nokserver
ENV     LD_LIBRARY_PATH $APPROOT/lib:$LD_LIBRARY_PATH
RUN     mkdir -p $APPROOT
SHELL   ["/bin/bash", "-c"]
WORKDIR $APPROOT

######### copy workspace ####begin######
####### edit .dockerignore if need #####
####### edit .dockerignore if need #####
COPY    .   . 
####### edit .dockerignore if need #####
####### edit .dockerignore if need #####
######### copy workspace #####end#######

RUN     pwd; ls -l;                                       \
        echo "======Start make in Dockerfile =========";  \
        make tar &&                                       \
        echo "======END make in Dockerfile ==========="


########################################
## 生成运行时image #################
########################################
FROM golang:1.19.2 
ENV         TZ              Asia/Shanghai
ENV         APPROOT         /usr/local/service/nokserver
ENV         OUTROOT         /usr/local/service/nokserver/_output/nokserver
ENV         LD_LIBRARY_PATH $APPROOT/lib:$LD_LIBRARY_PATH
WORKDIR     $APPROOT
COPY        --from=builder $OUTROOT .

ENTRYPOINT ["./bin/nokserver"]
