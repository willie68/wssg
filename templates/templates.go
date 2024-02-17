package templates

import _ "embed"

// Page the page template
var (
	//go:embed layout.html
	LayoutHTML string
	//go:embed reload.js
	AutoreloadJS string
	//go:embed cookiebanner.html
	Cookiebanner     string
	CookiebannerText = "Um Ihnen das beste Online-Erlebnis zu bieten, verwendet diese Website Cookies. Mit der Nutzung dieser Webseite erklären Sie sich damit einverstanden, dass wir Cookies verwenden, wie in unserer Datenschutzerklärung beschrieben."
)
