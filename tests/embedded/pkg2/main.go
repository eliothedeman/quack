package main

import (
	"log"

	"github.com/eliothedeman/quack"
	"github.com/eliothedeman/quack/tests/embedded/pkg1"
)

type Local struct {
	Yar string
}

type cmd struct {
	pkg1.SomeType
	Local
}

func (c cmd) Run([]string) {
	log.Println(c.Config)
	log.Println(c.Yar)

}

func main() {
	quack.Run("test", quack.WithCommand(new(cmd)))

}
