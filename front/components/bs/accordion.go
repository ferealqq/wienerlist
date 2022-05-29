package bs

import (
	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
)

func Accordion(markup ...vecty.MarkupOrChild) *vecty.HTML {
	markup = append(markup,
		vecty.Markup(
			vecty.Class("accordion"),
			// FIXME this is could be necessary for the accordion to work properly
			// vecty.Attribute("id", id),
		))
	//
	return elem.Div(markup...)
}

func AccordionItem(markup ...vecty.MarkupOrChild) *vecty.HTML {
	markup = append(markup, vecty.Markup(vecty.Class("accordion-item")))
	return elem.Div(markup...)
}

func AccordionHeader(markup ...vecty.MarkupOrChild) *vecty.HTML {
	// TODO maybe implement id?
	markup = append(markup, vecty.Markup(vecty.Class("accordion-header")))
	return elem.Heading2(markup...)
}

// TODO Create default values and more configurable options for setting args
func AccordionButton(target string, markup ...vecty.MarkupOrChild) *vecty.HTML {
	markup = append(markup,
		vecty.Markup(
			vecty.Class("accordion-button"),
			vecty.Class("collapsed"),
			vecty.Attribute("type", "button"),
			vecty.Attribute("data-bs-toggle", "collapse"),
			vecty.Attribute("data-bs-target", "#"+target),
			vecty.Attribute("aria-expanded", "true"),
			vecty.Attribute("aria-controls", target),
		))
	return elem.Button(markup...)
}

func AccordionBodyWrp(id string, markup ...vecty.MarkupOrChild) *vecty.HTML {
	markup = append(markup,
		vecty.Markup(
			vecty.Class("accordion-collapse"),
			vecty.Class("collapse"),
			vecty.Attribute("id", id),
			vecty.Attribute("aria-labelledby", "headingOne"),
			// FIXME this is maybe necessary to get accordion to work vecty.Attribute("data-bs-parent", "#"+parent),
		),
	)

	return elem.Div(markup...)
}

func AccordionBodyContent(markup ...vecty.MarkupOrChild) *vecty.HTML {
	markup = append(markup,
		vecty.Markup(
			vecty.Class("acordion-body"),
		),
	)

	return elem.Div(markup...)
}

func AccordionBody(id string, markup ...vecty.MarkupOrChild) *vecty.HTML {
	return AccordionBodyWrp(id, AccordionBodyContent(markup...))
}
