package main

import (
	"github.com/ciela/playground_golang/lgtm_maker/config"

	"github.com/naoina/kocha"
)

func main() {
	if err := kocha.Run(config.AppConfig); err != nil {
		panic(err)
	}
}
