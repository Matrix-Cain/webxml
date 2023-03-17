# webxml
## 介绍
这是一个简单的利用neo4j进行污点分析的练习项目, 主要是为了简单上手neo4j。 这个demo通过读取java web应用中的web.xml文件，解析其中的filter和servlet等信息。并将其以节点的形式保存在neo4j数据库中，通过neo4j的数据可视化后，我们可以直观的查看一个web应用的路由结构(这里只考虑了web.xml中的路由，通过注解指定的路由无法解析)

例如：在非verbose模式下到处对应的web.xml数据存入neo4j中，在neo4j bloom中执行
`MATCH p=()-[]->() RETURN p`

![image-1](https://github.com/zwh-china/webxml/raw/master/img/1.png)

通过简单查看可以发现存在一些鉴权filter，例如图中选中的`CheckIsLoggedFilter`

假设对代码进行代码审计后，我们发现某些servlet存在命令执行风险，例如一个名为`evalServlet`, 这时开启verbose模式，先在neo4j中执行

```neo4j
MATCH (n) DETACH DELETE n
```

删除先前插入数据，再导入数据。通过执行

```neo4j
MATCH path=()-[r:evalServlet]->()
RETURN path
```

或者(需要安装neo4j的apoc拓展)

```neo4j
match (source:App) // 添加where语句限制source函数
match (sink:Servlet {name:"evalServlet"}) // 添加where语句限制sink函数
call apoc.algo.allSimplePaths(source, sink, 'evalServlet>', 20) yield path // 查找具体路径,20代表深度，可以修改
return path
```

![image-2](https://github.com/zwh-china/webxml/raw/master/img/2.png)

![image-3](https://github.com/zwh-china/webxml/raw/master/img/img/3.png)

可以看到存在一条路径绕开了`CheckIsLoggedFilter`这个鉴权filter。通过此方式辅助查找未授权等漏洞

我们也可以在情况更加复杂的情况下执行

```neo4j
MATCH (n:Filter WHERE n.name="CheckIsLoggedFilter")<-[rels:evalServlet]-()
with collect(rels) as rels
MATCH path=()-[r:evalServlet where none (r1 in rels where r.url=r1.url)]->()
return path
```

筛选**未经过**`CheckIsLoggedFilter`节点的可达路径

![image-4](https://github.com/zwh-china/webxml/raw/master/img/4.png)

同理如果我们要查找**经过**`CheckIsLoggedFilter`的路径，我们可以执行

```neo4j
MATCH (n:Filter WHERE n.name="CheckIsLoggedFilter")<-[rels:evalServlet]-()
with collect(rels) as rels
MATCH path=()-[r:evalServlet where not none(r1 in rels where r.url=r1.url)]->()
return path
```

![image-5](https://github.com/zwh-china/webxml/raw/master/img/5.png)
