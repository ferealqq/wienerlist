package main

import (
	"github.com/ferealqq/wienerlist/front/components"
	"github.com/hexops/vecty"
)

func main() {
	// Move bootstrap to custom wasmserve so that you can also include bootstrap javascript etc
	vecty.AddStylesheet("https://cdn.jsdelivr.net/npm/bootstrap@5.2.0-beta1/dist/css/bootstrap.min.css")
	vecty.AddStylesheet("app.css")

	vecty.SetTitle("Wienerlist • Wiener boards!")
	p := &components.PageView{}
	vecty.RenderBody(p)
}
