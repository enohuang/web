package session

import (
	"dengming20240317/web"
	"github.com/google/uuid"
)

type Manager struct {
	Propagator
	Store
	CtxSessKey string
}

func (m *Manager) GetSession(ctx *web.Context) (Session, error) {
	if ctx.UserValues == nil {
		ctx.UserValues = make(map[string]any)
	}

	val, ok := ctx.UserValues[m.CtxSessKey /*"_sess"*/]
	if ok {
		return val.(Session), nil
	}

	//方法2 context 缓存
	ctx.Req.Context().Value(m.CtxSessKey)

	sessId, err := m.Extract(ctx.Req)
	if err != nil {
		return nil, err
	}
	sess, err := m.Get(ctx.Req.Context(), sessId)
	if err != nil {
		return nil, err
	}

	ctx.UserValues[m.CtxSessKey /*"_sess"*/] = sess

	//可以用context  复制性能很差
	// ctx.Req = ctx.Req.WithContext(context.WithValue(ctx.Req.Context(), m.CtxSessKey, sess))

	return sess, err

}

func (m *Manager) InitSession(ctx *web.Context) (Session, error) {
	id := uuid.New().String()
	sess, err := m.Generate(ctx.Req.Context(), id)
	if err != nil {
		return nil, err
	}
	//注入进去  HTTP响应里面
	m.Inject(id, ctx.Resp)
	return sess, nil
}

func (m *Manager) RemoveSession(ctx *web.Context) error {

	sess, err := m.GetSession(ctx)
	if err != nil {
		return err
	}
	m.Store.Remove(ctx.Req.Context(), sess.ID())
	m.Propagator.Remove(ctx.Resp)
	return nil
}

func (m *Manager) RefreshSession(ctx *web.Context) error {
	sess, err := m.GetSession(ctx)
	if err != nil {
		return err
	}

	return m.Refresh(ctx.Req.Context(), sess.ID())
}
