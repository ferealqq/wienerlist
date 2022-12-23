package components

import (
	"fmt"

	"github.com/ferealqq/wienerlist/front/components/bs"
	i "github.com/ferealqq/wienerlist/front/components/input"
	u "github.com/ferealqq/wienerlist/front/components/util"
	"github.com/ferealqq/wienerlist/front/util"
	v "github.com/hexops/vecty"
	e "github.com/hexops/vecty/elem"
	evt "github.com/hexops/vecty/event"
)

type ItemModal struct {
	v.Core
}

func (b *ItemModal) Render() v.ComponentOrHTML {
	id, err := util.GetIntVar("id")
	if err != nil {
		return v.Text("Invalid board id")
	}
	// itemId, err := util.GetIntVar("itemId")
	// if err != nil {
	// 	return v.Text("Invalid item id")
	// }
	return bs.FModal(
		[]v.MarkupOrChild{e.Heading3(v.Text("Item"))},
		[]v.MarkupOrChild{
			i.HoverTextInput(
				v.Markup(
					v.Attribute("value", "Otsikko"),
					v.Attribute("id", "title"),
				),
			),
			i.HoverTextAreaInput(
				v.Markup(
					v.Attribute("id", "description"),
				),
				v.Text("leipätekstiä komponenttiin"),
			),
		},
		[]v.MarkupOrChild{
			e.Div(
				u.Classes("justify-content-between", "d-flex", "w-100"),
				bs.Button2ry(v.Text("Close"), v.Markup(
					evt.Click(func(e *v.Event) {
						util.Redirect(fmt.Sprintf("/boards/%d", id))
					}),
				)),
				bs.ButtonPry(v.Text("Save changes")),
			),
		},
		v.Markup(
			v.Class("fade"),
			v.Class("show"),
			v.Class("item-modal"),
			v.Class("bg-secondary"),
			v.Class("bg-opacity-25"),
		),
	)
}
