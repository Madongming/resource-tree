# resource-tree
通用组织资源的树

# 目录说明:
* dao 对数据CURD
* cache 缓存，缓存整棵资源树（常驻），每个用户自己的资源树（LRU）及每个资源与其他资源的关系图（LRU），加快读取速度
* model 数据结构定义
* tools 一些工具
* global 需要初始化并全局使用的内容，如数据库链接，日志，配置，错误定义

# 简单介绍
## 功能
* 以树形结构管理资源
* 资源的属性这里不做管理，只有一个Key做为在其他系统中此资源的标识
* 用户组，用户的管理
* 将节点授权给用户组或用户
* 微服务的支持。给节点设置关联关系，自动生成此节点相关的图，展示微服务的调用关系。
## 特点
* 由于只存储简单的属性和key，每个树节点的大小一般在500字节以下，如果10000的节点，一共占500w字节=5000k=5M的内存，所以将他们全部放在缓存中。
* 对这些节点做了对与ID，树形，图形的索引放在缓存中。
* 将树的缓存放在全局变量中。并保证线程安全，以及最终一致性。
* 基于用户权限的树放在LRU缓存中。
* 基于节点的图。因为多个节点会映射到一张图中，所以保存一个中间Key值，将图放在两级LRU缓存中。
## 配置
目前均为必填，并且大于且接近真实数目最佳。后期再进行优化。
* `user_cache_size` 设置 用户自己的有权限的树的LRU缓存大小，一般为DAU即可
* `resource_cache_size` 设置 节点总数，
* `graph_cache_size` 设置 图的数目
