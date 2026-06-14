package main

import (
	"github.com/nikdotcodes/gator/internal/config"
)

func main() {
	c, _ := config.Read()
	c.SetUser("nik")
}
