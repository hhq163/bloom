# bloom-filter
bloom-filter 布隆过滤器，支持使用本地内存和redis缓存。在网上看了别人的实现，发现redis版本实现有问题，使用单个连接进行处理，有并发问题，也没有做性能测试

1、性能测试数据：  
（1）将bit数组存储在本地内存  
阿里云服务器，ecs.c6.xlarge实例：4核，8G内存  
Intel Xeon(Cascade Lake) Platinum 8269CY  

布隆过滤器，数据存内存，10000000条数据,占用44M内存，错误率0.1%，压力测试20000次，tps为732064次/s，平均每次查询消耗1506ns  

（2）将bit数组存储在redis中  

2、参考以下项目：  
（1）https://github.com/bculberson/bloom.git
