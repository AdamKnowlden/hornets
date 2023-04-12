package sessions

type Session struct {
	ClientID uint64
	Data     map[string]interface{}
}

var Sessions = make(map[uint64]Session)

func CreateSession(clientID uint64) Session {
	session := &Session{
		ClientID: clientID,
		Data:     make(map[string]interface{}),
	}

	Sessions[clientID] = *session

	return *session
}

func GetSession(clientID uint64) *Session {
	if CheckSession(clientID) {
		session := Sessions[clientID]

		return &session
	} else {
		return nil
	}
}

func (session *Session) UpdateData(key string, value interface{}) {
	session.Data[key] = value
}

func (session *Session) GetData(key string) interface{} {
	value, ok := session.Data[key]
	if ok {
		return value
	} else {
		return nil
	}
}

func EndSession(clientID uint64) {
	delete(Sessions, clientID)
}

func CheckSession(clientID uint64) bool {
	value, ok := Sessions[clientID]
	if ok && value.ClientID == clientID {
		return true
	} else {
		return false
	}
}
