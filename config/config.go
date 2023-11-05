package config

/*
//The place we get args

var table map[string]*singleTag = make(map[string]*singleTag)
var globalMutex *sync.RWMutex = &sync.RWMutex{}
var WrongType = fmt.Errorf("The value cant cover to the Target Type")

type singleTag struct {
	m      *sync.RWMutex
	config map[string]interface{}
}

// func init() {
// 	Tagging("global")
// }

func Tagging(tag string) *singleTag {
	globalMutex.Lock()
	defer globalMutex.Unlock()
	st := &singleTag{
		m:      &sync.RWMutex{},
		config: make(map[string]interface{}, 100),
	}
	table[tag] = st
	return st
}
func Global() *singleTag {
	return GetRegion("global")
}
func GetRegion(tag string) *singleTag {
	globalMutex.RLock()

	if table[tag] != nil {
		defer globalMutex.RUnlock()
		return table[tag]
	}
	globalMutex.RUnlock()
	return Tagging(tag)

}
func (st *singleTag) Set(s string, v any) {
	st.m.Lock()
	defer st.m.Unlock()
	st.config[s] = v
}
func (st *singleTag) GetBool(s string) (b bool, err error) {
	st.m.RLock()
	defer st.m.RUnlock()
	var ok bool
	if b, ok = st.config[s].(bool); ok {
		return
	}
	err = WrongType
	return
}
func (st *singleTag) GetString(s string) (str string, err error) {
	st.m.RLock()
	defer st.m.RUnlock()
	var ok bool
	if str, ok = st.config[s].(string); ok {
		return
	}
	err = WrongType
	return
}
func (st *singleTag) GetInt(s string) (i int, err error) {
	st.m.RLock()
	defer st.m.RUnlock()
	var ok bool
	if i, ok = st.config[s].(int); ok {
		return
	}
	err = WrongType
	return
}
func (st *singleTag) Get(s string) (v interface{}, err error) {
	st.m.RLock()
	defer st.m.RUnlock()
	var ok bool
	if v, ok = st.config[s].(interface{}); ok {
		return
	}
	err = WrongType
	return
}

*/
