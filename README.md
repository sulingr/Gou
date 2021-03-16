# Gou
仿Gin源码编写的Web框架

实现如下功能：

·设计上下文，并用以封装 Request 和 Response ，提供对 JSON、HTML 等返回类型的支持。
·使用 Trie 树实现动态路由解析，支持两种模式:name和*filepath。
·实现路由分组控制
