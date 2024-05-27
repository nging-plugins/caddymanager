package goth

import "fmt"

// ==========================
// ==== default instance ====
// ==========================

func RangeProviders(cb func(string, Provider) bool) bool {
	return providers.Range(cb)
}

func DeleteProvider(names ...string) {
	providers.Delete(names...)
}

// ==========================
// == Providers definition ==
// ==========================

func NewProviders() Providers {
	return Providers{
		m: map[string]Provider{},
	}
}

func (p *Providers) Add(viders ...Provider) {
	p.l.Lock()
	for _, provider := range viders {
		p.m[provider.Name()] = provider
	}
	p.l.Unlock()
}

func (p *Providers) Delete(names ...string) {
	p.l.Lock()
	for _, name := range names {
		delete(p.m, name)
	}
	p.l.Unlock()
}

func (p *Providers) Get(name string) (Provider, error) {
	p.l.RLock()
	provider := p.m[name]
	p.l.RUnlock()
	if provider == nil {
		return nil, fmt.Errorf("no provider for %s exists", name)
	}
	return provider, nil
}

func (p *Providers) All() map[string]Provider {
	cloned := map[string]Provider{}
	p.l.RLock()
	for name, provider := range p.m {
		cloned[name] = provider
	}
	p.l.RUnlock()
	return cloned
}

func (p *Providers) Range(cb func(string, Provider) bool) (ok bool) {
	p.l.RLock()
	for name, provider := range p.m {
		ok = cb(name, provider)
		if !ok {
			break
		}
	}
	p.l.RUnlock()
	return
}

func (p *Providers) Clear() {
	p.l.Lock()
	p.m = map[string]Provider{}
	p.l.Unlock()
}
