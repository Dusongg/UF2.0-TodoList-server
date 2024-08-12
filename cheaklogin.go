package main

type onlineMap struct {
	omap map[string]bool
}

func NewOnlineMap() *onlineMap {
	return &onlineMap{make(map[string]bool)}
}

var Omap *onlineMap = NewOnlineMap()

func (o *onlineMap) checkIfLoggedIn(userName string) bool {
	_, ok := o.omap[userName]
	return ok
}

func (o *onlineMap) loginUser(userName string) {
	o.omap[userName] = true
}

func (o *onlineMap) offlineUser(userName string) {
	delete(o.omap, userName)
}
