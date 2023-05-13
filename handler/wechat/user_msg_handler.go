package wechat

import (
	"strings"

	"wechatbot/config"
	"wechatbot/openai"
	"wechatbot/utils"

	"github.com/eatmoreapple/openwechat"
	log "github.com/sirupsen/logrus"
)

var _ MessageHandlerInterface = (*UserMessageHandler)(nil)

// UserMessageHandler 私聊消息处理
type UserMessageHandler struct {
}

// handle 处理消息
func (g *UserMessageHandler) handle(msg *openwechat.Message) error {
	if msg.IsText() {
		return g.ReplyText(msg)
	}
	return nil
}

// NewUserMessageHandler 创建私聊处理器
func NewUserMessageHandler() MessageHandlerInterface {
	return &UserMessageHandler{}
}

// ReplyText
func (g *UserMessageHandler) ReplyText(msg *openwechat.Message) error {
	// 接收私聊消息
	sender, err := msg.Sender()
	log.Printf("Received User %v Text Msg : %v", sender.NickName, msg.Content)
	// if UserService.ClearUserSessionContext(sender.ID(), msg.Content) {
	// 	_, err = msg.ReplyText("上下文已经清空了，你可以问下一个问题啦。")
	// 	if err != nil {
	// 		log.Printf("response user error: %v \n", err)
	// 	}
	// 	return nil
	// }

	wechat := config.GetWechatKeyword()
	requestText := msg.Content
	if wechat != nil {
		content, key := utils.ContainsI(requestText, *wechat)
		if len(key) == 0 {
			return nil
		}

		splitItems := strings.Split(content, key)
		if len(splitItems) < 2 {
			return nil
		}

		requestText = strings.TrimSpace(splitItems[1])
	}

	// 获取上下文，向GPT发起请求
	// requestText := strings.TrimSpace(msg.Content)
	requestText = strings.Trim(msg.Content, "\n")

	// requestText = UserService.GetUserSessionContext(sender.ID()) + requestText
	reply, err := openai.Completions(requestText)
	if err != nil {
		log.Printf("request error: %v \n", err)
		msg.ReplyText("出问题了，我一会儿去修理一下。")
		return err
	}
	result := *reply

	// 设置上下文，回复用户
	result = strings.TrimSpace(result)
	result = strings.Trim(result, "\n")
	// UserService.SetUserSessionContext(sender.ID(), requestText, reply)
	_, err = msg.ReplyText(result)
	if err != nil {
		log.Printf("response user error: %v \n", err)
	}
	return err
}
