package main

import "fmt"

func main() {
	notifier := NewMessageNotifier(":8000")
	notifier.AddCallback(func(chatId, userId, msgId, text string) {
		fmt.Printf("chatId=%s, userId=%s, msgId=%s, text=%s\n", chatId, userId, msgId, text)
	})
	_ = notifier.Run()
}
