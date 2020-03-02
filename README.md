# koala
koala是一个基于grpc的高并发、高性能的微服务框架，提供了非常丰富的功能。

# Installation



- Go Version >= 1.11
- Grpc Version: google.golang.org/grpc v1.24.0
- protoc Version >= 3.11.0, 安装地址：https://github.com/protocolbuffers/protobuf/releases

    ```
    
    go get github.com/ibinarytree/koala
    go get github.com/ibinarytree/koala/tools/koala

    ``````



# Usage 
    1. 生成服务端代码
    ```
    koala -s -f xxx.proto
    ```
    2. 生成客户端代码
    ```
    koala -c -f xxx.proto
    ```











