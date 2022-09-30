package gitcord

// import (
// 	"strings"
// 	"text/template"
// 	"time"

// 	"github.com/pkg/errors"
// )

// var (
// 	//go:embed pr_preview_message.tmpl.md
// 	prReviewMessage string
// 	//go:embed issue_message.tmpl.md
// 	issueMessage string
// )

// var (
// 	prReviewMessageTmpl = parseTmpl(prReviewMessage)
// 	issueMessageTmpl    = parseTmpl(issueMessage)
// )

// var funcMap = template.FuncMap{
// 	"now": time.Now,
// }

// func renderTmpl(tmpl *template.Template, v any) (string, error) {
// 	var s strings.Builder
// 	err := tmpl.Execute(&s, v)
// 	return s.String(), errors.Wrap(err, "cannot render template")
// }

// func parseTmpl(str string) *template.Template {
// 	return template.Must(template.New("").Funcs(funcMap).Parse(str))
// }
