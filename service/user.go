package service

import (
	dto "GinProject/dto/user"
	"GinProject/model"
	"GinProject/query"
	"GinProject/utils"
	"context"
	"github.com/go-gomail/gomail"
	"log"
	"math/rand"
	"time"
)

func UserLogin(loginKey string, password string) bool {
	user1 := getByTelephone(loginKey)
	user2 := getByEmail(loginKey)
	if user1 != nil {
		return *user1.Password == password
	} else if user2 != nil {
		return *user2.Password == password
	} else {
		return false
	}
}

func getByTelephone(tel string) *model.User {
	dUser := query.Q.User
	ctx := context.Background()
	user1, err := dUser.WithContext(ctx).Where(dUser.Telephone.Eq(tel)).First()
	if err != nil {
		return nil
	} else {
		return user1
	}
}

func getByEmail(email string) *model.User {
	dUser := query.Q.User
	ctx := context.Background()
	user1, err := dUser.WithContext(ctx).Where(dUser.Email.Eq(email)).First()
	if err != nil {
		return nil
	} else {
		return user1
	}
}

func randomCode() string {
	// 设置随机数种子
	rand.New(rand.NewSource(time.Now().UnixNano()))

	// 定义字符集合
	charset := "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	// 生成6位随机字符串
	randomString := make([]byte, 6)
	for i := range randomString {
		randomString[i] = charset[rand.Intn(len(charset))]
	}

	return string(randomString)
}

func SendEmail(email string) bool {
	// 设置发件人信息
	from := "2212340514@qq.com"
	password := "wnqdhkuazojndjae"

	// 设置收件人信息
	to := email

	//生成6位验证码
	registerCode := randomCode()

	//存入redis
	err := utils.Client.Set(email, registerCode, time.Minute*10).Err()
	if err != nil {
		log.Print("write redis fail")
		return false
	}
	// 创建邮件内容
	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", "Hello, This Is A New Register Mail!")
	m.SetBody("text/plain", "Your Register Code is : "+registerCode)

	// 创建 SMTP 发送器
	d := gomail.NewDialer("smtp.qq.com", 587, from, password)

	// 发送邮件
	err = d.DialAndSend(m)
	if err != nil {
		return false
	}
	return true
}

func Register(form *dto.RegisterForm) bool {
	code, err := utils.Client.Get(form.Email).Result()
	if err != nil {
		return false
	}
	if code != form.RegisterCode {
		return false
	} else {
		dUser := query.User
		ctx := context.Background()
		user := model.User{
			Email:     &form.Email,
			Telephone: &form.Telephone,
			Password:  &form.Password,
		}
		insertErr := dUser.WithContext(ctx).Create(&user)
		if insertErr != nil {
			return false
		} else {
			utils.Client.Del(form.Email)
			return true
		}
	}
}

func UserLogout(userKey string) bool {
	err := utils.Client.Del("login:jwt:" + userKey).Err()
	if err != nil {
		return false
	} else {
		return true
	}
}

func UserSignIn(userId string) bool {
	//用户签到实现
	today := int64(time.Now().Day())
	err := utils.Client.SetBit("sign:user:"+userId, today, 1).Err()
	if err != nil {
		return false
	}
	return true
}

func GetUserSign(userId string, day int64) int64 {
	res, err := utils.Client.GetBit("sign:user:"+userId, day).Result()
	if err != nil {
		return 0
	}
	return res
}

func GetUserById(id int64) *model.User {
	U := query.User
	ctx := context.Background()
	user, err := U.WithContext(ctx).Where(U.UserID.Eq(id)).First()
	if err != nil {
		return nil
	}
	return user
}

func DeleteUser(id int64) bool {
	U := query.User
	ctx := context.Background()
	_, err := U.WithContext(ctx).Where(U.UserID.Eq(id)).Delete()
	return err == nil
}

func ChangeUserStatus(id int64) bool {
	U := query.User
	ctx := context.Background()
	user, err := U.WithContext(ctx).Where(U.UserID.Eq(id)).Select(U.State).First()
	if err != nil {
		return false
	}
	if *user.State == "正常" {
		_, err = U.WithContext(ctx).Where(U.UserID.Eq(id)).Update(U.State, "封禁")
		return err == nil
	} else {
		_, err = U.WithContext(ctx).Where(U.UserID.Eq(id)).Update(U.State, "正常")
		return err == nil
	}
}

func UserRealName(real *dto.RealName) bool {
	U := query.User
	ctx := context.Background()
	isReal := int64(1)
	user := model.User{
		IsReal:   &isReal,
		RealName: &real.Name,
		RealID:   &real.RealId,
	}
	_, err := U.WithContext(ctx).Where(U.UserID.Eq(real.Id)).Updates(user)
	return err == nil
}

func UserAdminAuth(id int64) bool {
	U := query.User
	ctx := context.Background()
	_, err := U.WithContext(ctx).Where(U.UserID.Eq(id)).Update(U.IsAdmin, 1)
	return err == nil
}

func UpdateUserInfo(form *dto.UserUpdateForm) bool {
	ctx := context.Background()
	U := query.User
	_, err := U.WithContext(ctx).Where(U.UserID.Eq(form.Id)).Updates(form)
	return err == nil
}

func UpdateUserPsw(form *dto.UserPswForm) bool {
	ctx := context.Background()
	U := query.User
	user, err := U.WithContext(ctx).Where(U.UserID.Eq(form.Id)).First()
	if err != nil {
		return false
	}
	if *user.Password != form.OldPassword {
		return false
	}
	_, err = U.WithContext(ctx).Where(U.UserID.Eq(form.Id)).Update(U.Password, form.NewPassword)
	return err == nil
}
