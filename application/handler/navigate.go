package handler

import (
	"github.com/coscms/webcore/library/navigate"
	"github.com/webx-top/echo"
)

var LeftNavigate = &navigate.Item{
	Display: true,
	Name:    echo.T(`网站管理`),
	Action:  `caddy`,
	Icon:    `sitemap`,
	Children: &navigate.List{
		{
			Display: false,
			Name:    echo.T(`Caddy日志`),
			Action:  `log_show`,
		},
		{
			Display: true,
			Name:    echo.T(`网站列表`),
			Action:  `vhost`,
		},
		{
			Display: true,
			Name:    echo.T(`添加网站`),
			Action:  `vhost_add`,
			Icon:    `plus`,
		},
		{
			Display: false,
			Name:    echo.T(`重启Caddy`),
			Action:  `restart`,
			Icon:    ``,
		},
		{
			Display: false,
			Name:    echo.T(`停止Caddy`),
			Action:  `stop`,
			Icon:    ``,
		},
		{
			Display: false,
			Name:    echo.T(`查看网站日志`),
			Action:  `vhost_log`,
		},
		{
			Display: false,
			Name:    echo.T(`查看网站动态`),
			Action:  `log`,
		},
		{
			Display: false,
			Name:    echo.T(`配置表单`),
			Action:  `addon_form`,
			Icon:    ``,
		},
		{
			Display: false,
			Name:    echo.T(`修改网站`),
			Action:  `vhost_edit`,
			Icon:    ``,
		},
		{
			Display: false,
			Name:    echo.T(`删除网站`),
			Action:  `vhost_delete`,
			Icon:    ``,
		},
		{
			Display: false,
			Name:    echo.T(`管理网站文件`),
			Action:  `vhost_file`,
			Icon:    ``,
		},
		{
			Display: false,
			Name:    echo.T(`生成Caddyfile`),
			Action:  `vhost_build`,
			Icon:    ``,
		},

		{
			Display: true,
			Name:    echo.T(`分组管理`),
			Action:  `group`,
		},
		{
			Display: true,
			Name:    echo.T(`添加分组`),
			Action:  `group_add`,
			Icon:    `plus`,
		},
		{
			Display: false,
			Name:    echo.T(`修改分组`),
			Action:  `group_edit`,
			Icon:    ``,
		},
		{
			Display: false,
			Name:    echo.T(`删除分组`),
			Action:  `group_delete`,
			Icon:    ``,
		},

		{
			Display: true,
			Name:    echo.T(`引擎配置`),
			Action:  `server`,
		},
		{
			Display: false,
			Name:    echo.T(`添加引擎配置`),
			Action:  `server_add`,
			Icon:    `plus`,
		},
		{
			Display: false,
			Name:    echo.T(`修改引擎配置`),
			Action:  `server_edit`,
			Icon:    ``,
		},
		{
			Display: false,
			Name:    echo.T(`删除引擎配置`),
			Action:  `server_delete`,
			Icon:    ``,
		},
		{
			Display: false,
			Name:    echo.T(`更新HTTPS证书`),
			Action:  `server_renew_cert`,
			Icon:    ``,
		},
	},
}
