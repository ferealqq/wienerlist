package bs

import (
	v "github.com/hexops/vecty"
	e "github.com/hexops/vecty/elem"
)

func Modal(markup ...v.MarkupOrChild) *v.HTML {
	markup = append(markup,
		v.Markup(
			// v.Class("modal")))
			v.Class("modal"),
			v.Attribute("tabindex", "-1")))

	return e.Div(markup...)
}

func ModalDialog(markup ...v.MarkupOrChild) *v.HTML {
	markup = append(markup, v.Markup(v.Class("modal-dialog")))
	return e.Div(markup...)
}

func ModalContent(markup ...v.MarkupOrChild) *v.HTML {
	markup = append(markup, v.Markup(v.Class("modal-content")))
	return e.Div(markup...)
}

func ModalHeader(markup ...v.MarkupOrChild) *v.HTML {
	markup = append(markup, v.Markup(v.Class("modal-header")))
	return e.Div(markup...)
}

func ModalBody(markup ...v.MarkupOrChild) *v.HTML {
	markup = append(markup, v.Markup(v.Class("modal-body")))
	return e.Div(markup...)
}

func ModalFooter(markup ...v.MarkupOrChild) *v.HTML {
	markup = append(markup, v.Markup(v.Class("modal-footer")))
	return e.Div(markup...)
}

func FModal(header []v.MarkupOrChild, body []v.MarkupOrChild, footer []v.MarkupOrChild, modalConf ...v.MarkupOrChild) *v.HTML {
	modalConf = append(modalConf, ModalDialog(
		ModalContent(
			ModalHeader(header...),
			ModalBody(body...),
			ModalFooter(footer...))))
	return Modal(modalConf...)
}

// <div class="modal" tabindex="-1">
//   <div class="modal-dialog">
//     <div class="modal-content">
//       <div class="modal-header">
//         <h5 class="modal-title">Modal title</h5>
//         <button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
//       </div>
//       <div class="modal-body">
//         <p>Modal body text goes here.</p>
//       </div>
//       <div class="modal-footer">
//         <button type="button" class="btn btn-secondary" data-bs-dismiss="modal">Close</button>
//         <button type="button" class="btn btn-primary">Save changes</button>
//       </div>
//     </div>
//   </div>
// </div>
