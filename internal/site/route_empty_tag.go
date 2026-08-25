package site

func EmptyTagPage(tag string, preview bool) interfaceComponent {
	return emptyTagPage(tag, PageData{Preview: preview}, nil)
}
