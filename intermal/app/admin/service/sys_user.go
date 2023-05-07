package service

import (
	"strings"

	"github.com/gogoclouds/go-web/intermal/app/admin/dao"
	"github.com/gogoclouds/go-web/intermal/app/admin/enum"
	"github.com/gogoclouds/go-web/intermal/app/admin/model"
	"github.com/gogoclouds/gogo/g"
	"github.com/gogoclouds/gogo/pkg/util"
	"github.com/gogoclouds/gogo/web/r"
	"github.com/jinzhu/copier"
)

type IUser interface {
	PageList(req model.UserPageQuery) (*r.PageAnyResp, *g.Error)
	Details(id string) (model.SysUser, *g.Error)
	Create(req model.UserCreateReq) *g.Error
	Updates(req model.SysUser) *g.Error
	UpdateStatus(req model.UserUpdateStatusReq) *g.Error
	UpdatePassword(req model.UserUpdatePasswdReq) *g.Error
	Delete(ID string) *g.Error
}

var userDao dao.IUser = new(dao.User)

type User struct{}

func (svc *User) PageList(req model.UserPageQuery) (*r.PageAnyResp, *g.Error) {
	return userDao.PageList(req)
}

func (svc *User) Details(id string) (model.SysUser, *g.Error) {
	return userDao.FindByID(id)
}

func (svc *User) Create(req model.UserCreateReq) *g.Error {
	var user model.SysUser
	copier.Copy(&user, &req)

	user.Username = strings.Trim(user.Username, " ")
	user.Nickname = strings.Trim(user.Nickname, " ")
	user.Status = enum.UserStatusEnable

	hashPwd := svc.bcrypt().Hash(req.Password)
	user.Password = hashPwd

	return userDao.Create(user)
}

func (svc *User) Updates(req model.SysUser) *g.Error {
	return userDao.Update(req)
}

func (svc *User) UpdateStatus(req model.UserUpdateStatusReq) *g.Error {
	var user model.SysUser
	copier.Copy(&user, &req)
	return userDao.Update(user)
}

func (svc *User) UpdatePassword(req model.UserUpdatePasswdReq) *g.Error {
	u, gerr := userDao.FindByID(req.ID)
	if gerr != nil {
		return gerr
	}

	bcrypt := svc.bcrypt()
	if !bcrypt.Check(u.Password, req.OldPassword) {
		return g.NewError("旧密码验证不通过")
	}

	var user model.SysUser
	copier.Copy(&user, &req)
	user.Password = bcrypt.Hash(req.NewPassword)
	return userDao.Update(user)
}

func (svc *User) Delete(ID string) *g.Error {
	return userDao.Delete(ID)
}

func (svc *User) bcrypt() util.Bcrypt {
	serviceKV := g.Conf.AppServiceKV()
	bcrypt := util.NewBcrypt(serviceKV["salt"].(string))
	return bcrypt
}
