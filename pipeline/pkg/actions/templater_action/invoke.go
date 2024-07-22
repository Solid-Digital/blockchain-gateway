package templater_action

import (
	"bytes"
	"github.com/Masterminds/sprig"
	"github.com/unchain/pipeline/pkg/domain"
	"github.com/unchainio/pkg/errors"
	"text/template"
)

const (
	// Any constants used in the application
	InputVariables = "variables"
	InputTemplate  = "template"
	TemplateResult = "result"
)

func Invoke(stub domain.Stub, input map[string]interface{}) (output map[string]interface{}, err error) {
	var buf bytes.Buffer

	msg, err := NewMessage(input)
	if err != nil {
		return nil, err
	}

	templ, err := template.New("template").Funcs(getFuncMap()).Parse(msg.Template)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse template")
	}

	err = templ.Execute(&buf, msg.Variables)
	if err != nil {
		return nil, errors.Wrap(err, "failed to render template")
	}

	return map[string]interface{}{TemplateResult: buf.String()}, nil
}

func getFuncMap() template.FuncMap {
	return sprig.TxtFuncMap()
}
