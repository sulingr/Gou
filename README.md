# Gou
使用Go编写，仿造Gin源码编写的Web框架

实现如下功能：
* 提供对 JSON、HTML 等返回类型的支持。
* 使用 Trie 树实现动态路由解析，支持两种模式:name和*filepath。
* 实现路由分组控制。
* 设计并实现 Web 框架的中间件机制，并实现Logger中间件（记录请求到响应的时间）
* 实现静态资源服务
* 支持HTML模板渲染