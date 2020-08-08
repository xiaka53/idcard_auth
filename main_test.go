package idcard

import "testing"

func Test_IdCard(t *testing.T) {
	sex, address, ok := IdCard("")
	if !ok {
		t.Log("验证失败")
	}
	if sex == "男" || sex == "女" {
		t.Log("性别验成功：", sex)
	} else {
		t.Error("性别验证失败")
	}
	t.Log("归属地验证结果：", address)

}
