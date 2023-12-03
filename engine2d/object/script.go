package Object

import "IJustWantToEscape/script"

/* TODO:
A multi Object Scipter
*/
type ScriptRunner struct {
	Channel chan script.Script
	Object  *StatefullObject
}

type Text []struct {
}
