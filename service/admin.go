package service

import (
	dto "GinProject/dto/admin"
	"GinProject/model"
	"GinProject/query"
	"GinProject/utils"
	"context"
	"strconv"
)

func AdminLogin(loginForm *dto.AdminLoginForm) bool {
	admin := GetAdminById(loginForm.AdminId)
	if admin == nil {
		return false
	} else {
		return *admin.Password == loginForm.Password
	}
}

func GetAdminById(id string) *model.Admin {
	A := query.Admin
	ctx := context.Background()
	adminId, _ := strconv.Atoi(id)
	admin, err := A.WithContext(ctx).Where(A.AdminID.Eq(int64(adminId))).First()
	if err != nil {
		return nil
	} else {
		return admin
	}
}

func AdminLogout(id string) bool {
	err := utils.Client.Del("login:jwt:" + id).Err()
	if err != nil {
		return false
	} else {
		return true
	}
}

func UpdateAdminInfo(form *dto.AdminUpdateForm) bool {
	A := query.Admin
	ctx := context.Background()
	_, err := A.WithContext(ctx).Where(A.AdminID.Eq(form.Id)).Updates(form)
	return err == nil
}

func UpdateAdminPassword(form *dto.AdminPswForm) bool {
	A := query.Admin
	ctx := context.Background()
	admin, err := A.WithContext(ctx).Where(A.AdminID.Eq(form.Id)).First()
	if err != nil {
		return false
	}
	if *admin.Password != form.OldPassword {
		return false
	}
	_, err = A.WithContext(ctx).Where(A.AdminID.Eq(form.Id)).Update(A.Password, form.NewPassword)
	return err == nil
}
