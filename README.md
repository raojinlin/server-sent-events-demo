# server-sent-events-demo

Server-Sent Events go语言实现demo。SSE的具体的使用可以参考: [Clogs](https://github.com/raojinlin/clogs)

## web实时通信机制

* Websocket
* Server-Sent Events
* 轮询

### 轮询

顾名思义，轮询就是在某个时间间隔内定期向服务器发送请求。其中轮询有分为短轮询和长轮询。这是一种客户端主动请求的方式。

#### 短轮询

定期向服务器请求，无论请求的资源是否可用，服务器都会尽快响应，客户端再次发起下次轮询。这种方式会比较消耗网络带宽，如果资源一直不可用就会有很多不必要的请求发送到服务器。

![](https://article.biliimg.com/bfs/article/04d62b2b38b5f447d1eeb20110cc4111dad8cc86.png)


#### 长轮询
定期向服务器发送请求，与短轮询不同的是，在资源不可用时长轮询不会立即将连接关闭，而是会等待资源可用后在响应客户端。或者等待了一段时间资源任然不可用（超时）服务器将连接关闭，客户端等待一段时间后再次发起请求。
与短轮询相比，长轮询更高效一些，请求数量减少了很多。

![](https://article.biliimg.com/bfs/article/3c6811555f7109596b6d4e205360d7c6bf8eb67c.png)

#### Websocket

![](https://article.biliimg.com/bfs/article/3c6811555f7109596b6d4e205360d7c6bf8eb67c.png)

在客户端和服务器打开交互式的通信会话。这是一种全双工通信，客户端与服务器会建立一个持久连接，服务器可以主动发送数据给客户端。客户端可以通过监听事件来处理来自服务器的消息。与轮询的方式相比，大大减少了延迟，没有了数据更新的往返时间。

#### Server-Sent Events

服务器发送事件，SSE会建立一个持久的HTTP连接。建立连接后服务器可以主动往客户端推送数据。与websocket不同，这是一种单向通信的方式，即建立连接后客户端不能向服务器发送数据。

### Server-Sent Events
> 通常来说，一个网页获取新的数据通常需要发送一个请求到服务器，也就是向服务器请求的页面。使用服务器发送事件，服务器可以随时向我们的 Web 页面推送数据和信息。这些被推送进来的信息可以在这个页面上以 事件 + 数据 的形式来处理。

![](https://article.biliimg.com/bfs/article/e34404062b6adc3355411ed2a1610052640b3694.png)

#### 使用场景

实时数据更新：SSE适用于需要在客户端实时更新数据的场景。例如，股票市场行情、实时新闻更新、社交媒体的实时通知等。

实时监控和通知：SSE可用于监控系统、设备或传感器的实时状态。服务器可以将实时数据推送给客户端，以便实时监测并及时采取相应行动。例如，实时监控温度、湿度、能源消耗等。

实时协作和协同编辑：SSE可以在协作和协同编辑工具中提供实时更新功能。多个用户可以同时编辑同一个文档，并通过SSE实时更新对方的更改，以实现实时协作和协同编辑的效果。

消息推送和通知：SSE可用于发送实时消息和通知给订阅的用户。例如，推送新邮件通知、提醒用户活动提醒、推送定制的实时提醒等。

事件流是一个简单的文本数据流，文本应该用UTF8格式编码。事件流的消息由两个换行符分开，以冒号开头的行为注释行，会被忽略。

#### 事件流字段

![](https://article.biliimg.com/bfs/article/a566e50b0051cd78a28c3ec641865d54231214f4.png)

* event: 用于标识事件类型的字符串，如果没有指定event，浏览器默认认为是message。
* data: 消息的数据字段，当EventSource收到多个已```data:```开头的连续行是，会将它们连接起来，在它们之间插入一个换行符。，末尾的换行符也会被删除。
* id: 事件ID，会被设置为当前EventSource对象的内部属性“最后一个事件ID”的值。
* retry: 重新连接的时间。如果与服务器的连接丢失，浏览器会等待指定的时间，然后重新连接。retry必须是一个整数，它的单位是毫秒。

#### 使用

SSE的用法比较简单，只需要在服务器编写一些代码将事件流传输到前端。前端使用EventSource来监听这些事件流。

使用JavaScript创建一个EventSource对象，连接到localhost:8083/events。监听消息流的message事件（如果没有指定event，那么event默认是message）

EventSource使用


#### 服务器发送事件流

服务器端发送事件流也比较简单，这里使用Go语言做了一个demo。

首先我们创建一个名字为Event的struct，表示SSE中的一个事件。这个struct有两个方法

* `String() string`，将Event编码为字符串
* `Bytes() []byte`，将Event编码为字节切片

![](https://article.biliimg.com/bfs/article/b371e1ed17e63916d4862fd239667cfb2d8f93a2.png)

创建一个HTTP 控制器，负责输出事件流。下面这个stream函数会输出100个事件到客户端，每隔1秒钟输出一个事件。这里注意，一个事件写入到连接后不要忘记把缓存数据flush掉，这样客户端才能看到最新推送的数据。


![x](https://article.biliimg.com/bfs/article/64ceff601c636f4426809085c6f150a5b0d1a4cb.png)


### 注意事项
* nginx配置
  * 使用nginx做反向代理时需要将proxy_buffering关闭
    * proxy_buffering off
  * 或者加上响应头部x-accel-buffering，这样nginx就不会给后端响应数据加buffer
    * x-accel-buffering: no
* EventSource
  * 连接关闭后会自动重连
  * 需要显示的调用close方法
  * EventSource.prototype.close
 


### 参考
[MDN - Using_server-sent_events](https://developer.mozilla.org/zh-CN/docs/Web/API/Server-sent_events/Using_server-sent_events)
