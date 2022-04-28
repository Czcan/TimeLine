package email

import (
	"strings"
	"sync"
	"testing"
)

func TestSendEmail(t *testing.T) {
	body := strings.Replace(`
		<html>
		<body>
		<h4>TimeLine</h3>
		<p>[TimeLine]时间轴手账</p>
		<p>$，请在三分钟内输入验证码。为了你的账号安全，请勿将验证码告知他人</p>
		</body>
		</html>
	`, "$", "123456", -1)
	var wg sync.WaitGroup
	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			email := New("smtp.qq.com:25", "1479039156@qq.com", "[TimeLine]时间轴手账", "dtoeddyfqbjsghgg")
			err := email.Send("1479039156@qq.com", body, "[TimeLine]时间轴手账验证码", 1)
			if err != nil {
				t.Errorf("send email failed error: %v", err.Error())
			}
		}()
	}
	wg.Wait()
}
