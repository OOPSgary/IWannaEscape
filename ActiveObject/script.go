package ActiveObject

import "IJustWantToEscape/script"

/* TODO:
A multi Object Scipter
*/
type ScriptRunner struct {
	Channel chan script.Script
	Ticker  chan struct{}
	Object  *Object
}

type Text []struct {
}
