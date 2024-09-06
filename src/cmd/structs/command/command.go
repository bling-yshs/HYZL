package command

type Parameters struct {
	// 运行程序的名称(或路径)
	Name string
	// 运行程序的工作目录
	WorkDir string
	// 运行程序的参数
	Args []string
	// 是否自动关闭窗口
	AutoCloseWindow bool
}
