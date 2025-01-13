package main

import (
	"github.com/sirupsen/logrus"
	"github.com/zjyl1994/yusifubot/infra/startup"
)

func main() {
	if err := startup.Start(); err != nil {
		logrus.Fatalln(err.Error())
	}
}
