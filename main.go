package idcard

import (
	"log"
	"strconv"
	"time"
)

//card 身份证号
//sex 性别
//address 归属地
//ok 验证是否成功
func IdCard(card string) (sex, address string, ok bool) {
	if len(card) == 18 {
		idcardJudge := idcardJudge(card)
		if idcardJudge {
			address = queryAddress(card)
			sex = querySex(card)
			return sex, address, true
		} else {
			return "", "", false
		}
	} else if len(card) == 15 {
		//15位身份证号码转18位
		card = Citizen15To18(card)
		idcardJudge := idcardJudge(card)
		if idcardJudge {
			address = queryAddress(card)
			sex = querySex(card)
			return sex, address, true
		} else {
			return "", "", false
		}
	}
	return "", "", false
}

//通过身份证查询归属地
func queryAddress(str string) string {
	str = str[:6]
	var (
		city string
		ok   bool
	)
	if city, ok = address[str]; !ok {
		city = "未更新或已删减该六位行政代码或该六位行政代码错误！"
	}
	return city
}

//通过身份证判断性别
func querySex(str string) string {
	//截取倒数第二位，奇数表示男性偶数表示女性
	str = str[len(str)-2 : len(str)-1]
	strInt, err := strconv.Atoi(str)
	if err != nil {
		log.Fatal(err)
	}
	if strInt%2 == 0 {
		return "女"
	} else {
		return "男"
	}
}

//合法身份证判断
var weight = [17]int{7, 9, 10, 5, 8, 4, 2, 1, 6, 3, 7, 9, 10, 5, 8, 4, 2}
var validValue = [11]byte{'1', '0', 'X', '9', '8', '7', '6', '5', '4', '3', '2'}
var validProvince = []string{
	"11", // 北京市
	"12", // 天津市
	"13", // 河北省
	"14", // 山西省
	"15", // 内蒙古自治区
	"21", // 辽宁省
	"22", // 吉林省
	"23", // 黑龙江省
	"31", // 上海市
	"32", // 江苏省
	"33", // 浙江省
	"34", // 安徽省
	"35", // 福建省
	"36", // 山西省
	"37", // 山东省
	"41", // 河南省
	"42", // 湖北省
	"43", // 湖南省
	"44", // 广东省
	"45", // 广西壮族自治区
	"46", // 海南省
	"50", // 重庆市
	"51", // 四川省
	"52", // 贵州省
	"53", // 云南省
	"54", // 西藏自治区
	"61", // 陕西省
	"62", // 甘肃省
	"63", // 青海省
	"64", // 宁夏回族自治区
	"65", // 新疆维吾尔自治区
	"71", // 台湾省
	"81", // 香港特别行政区
	"82", // 澳门特别行政区
}

func idcardJudge(str string) bool {
	//出生年
	nYear, _ := strconv.Atoi(str[6:10])
	//出生月
	nMonth, _ := strconv.Atoi(str[10:12])
	//出生日
	nDay, _ := strconv.Atoi(str[12:14])
	//验证区域码
	judge1 := CheckProvinceValid(str)
	//验证校验码
	judge2 := IsValidCitizenNo18(str)
	//验证生日
	judge3 := CheckBirthdayValid(nYear, nMonth, nDay)
	return judge1 && judge2 && judge3
}

//15位身份证转为18位
func Citizen15To18(citizenNo15 string) string {
	nLen := len(citizenNo15)
	if nLen != 15 {
		return "身份证不是15位！"
	}
	citizenNo18 := make([]byte, 0)
	citizenNo18 = append(citizenNo18, citizenNo15[:6]...)
	citizenNo18 = append(citizenNo18, '1', '9')
	citizenNo18 = append(citizenNo18, citizenNo15[6:]...)

	sum := 0
	for i, v := range citizenNo18 {
		n, _ := strconv.Atoi(string(v))
		sum += n * weight[i]
	}
	mod := sum % 11
	citizenNo18 = append(citizenNo18, validValue[mod])
	return string(citizenNo18)
}

//18位身份证校验码
func IsValidCitizenNo18(idcard string) bool {
	//string -> []byte
	citizenNo18 := []byte(idcard)
	nSum := 0
	for i := 0; i < len(citizenNo18)-1; i++ {
		n, _ := strconv.Atoi(string(citizenNo18[i]))
		nSum += n * weight[i]
	}
	//mod得出18位身份证校验码
	mod := nSum % 11
	if validValue[mod] == citizenNo18[17] {
		return true
	}
	return false
}

//出生年为闰年时2月有29号
func IsLeapYear(nYear int) bool {
	if nYear <= 0 {
		return false
	}
	if (nYear%4 == 0 && nYear%100 != 0) || nYear%400 == 0 {
		return true
	}
	return false
}

//验证生日
func CheckBirthdayValid(nYear, nMonth, nDay int) bool {
	if nYear < 1900 || nMonth <= 0 || nMonth > 12 || nDay <= 0 || nDay > 31 {
		return false
	}
	//出生日期大于现在的日期
	curYear, curMonth, curDay := time.Now().Date()
	if nYear == curYear {
		if nMonth > int(curMonth) {
			return false
		} else if nMonth == int(curMonth) && nDay > curDay {
			return false
		}
	}
	//出生日期在2月份
	if 2 == nMonth {
		//闰年2月只有29号
		if IsLeapYear(nYear) && nDay > 29 {
			return false
		} else if nDay > 28 { //非闰年2月只有28号
			return false
		}
	} else if 4 == nMonth || 6 == nMonth || 9 == nMonth || 11 == nMonth { //小月只有30号
		if nDay > 30 {
			return false
		}
	}
	return true
}

//验证区域码
func CheckProvinceValid(str string) bool {
	citizenNo := []byte(str)
	provinceCode := make([]byte, 0)
	provinceCode = append(provinceCode, citizenNo[:2]...)
	provinceStr := string(provinceCode)
	for i, _ := range validProvince {
		if provinceStr == validProvince[i] {
			return true
		}
	}
	return false
}
