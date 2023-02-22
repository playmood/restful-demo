# Host 服务模块

## IMPL

Host Service的具体实现，上层业务基于Service进行编程，面向接口

Host Service使用方式：
+ 用于内部模块调用，基于此封装更高一层的业务逻辑，比如发布服务
+ Host Service对外暴露： 
+ http协议(暴露给用户)
+ gRPC(暴露给内部服务)