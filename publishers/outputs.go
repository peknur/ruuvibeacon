package publishers

import "github.com/peknur/ruuvibeacon"

// Add publisher
func Add(name string, p ruuvibeacon.Publisher) {
	ruuvibeacon.Publishers[name] = p
}
