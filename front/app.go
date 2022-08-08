package main

import (
	"syscall/js"

	"github.com/ferealqq/wienerlist/front/components"
	"github.com/hexops/vecty"
)

func main() {
	// Move bootstrap to custom wasmserve so that you can also include bootstrap javascript etc
	vecty.AddStylesheet("https://cdn.jsdelivr.net/npm/bootstrap@5.2.0-beta1/dist/css/bootstrap.min.css")
	vecty.AddStylesheet("https://cdn.jsdelivr.net/npm/bootstrap-icons@1.8.3/font/bootstrap-icons.css")
	AddScript("https://cdn.jsdelivr.net/npm/bootstrap@5.0.2/dist/js/bootstrap.min.js")
	vecty.AddStylesheet("app.css")

	vecty.SetTitle("Wienerlist • Wiener boards!")
	p := &components.PageView{}
	vecty.RenderBody(p)
}

func AddScript(url string) {
	script := js.Global().Get("document").Call("createElement", "script")
	script.Set("src", url)
	js.Global().Get("document").Get("head").Call("appendChild", script)
}
