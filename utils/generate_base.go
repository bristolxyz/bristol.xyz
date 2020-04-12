package utils

import (
	"bytes"
	"github.com/bristolxyz/bristol.xyz/models"
	"github.com/getsentry/sentry-go"
	"html/template"
	"io/ioutil"
)

var baseTemplate *template.Template

func init() {
	f, err := ioutil.ReadFile("./templates/base.html")
	if err != nil {
		sentry.CaptureException(err)
		panic(err)
	}
	t, err := template.New("base").Parse(string(f))
	if err != nil {
		sentry.CaptureException(err)
		panic(err)
	}
	baseTemplate = t
}

// GenerateBase is used to generate the base HTML around the provided content.
func GenerateBase(Title, Description, AdditionalHeadHTML, ContentHTML string, User *models.User) string {
	b := bytes.Buffer{}
	err := baseTemplate.Execute(&b, map[string]interface{}{
		"Title": Title,
		"Description": Description,
		"AdditionalHeadHTML": template.HTML(AdditionalHeadHTML),
		"ContentHTML": template.HTML(ContentHTML),
		"User": User,
	})
	if err != nil {
		sentry.CaptureException(err)
		panic(err)
	}
	return b.String()
}
