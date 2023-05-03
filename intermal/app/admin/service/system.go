package service

type ISystem interface {
	Login()
	Logout()
	Captcha()
}

type System struct{}

func (s System) Login() {
	//TODO implement me
	panic("implement me")
}

func (s System) Logout() {
	//TODO implement me
	panic("implement me")
}

func (s System) Captcha() {
	//TODO implement me
	panic("implement me")
}