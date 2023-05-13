package wechat

import (
	"github.com/eatmoreapple/openwechat"
	log "github.com/sirupsen/logrus"
)

type MessageHandlerInterface interface {
	handle(*openwechat.Message) error
	ReplyText(*openwechat.Message) error
}

type Type string

const (
	GroupHandler = "group"
	UserHandler  = "user"
)

var handlers map[Type]MessageHandlerInterface

func init() {
	handlers = make(map[Type]MessageHandlerInterface)
	handlers[GroupHandler] = NewGroupMessageHandler()
	handlers[UserHandler] = NewUserMessageHandler()
}

func Handler(msg *openwechat.Message) {
	if msg.IsSendByGroup() {
		err := handlers[GroupHandler].handle(msg)
		if err != nil {
				log.Errorf("GroupHandler handle error: %s\n", err.Error())
				return
		}
	}else if msg.IsSendByFriend() {
		err := handlers[UserHandler].handle(msg)
		if err != nil {
				log.Errorf("FriendHandler handle error: %s\n", err.Error())
				return
		}
	}
}
