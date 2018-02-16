package jupiter

import (
	srvModels "github.com/home-IoT/jupiter/server/models"
	"net/http"
)

// ServerCreateLinksWithSelf a links object with a self URL reference
func ServerCreateLinksWithSelf(request *http.Request) *srvModels.GenericLinks {
	links := new(srvModels.GenericLinks)

	selfLink := request.URL.RequestURI()

	//if !request.URL.IsAbs() {
	//	selfLink = request.Host + selfLink
	//}

	links.Self = &selfLink

	return links
}
