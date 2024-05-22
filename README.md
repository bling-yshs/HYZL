<p align="center"><img src="https://cdn.jsdelivr.net/gh/bling-yshs/ys-image-host@main/img/hyzl-icon.jpg" width="300" alt="icon" /></p>
<p align="center"><b>基于 Golang 的 <a href="https://github.com/yoimiya-kokomi/Miao-Yunzai">云崽</a> 启动器</b></p>
<p align="center">
  <a href="https://www.gnu.org/licenses/gpl-3.0.html"><img src="https://shields.io/github/license/bling-yshs/YzLauncher-windows?color=%231890FF" alt="License"></a>
  <a href="https://gitee.com/bling_yshs/Yunzai-v3-Installation-Steps"><img src="https://gitee.com/bling_yshs/Yunzai-v3-Installation-Steps/badge/star.svg?theme=dark" alt="Stars"></a>
  <a href="https://github.com/badges/shields/pulse"><img src="https://img.shields.io/github/commit-activity/m/bling-yshs/YzLauncher-windows" alt="Activity"/></a>
  <a href="https://github.com/bling-yshs/YzLauncher-windows/releases"><img src="https://img.shields.io/github/v/release/bling-yshs/YzLauncher-windows" alt="GitHub release"></a>
</p>

------------------------------

## 关于名字

HYZL 全名为 Hello Yunzai Launcher

由于一开始本项目并没有打算长期更新， 所以最早并没有取名。

后来签名API被封杀，本项目也停摆。着手开发于适配 QQNT 的，有UI界面的重制版 [HYZL](https://github.com/bling-yshs/HYZL-Tauri) 。

但是后来发现 Ws-Plugin 所适配的程度有限，QQNT的云崽用起来并不舒服。并且考虑到大部分云崽都是部署在服务器上，内存并不富裕，所以终止了 HYZL 的开发。

但是如今签名API复活，故本项目重新启用，为了纪念 HYZL ，本项目最终决定也叫 **HYZL**

## 快速开始

### 系统要求

HYZL 仅支持 Windows (Windows 10+，Windows Server 2016+) ，Windows 7，Windows Server 2012 可能会有存在兼容性问题

### 下载

从 [Github release](https://github.com/bling-yshs/YzLauncher-windows/releases/latest) 或 [Gitee release](https://gitee.com/bling_yshs/Yunzai-v3-Installation-Steps/releases/latest) 下载最新版本的 HYZL

### 依赖

HYZL 基于 Golang ，所以无需任何额外运行库。但是云崽本体需要 Git 和 Node.js ，所以你需要提前安装好它们

### 使用

#### 从零开始

随便创建一个文件夹，将 HYZL 放入其中并运行。进入主界面后，选择 `安装云崽` ，完成以后选择 `云崽管理` -> `启动云崽` 即可

#### 接入现有云崽

将 HYZL 放在 `Yunzai-Bot` 或 `Miao-Yunzai` 文件夹的同级目录即可，如图所示

![接入现有云崽](https://cdn.jsdelivr.net/gh/bling-yshs/ys-image-host@main/img/接入现有云崽.png)

## 项目开发

### 语言

[Golang](https://go.dev/) 1.22.3

### IDE

[Goland](https://www.jetbrains.com/go/) 或 [VS Code](https://code.visualstudio.com/)

### 格式化

在 Pull Request 之前，请务必使用 gofmt 进行代码格式化
