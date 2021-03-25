package webreal

func (c *Client) CtxLoad(key interface{}) (interface{}, bool) {
	return c.ctx.Load(key)
}

func (c *Client) CtxStore(key interface{}, value interface{}) {
	c.ctx.Store(key, value)
}

func (c *Client) CtxDelete(key interface{}) {
	c.ctx.Delete(key)
}

func (c *Client) CtxRange(f func(key, value interface{}) bool) {
	c.ctx.Range(f)
}
