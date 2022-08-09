package main

//#include <stdint.h>
//#include <stdlib.h>
//typedef struct  { void* message; int size; char* error; } FetchWebsiteReturn;
import "C"

import (
	"github.com/massalabs/thyra/pkg/node"
	"github.com/massalabs/thyra/pkg/onchain/website"
)

//export fetchWebsite
func fetchWebsite(address *C.char, filename *C.char) *C.FetchWebsiteReturn {
	// instanciate here, must be free on caller side
	output := (*C.FetchWebsiteReturn)(C.malloc(C.size_t(C.sizeof_FetchWebsiteReturn)))

	c := node.NewDefaultClient() // server shall be set by the caller

	res, err := website.Fetch(c, C.GoString(address), C.GoString(filename))
	if err != nil {
		output.error = C.CString(err.Error())
		return output
	}

	output.error = nil
	output.message = C.CBytes(res)
	output.size = C.int(len(res))

	return output
}

func main() {}
