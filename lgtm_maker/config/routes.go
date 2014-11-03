package config

import (
	"github.com/ciela/playground_golang/lgtm_maker/app/controller"
	"github.com/naoina/kocha"
)

type RouteTable kocha.RouteTable

var routes = RouteTable{
	{
		Name:       "root",
		Path:       "/",
		Controller: &controller.Root{},
	}, {
		Name:       "images",
		Path:       "/images",
		Controller: &controller.Images{},
	}, {
		Name:       "imagesId",
		Path:       "/images/:imageId",
		Controller: &controller.Images{},
	},
}

func init() {
	AppConfig.RouteTable = kocha.RouteTable(append(routes, RouteTable{
		{
			Name:       "static",
			Path:       "/*path",
			Controller: &kocha.StaticServe{},
		},
	}...))
}
