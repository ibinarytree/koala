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

hello.proto
```
syntax = "proto3";
package hello;
//使用go_package指定生成的代码路径，下面的代码将会生成到github.com/ibinarytree/hello/proto_gen/hello这个目录下；
//大家可以根据自己项目的目录指定这个路径
//一般建议大家把proto文件和生成的代码统一进行管理，放到一个单独的代码库中；
option go_package="github.com/ibinarytree/hello/proto_gen/hello";

message HelloRequest {
  string name = 1;
}

message HelloResponse {
  string reply = 1;
}

service HelloService {
  //sayhello 方法
  rpc SayHello(HelloRequest) returns (HelloResponse){}
}

```
生成代码
```
koala -s -f hello.proto
```
注意：如果hello.proto有import其他proto文件，需要通过proto_path指定引入的proto文件的目录，通过-i文件生成引入的proto文件代码。

2. 生成客户端代码
```
koala -c -f hello.proto
```
3. 编译，生成代码之后，执行build.sh即可以编译。window执行build.bat
```
./build.sh 
```
4. 运行，编译之后的二进制程序统一放到了output目录，执行output目录下的start.sh即可以启动运行。window执行start.bat
```
cd output
./start.sh
```
5. 部署，把output目录打包成output.tar.gz，拷贝到服务端上，然后后台运行即可
```
nohup ./start.sh&
```
# Question
- 编译时，报etcd相关错误，比如:undefined: balancer.PickOptions；是因为grpc版本和etcd不兼容导致。解决方法如下：

打开go.mod, 把grpc的版本修改成v1.24.0
```
require (
    github.com/golang/protobuf v1.3.4
    github.com/ibinarytree/koala v1.9.9
    github.com/smartystreets/goconvey v1.6.4 // indirect
    golang.org/x/net v0.0.0-20200301022130-244492dfa37a
    google.golang.org/grpc v1.27.1
)
```
修改后的go.mod内容如下：
```
require (
    github.com/golang/protobuf v1.3.4
    github.com/ibinarytree/koala v1.9.9
    github.com/smartystreets/goconvey v1.6.4 // indirect
    golang.org/x/net v0.0.0-20200301022130-244492dfa37a
    google.golang.org/grpc v1.24.0
)
```
重新编译，问题解决了


