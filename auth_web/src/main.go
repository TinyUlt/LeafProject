package main
import (
	"flag"

	"LeafProject/sky/web"

)


func LoginHandler(ctx *web.Context) string {

	return "{\"error\":\"0\"}"
}
func main() {

	flag.Parse()

	web.Config.StaticDir = "./static"

	web.Post("/login", LoginHandler)

	web.Run("127.0.0.1:8080")

}
