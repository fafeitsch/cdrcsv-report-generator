package report

import (
	"flag"
	"fmt"
)

func main() {
	definition := flag.String("definition", "./definition.json", "Path to the json file containing the report definition.")
	template := flag.String("template", "./template.gohtml", "Path to the html file containing the template.")
	flag.Parse()
	_ = fmt.Sprintf("%v%v", definition, template)
}
