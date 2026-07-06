package asky // Имя пакета должно совпадать с именем вашего модуля из go.mod

import "embed"

//go:embed web/templates/*.html
var TemplateFS embed.FS
