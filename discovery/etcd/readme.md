#### Etcd 服务发现

1、Key：`/smartproxy/${proxy_name}/${addr}-${randomstr}`<br>
2、Value：<br>

```json
{ "addr": "127.0.0.1:10003", "weight": 1 }
```

3、操作<br>

```bash
etcdctl put /smartproxy/proxy1/127.0.0.1:10003-abcdefg '{"addr": "127.0.0.1:10003", "weight":1}'
```
