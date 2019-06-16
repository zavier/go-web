package memory

import (
	"../session"
	"container/list"
	"sync"
	"time"
)

type Provider struct {
	lock     sync.Mutex
	sessions map[string]*list.Element
	// 回收时，只判断最后一个原始是否过期，每次访问节点会把对应节点放到头部
	list *list.List
}

func init() {
	pder.sessions = make(map[string]*list.Element, 0)
	session.Register("memory", pder)
}

func (pder *Provider) SessionInit(sid string) (session.Session, error) {
	pder.lock.Lock()
	defer pder.lock.Unlock()
	v := make(map[interface{}]interface{}, 0)
	newSession := &SessionStore{sid: sid, timeAccessed: time.Now(), value: v}
	element := pder.list.PushBack(newSession)
	pder.sessions[sid] = element
	return newSession, nil
}

func (pder *Provider) SessionRead(sid string) (session.Session, error) {
	if element, ok := pder.sessions[sid]; ok {
		return element.Value.(*SessionStore), nil
	} else {
		session, err := pder.SessionInit(sid)
		return session, err
	}
	return nil, nil
}

func (pdef *Provider) SessionDestroy(sid string) error {
	if elemenet, ok := pder.sessions[sid]; ok {
		delete(pder.sessions, sid)
		pder.list.Remove(elemenet)
		return nil
	}
	return nil
}

func (pdef *Provider) SessionGC(maxlifetime int64) {
	pdef.lock.Lock()
	defer pder.lock.Unlock()

	for {
		// 判断最后一个，过期则删除
		element := pder.list.Back()
		if element == nil {
			break
		}
		if (element.Value.(*SessionStore).timeAccessed.Unix() + maxlifetime) < time.Now().Unix() {
			pder.list.Remove(element)
			delete(pder.sessions, element.Value.(*SessionStore).sid)
		} else {
			break
		}
	}
}

// 将节点前置
func (pder *Provider) SessionUpdate(sid string) error {
	pder.lock.Lock()
	defer pder.lock.Unlock()
	if element, ok := pder.sessions[sid]; ok {
		element.Value.(*SessionStore).timeAccessed = time.Now()
		pder.list.MoveToFront(element)
		return nil
	}
	return nil
}
