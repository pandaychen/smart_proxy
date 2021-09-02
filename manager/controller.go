package manager

type Controller struct {
	BindAddr string
}

// NewController 创建一个Controller
func NewController() *Controller {
	return &Controller{}
}
