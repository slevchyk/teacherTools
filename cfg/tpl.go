package cfg

import "html/template"


var Tpl *template.Template


func init() {

	Tpl = template.Must(template.ParseGlob("templates/*.gohtml"))
}