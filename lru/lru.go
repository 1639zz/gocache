package lru

import "container/list"

//创建cache结构体
type Cache struct {
	//最大内存
	maxBytes int64
	//已使用的内存
	nbytes int64
	//双向链表
	ll *list.List
	//缓存map
	cache map[string]*list.Element
	//被删除的回调函数
	OnEvicted func(key string,value Value)
}
type entry struct {
	key string
	value Value
}
type Value interface {
	Len() int
}
//实现New()函数
func New(maxBytes int64, onEvicted func(string, Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}
//实现查找功能
func (c *Cache) Get(key string)(value Value,ok bool)  {
	//从字典中查找对应双向链表的节点
	if ele, ok := c.cache[key]; ok {
		//如果存在，将该节点移到队尾
		c.ll.MoveToFront(ele)
		//得到kv键值对
		kv := ele.Value.(*entry)
		//返回查到的value值
		return kv.value,true
	}
	return
}
//实现删除操作,移除最近访问最少的节点
func (c *Cache) RemoveOldest()  {
	//得到首节点
	ele := c.ll.Back()
	//判断ele是否为空，为空则删除
	if ele != nil {
		//删除节点
		c.ll.Remove(ele)
		//得到kv键值对
		kv := ele.Value.(*entry)
		//删除c.cache的映射关系
		delete(c.cache,kv.key)
		//更新已使用内存长度
		c.nbytes-=int64(len(kv.key))+int64(kv.value.Len())
		//如果回调函数不等于nil。则调用回调函数
		if c.OnEvicted !=nil {
			c.OnEvicted(kv.key,kv.value)
		}
	}
}
//实现新增修改操作
func (c *Cache)Add(key string,value Value)  {
	//如果键存在，则更新对应节点的值
	if ele,ok:=c.cache[key];ok {
		//将该节点移到队尾
		c.ll.MoveToFront(ele)
		kv :=ele.Value.(*entry)
		//更新内存长度
		c.nbytes+=int64(len(kv.key))-int64(kv.value.Len())
		kv.value=value
	}else{
		//不存在则新增 else
		c.ll.PushFront(&entry{key,value})
		c.cache[key]=ele
		//更新 c.nbytes
		c.nbytes+=int64(len(key))+int64(value.Len())
	}
	//如果超过了设定的最大值 c.maxBytes，则移除最少访问的节点。
	for c.maxBytes!=0&&c.maxBytes<c.nbytes {
		c.RemoveOldest()
	}
}
//得到添加数据条数
func (c *Cache)Len() int{
	return c.ll.Len()
}

