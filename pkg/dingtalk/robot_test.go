package dingtalk

import (
	"os"
	"testing"
)

func TestRobot_SendMessage(t *testing.T) {
	//t.SkipNow()

	msg := map[string]interface{}{
		"msgtype": "text",
		"text": map[string]string{
			"content": "这是一条golang钉钉消息测试.",
		},
		"at": map[string]interface{}{
			"atMobiles": []string{},
			"isAtAll":   false,
		},
	}

	robot := NewRobot(os.Getenv("ROBOT_TOKEN"), os.Getenv("ROBOT_SECRET"))
	if err := robot.SendMessage(msg); err != nil {
		t.Error(err)
	}
}

func TestRobot_SendTextMessage(t *testing.T) {
	robot := NewRobot(os.Getenv("ROBOT_TOKEN"), os.Getenv("ROBOT_SECRET"))
	if err := robot.SendTextMessage("普通文本消息", []string{}, false); err != nil {
		t.Error(err)
	}
}

func TestRobot_SendMarkdownMessage(t *testing.T) {
	robot := NewRobot(os.Getenv("ROBOT_TOKEN"), os.Getenv("ROBOT_SECRET"))
	err := robot.SendMarkdownMessage(
		"Markdown Test Title",
		"### Markdown 测试消息\n* 谷歌: [Google](https://www.google.com/)\n* 一张图片\n ![](https://avatars0.githubusercontent.com/u/40748346)",
		[]string{},
		false,
	)
	if err != nil {
		t.Error(err)
	}
}

func TestRobot_SendLinkMessage(t *testing.T) {
	robot := NewRobot(os.Getenv("ROBOT_TOKEN"), os.Getenv("ROBOT_SECRET"))
	err := robot.SendLinkMessage(
		"Link Test Title",
		"这是一条链接测试消息",
		"https://github.com/JetBlink",
		"https://avatars0.githubusercontent.com/u/40748346",
	)

	if err != nil {
		t.Error(err)
	}
}
