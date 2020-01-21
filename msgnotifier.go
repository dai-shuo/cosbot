package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"reflect"
	"sync"
)

type MessageCallbackFunc func(chatId, userId, msgId, text string)

type MessageNotifier struct {
	sync.RWMutex
	engine *gin.Engine
	addr []string
	callbacks []MessageCallbackFunc
}

func NewMessageNotifier(addr...string) *MessageNotifier {
	n := &MessageNotifier{
		engine: gin.Default(),
		addr: addr,

	}
	n.engine.POST("/lark", n.handle)
	return n
}

func (n *MessageNotifier) Run() error {
	return n.engine.Run(n.addr...)
}

func (n *MessageNotifier) handle(c *gin.Context) {
	d := make(map[string]interface{})

	if c.BindJSON(&d) == nil {
		switch d["type"] {
		case "url_verification":
			c.JSON(http.StatusOK, gin.H{
				"challenge": d["challenge"],
			})
		case "event_callback":
			if e, ok := d["event"].(map[string]interface{}); ok {
				if e["type"] == "message" || e["msg_type"] == "text" {
					n.notifyMessageEvent(e)
				}
			}
			c.String(http.StatusOK, ":)")
		default:
			c.String(http.StatusOK, ":)")
		}
	}
}

func (n *MessageNotifier) notifyMessageEvent(e map[string]interface{}) {
	defer func() {
		recover()
	}()
	chat := e["open_chat_id"].(string)
	user := e["open_id"].(string)
	msg := e["open_message_id"].(string)
	text := e["text_without_at_bot"].(string)
	n.RLock()
	defer n.RUnlock()
	for _, c := range n.callbacks {
		go c(chat, user, msg, text)
	}
}

func (n *MessageNotifier) AddCallback(f MessageCallbackFunc) {
	n.Lock()
	defer n.Unlock()
	n.callbacks = append(n.callbacks, f)
}

func (n *MessageNotifier) RemoveCallback(f MessageCallbackFunc) {
	n.Lock()
	defer n.Unlock()
	newCallbacks := make([]MessageCallbackFunc, 0, len(n.callbacks))
	fPtr := reflect.ValueOf(f).Pointer()
	for _, c := range n.callbacks {
		if reflect.ValueOf(c).Pointer() == fPtr {
			continue
		}
		newCallbacks = append(newCallbacks, c)
	}
	n.callbacks = newCallbacks
}

func (n *MessageNotifier) RemoveAllCallbacks() {
	n.Lock()
	defer n.Unlock()
	n.callbacks = nil
}
