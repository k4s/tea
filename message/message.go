package message

import (
	"sync"
	"sync/atomic"
)

//Message封装了来回交换的messages
type Message struct {
	Header []byte
	ConnID uint32
	Body   []byte
	bsize  int
	refcnt int32
	pool   *sync.Pool
}

type msgCacheInfo struct {
	maxbody int
	pool    *sync.Pool
}

func newMsg(sz int) *Message {
	m := &Message{}
	m.Header = make([]byte, 0, 4)
	m.Body = make([]byte, 0, sz)
	m.bsize = sz
	return m
}

// 按size切割多个Message的池
var messageCache = []msgCacheInfo{
	{
		maxbody: 64,
		pool: &sync.Pool{
			New: func() interface{} { return newMsg(64) },
		},
	}, {
		maxbody: 128,
		pool: &sync.Pool{
			New: func() interface{} { return newMsg(128) },
		},
	}, {
		maxbody: 256,
		pool: &sync.Pool{
			New: func() interface{} { return newMsg(256) },
		},
	}, {
		maxbody: 512,
		pool: &sync.Pool{
			New: func() interface{} { return newMsg(512) },
		},
	}, {
		maxbody: 1024,
		pool: &sync.Pool{
			New: func() interface{} { return newMsg(1024) },
		},
	}, {
		maxbody: 4096,
		pool: &sync.Pool{
			New: func() interface{} { return newMsg(4096) },
		},
	}, {
		maxbody: 8192,
		pool: &sync.Pool{
			New: func() interface{} { return newMsg(8192) },
		},
	}, {
		maxbody: 65536,
		pool: &sync.Pool{
			New: func() interface{} { return newMsg(65536) },
		},
	},
}

//Free 在message上减少引用计数，如果没有进一步的引用，则释放它的资源
//这样做可以让资源在不使用GC的情况下回收，这对性能有相当大的好处
func (m *Message) Free() {
	if v := atomic.AddInt32(&m.refcnt, -1); v > 0 {
		return
	}
	for i := range messageCache {
		if m.bsize == messageCache[i].maxbody {
			messageCache[i].pool.Put(m)
			return
		}
	}
}

//Dup 增加消息上的引用计数
func (m *Message) Dup() *Message {
	atomic.AddInt32(&m.refcnt, 1)
	return m
}

//NewMessage是获得新Message的方式
//使用池缓存，极大地减少了垃圾收集器的负载
func NewMessage(sz int) *Message {
	var m *Message
	for i := range messageCache {
		if sz < messageCache[i].maxbody {
			m = messageCache[i].pool.Get().(*Message)
			break
		}
	}
	if m == nil {
		m = newMsg(sz)
	}

	m.refcnt = 1
	return m
}
