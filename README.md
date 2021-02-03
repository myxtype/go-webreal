# webreal
基于websocket的客户端级的订阅和发布服务

语言表达有限，望见谅

# 背景

因为有些项目都用到这种模式：客户端需要实时得到服务端的一些数据，并且由客户端主动订阅。

比如：用户账户余额变动，服务器监听MySQL的BinLog日志，然后通过Redis订阅模型发布消息，推送服务器就可以监听Redis的订阅，通过此Lib来向客户端发出余额变动的消息。

# 流程

1. 客户端连接
2. 客户端订阅
3. 服务器向订阅数据的客户端发出消息

![img.png](img.png)

# 验证客户端

有些时候要对客户端的订阅信息进行Token验证，比如User1仅能订阅他自己的账户余额变动，User1订阅时需要携带他的Token信息，推送服务器进行验证，验证成功就订阅。
