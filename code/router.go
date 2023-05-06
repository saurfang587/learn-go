package code

import (
	"fmt"
	"net/http"
)

//type ServeMux struct {
//	mu    sync.RWMutex        //读写互斥锁 用来保护map和切片的并发安全
//	m     map[string]muxEntry //装载路径与方法的MAP，Key是路径，Value是一个结构体。
//	es    []muxEntry          //这个切片里的结构体是按照带路径的长度从长到短排好序的，作用是将一类请求归到同一方法中
//	hosts bool                //标志位，当前的路径中有没有根路径，只起到一个匹配优先级的问题
//}
//
//type muxEntry struct {
//	h       Handler //路径对应的方法
//	pattern string  //这里冗余了一分路径
//}

func TestRouter1() {
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		_, _ = fmt.Fprintf(writer, "关注香香编程喵喵喵，关注香香编程谢谢喵喵喵！")
	})
	http.HandleFunc("localhost/a123", func(writer http.ResponseWriter, request *http.Request) {
		//会使 mux的hosts变位true  访问http://localhost:8081/a123 会解析到这个方法
		_, _ = fmt.Fprintf(writer, "localhost/a123")
	})
	http.HandleFunc("/b/c", func(writer http.ResponseWriter, request *http.Request) {
		_, _ = fmt.Fprintf(writer, "/index/a/b/c/")
	})
	http.HandleFunc("/index/a/b/c", func(writer http.ResponseWriter, request *http.Request) {
		_, _ = fmt.Fprintf(writer, "/index/a/b/c/")
	})
	http.HandleFunc("/z/x/", func(writer http.ResponseWriter, request *http.Request) {
		//路径 /z/x/111 和 /z/x/222 都会解析到这个方法
		_, _ = fmt.Fprintf(writer, "/z/x/")
	})
	fmt.Printf("http mux:%+v\n", http.DefaultServeMux)
	panic(http.ListenAndServe(":8081", nil))
}

//func (mux *ServeMux) Handle(pattern string, handler Handler) {
//	mux.mu.Lock() //互斥锁的经典用法，后续有对MAP和切片的操作，需要保证并发安全
//	defer mux.mu.Unlock()
//
//	if pattern == "" {
//		panic("http: invalid pattern") //路径不能为空
//	}
//	if handler == nil {
//		panic("http: nil handler") //方法不能为空
//	}
//	if _, exist := mux.m[pattern]; exist {
//		panic("http: multiple registrations for " + pattern) //路径不能重复
//	}
//
//	if mux.m == nil {
//		mux.m = make(map[string]muxEntry) //创建MAP
//	}
//	e := muxEntry{h: handler, pattern: pattern} //创建存储主体
//	mux.m[pattern] = e                          //赋值
//	if pattern[len(pattern)-1] == '/' {         //如果当前路径的结尾是'/'
//		mux.es = appendSorted(mux.es, e) //把当前路径放到es里，并按照字符串长度排序
//	}
//
//	if pattern[0] != '/' { //如果路径的头部不是'/'
//		mux.hosts = true //就把hosts的标志位设置为true
//	}
//}

//func (mux *ServeMux) handler(host, path string) (h Handler, pattern string) {
//	mux.mu.RLock()         //同样的，先把锁加上。
//	defer mux.mu.RUnlock() //注意这里的逻辑，他并没有在match函数中加锁，而是在这一层就开始加锁了。
//
//	// 如果存在主机名，就优先采用全匹配的方式 参见后续案例
//	if mux.hosts {
//		h, pattern = mux.match(host + path)
//	}
//	// 正常匹配，如果没有匹配到，就返回一个默认方法
//	if h == nil {
//		h, pattern = mux.match(path)
//	}
//	if h == nil {
//		h, pattern = NotFoundHandler(), ""
//	}
//	return
//}

//func (mux *ServeMux) match(path string) (h Handler, pattern string) {
//	// 先判断是否Map中是否存在这个路由
//	v, ok := mux.m[path]
//	if ok {
//		return v.h, v.pattern
//	}
//
//	//如果Map里不存在，就在切片中匹配下前缀路由
//	for _, e := range mux.es {
//		if strings.HasPrefix(path, e.pattern) {
//			return e.h, e.pattern
//		}
//	}
//	//都不存在，返回空
//	return nil, ""
//}

//type RouterGroup struct {
//	// Handlers 一个函数链，底层是一个数组，装的是中间件的方法，子Group会继承父的中间件
//	Handlers HandlersChain
//	// basePath 定义的Group的主路径
//	basePath string
//	// engine 核心引擎的指针。需要注意，基于主Group创造出来的子Group们不仅会使用主中间件还会共用一个主引擎
//	engine *Engine
//	// root 是否为根Group，只有我们New出来的主引擎自带的Group这里才会是true
//	root bool
//}
//
//func (group *RouterGroup) Use(middleware ...HandlerFunc) IRoutes {
//	//非常简单的操作，将中间件方法Append到Handlers中
//	group.Handlers = append(group.Handlers, middleware...)
//	return group.returnObj()
//}
//
//// Group 该方法会返回一个全新的RouterGroup，并且会继承父Group的一些属性
//func (group *RouterGroup) Group(relativePath string, handlers ...HandlerFunc) *RouterGroup {
//	return &RouterGroup{
//		//将多个中间件按照顺序组合到一起，底层使用的是内存copy,我们讲到中间件时会展开
//		Handlers: group.combineHandlers(handlers),
//		//将父级的路径组合到一起
//		basePath: group.calculateAbsolutePath(relativePath),
//		//把主引擎的指针继续传递下去
//		engine: group.engine,
//	}
//}
//
//// handle 核心中的核心。POST，GET等注册路由器方法的底层函数。
//func (group *RouterGroup) handle(httpMethod, relativePath string, handlers HandlersChain) IRoutes {
//	//获取当前的绝对路径相当于 Group的路径加上传入的路径
//	absolutePath := group.calculateAbsolutePath(relativePath)
//	//将当前函数和Group的中间件函数们串起来，得到一个新的切片
//	handlers = group.combineHandlers(handlers)
//	//将上述获得的绝对路径+方法链注册到路由树里。注意这里调用的是主引擎的addRoute方法。也就是说，路由树的管理权在主引擎手上。
//	group.engine.addRoute(httpMethod, absolutePath, handlers)
//	return group.returnObj()
//}
//
//// returnObj 函数定义是返回一个接口，实际返回的是自身。这是实现链式操作的一种方式
//func (group *RouterGroup) returnObj() IRoutes {
//	//root 字段的唯一用处，在根Group的情况下，返回主引擎。注意，主引擎其实也是IRoutes的实现，你懂得：）。
//	if group.root {
//		return group.engine
//	}
//	return group
//}

type methodTree struct {
	method string //方法类型：GET，POST等
	//root   *node  //方法类型对应的具体路由树
}

type methodTrees []methodTree //主引擎中存储的路由树

// node 树的节点结构，适当调整了下字段顺序
//type node struct {
//	path      string // 节点路径，比如上面的blog，a，bout等等
//	fullPath  string //完整路径
//	wildChild bool   //通配符标识,比如带有:id的节点，这个值就是true
//
//	//字母索引，对应的就是上边的a，下面的about和article的b和r
//	//也就意味着，这个分支下面还有至少两个分支
//	indices string
//	//优先级，当前节点的子节点越多，这个值会越大
//	priority uint32
//	//当前节点的所有子节点
//	children []*node
//	//节点类型 有四种 static root param catchAll
//	// static: 静态节点（默认），
//	// root: 根节点
//	// catchAll: 带有通配符的节点，也就是*
//	// param: 参数节点
//	nType nodeType
//
//	handlers HandlersChain //函数链
//
//}

//func (engine *Engine) addRoute(method, path string, handlers HandlersChain) {
//	...
//	// 获取请求方法对应的树
//	root := engine.trees.get(method)
//	if root == nil {
//		//不存在的话，就创建一个该类型的树。
//		root = new(node)
//		root.fullPath = "/"
//		engine.trees = append(engine.trees, methodTree{method: method, root: root})
//	}
//	//将当前路径和方法添加到前缀树里。
//	root.addRoute(path, handlers)
//	...
//}

// addRoute 这个无法并发安全。言外之意就是，这个前缀树结构不是并发安全的。
//注意，gin框架经历了多次优化更新，这里采用的是1.9的版本。比着之前，代码量和整体逻辑简化不少
//func (n *node) addRoute(path string, handlers HandlersChain) {
//	fullPath := path
//	//对路由树的 priority 属性进行自增操作，表示在此路由树中添加了一个新的节点。
//	n.priority++
//
//	// 如果路由树还没有任何节点，那么直接将传入的路由路径和处理函数添加到当前节点，此时路由树的类型被设置为 root。
//	if len(n.path) == 0 && len(n.children) == 0 {
//		n.insertChild(path, fullPath, handlers)
//		n.nType = root
//		return
//	}
//
//	parentFullPathIndex := 0
//
//walk:
//	//walk循环 Go代码块的一种用法，通常会结合着goto使用。一般业务代码中不建议使用
//	//该循环的作用是寻找当前路径最合适的插入位置，本质上就是一种树的遍历，遍历的同时，适当调整树的结构。
//	for {
//		//查找当前路由节点和传入的路由路径的最长公共前缀
//		i := longestCommonPrefix(path, n.path)
//
//		//如果发现有前缀，并且前缀和当前前缀还小，相当于从 abc->ab,就会触发节点的分裂
//		if i < len(n.path) {
//			//分裂为新的节点，主要是path发生变化
//			...
//		}
//
//		// 创建一个当前路径的节点
//		if i < len(path) {
//			path = path[i:]
//			c := path[0]
//
//			// 处理参数类型后 还带有/的特殊情况，直接略过
//			if n.nType == param && c == '/' && len(n.children) == 1 {
//				...
//				continue walk
//			}
//
//			// 按照前缀继续往下查，直到没有重复前缀为止 此时的节点就是我们需要添加的位置了
//			for i, max := 0, len(n.indices); i < max; i++ {
//				...
//			}
//
//			// 如果此时的路径还有东西，并且类型不是通配符，可以认为是一个正常类型 我们的about 就会走到这里
//			if c != ':' && c != '*' && n.nType != catchAll {
//				// 向父节点上添加上子路径的首字母
//				n.indices += bytesconv.BytesToString([]byte{c})
//				//构造一下子节点
//				child := &node{
//					fullPath: fullPath,
//				}
//				//向当前节点上增加一个子节点
//				n.addChild(child)
//				n.incrementChildPrio(len(n.indices) - 1)
//				//将游标指向新节点，准备结束循环
//				n = child
//			} else if n.wildChild {
//				// 如果是通配符类型的，需要检查下是否有冲突，整体逻辑类似
//				...
//			}
//			//向当前节点插入一个路径，以上所有的操作都只是找到位置，这里才是真正的插入节点
//			n.insertChild(path, fullPath, handlers)
//			return
//		}
//
//		// 当前节点已经被注册过了
//		if n.handlers != nil {
//			panic("handlers are already registered for path '" + fullPath + "'")
//		}
//		//如果当前节点根本就不需要走前缀索引，那么直接就可以创建节点了
//		n.handlers = handlers
//		n.fullPath = fullPath
//		return
//	}
//}
//
//func (n *node) getValue(path string, params *Params, skippedNodes *[]skippedNode, unescape bool) (value nodeValue) {
//	var globalParamsCount int16
//
//walk:
//	// 遍历整棵树，第一次遍历时，n是根路径
//	for {
//		prefix := n.path
//		//在前缀树里查一下当前路径的位置
//		if len(path) > len(prefix) {
//			//前缀能匹配上
//			if path[:len(prefix)] == prefix {
//				...
//				// 遍历所有非通配符的节点
//				for i, c := range []byte(n.indices) {
//					...
//				}
//
//				// 非通配符节点
//				if !n.wildChild {
//					...
//				}
//
//				// 处理通配符节点
//				n = n.children[len(n.children)-1]
//				globalParamsCount++
//
//				switch n.nType {
//				case param:
//					...
//				case catchAll:
//					...
//				default:
//					panic("invalid node type")
//				}
//			}
//		}
//
//		//匹配到了具体的节点
//		if path == prefix {
//			//如果当前节点没有handler 或者当前路径不是 就需要回滚到最后一个有效的skippedNode
//			if n.handlers == nil && path != "/" {
//				...
//			}
//
//			//如果当前节点有handler 就算是匹配到具体的路由了
//			if value.handlers = n.handlers; value.handlers != nil {
//				value.fullPath = n.fullPath
//				return
//			}
//
//			...
//
//			//继续查询带有/的节点
//			for i, c := range []byte(n.indices) {
//				...
//			}
//
//			return
//		}
//
//		//判定下value.tsr 的值，这是个赋值语句……
//		value.tsr = path == "/" ||
//			(len(prefix) == len(path)+1 && prefix[len(path)] == '/' &&
//				path == prefix[:len(prefix)-1] && n.handlers != nil)
//
//		// 回退到最近一个有效的节点
//		if !value.tsr && path != "/" {
//			...
//		}
//
//		return
//	}
//}
