package config

import "sync"

func NewTransformConfig() *TransformConfig {
	return &TransformConfig{
		zmodemConfig: &ZModemConfig{},
	}
}

type TransformConfig struct {
	BufferSize   int
	zmodemConfig *ZModemConfig
}

func (t *TransformConfig) ZModemConfig() *ZModemConfig {
	return t.zmodemConfig
}

func (t *TransformConfig) SetDefaults() {
	if t.BufferSize < 1 {
		t.BufferSize = 8192
	}
	if t.zmodemConfig == nil {
		t.zmodemConfig = &ZModemConfig{}
	}
}

type ZModemConfig struct {
	DisableZModemSZ, DisableZModemRZ bool
	ZModemSZ, ZModemRZ, ZModemSZOO   bool

	mutexDisableZModemSZ, mutexDisableZModemRZ    sync.RWMutex
	mutexZModemSZ, mutexZModemRZ, mutexZModemSZOO sync.RWMutex
}

func (z *ZModemConfig) GetDisableZModemSZ() bool {
	z.mutexDisableZModemSZ.RLock()
	v := z.DisableZModemSZ
	z.mutexDisableZModemSZ.RUnlock()
	return v
}

func (z *ZModemConfig) GetDisableZModemRZ() bool {
	z.mutexDisableZModemRZ.RLock()
	v := z.DisableZModemRZ
	z.mutexDisableZModemRZ.RUnlock()
	return v
}

func (z *ZModemConfig) GetZModemSZ() bool {
	z.mutexZModemSZ.RLock()
	v := z.ZModemSZ
	z.mutexZModemSZ.RUnlock()
	return v
}

func (z *ZModemConfig) GetZModemRZ() bool {
	z.mutexZModemRZ.RLock()
	v := z.ZModemRZ
	z.mutexZModemRZ.RUnlock()
	return v
}

func (z *ZModemConfig) GetZModemSZOO() bool {
	z.mutexZModemSZOO.RLock()
	v := z.ZModemSZOO
	z.mutexZModemSZOO.RUnlock()
	return v
}

func (z *ZModemConfig) SetDisableZModemSZ(on bool) {
	z.mutexDisableZModemSZ.Lock()
	z.DisableZModemSZ = on
	z.mutexDisableZModemSZ.Unlock()
}

func (z *ZModemConfig) SetDisableZModemRZ(on bool) {
	z.mutexDisableZModemRZ.Lock()
	z.DisableZModemRZ = on
	z.mutexDisableZModemRZ.Unlock()
}

func (z *ZModemConfig) SetZModemSZ(on bool) {
	z.mutexZModemSZ.Lock()
	z.ZModemSZ = on
	z.mutexZModemSZ.Unlock()
}

func (z *ZModemConfig) SetZModemRZ(on bool) {
	z.mutexZModemRZ.Lock()
	z.ZModemRZ = on
	z.mutexZModemRZ.Unlock()
}

func (z *ZModemConfig) SetZModemSZOO(on bool) {
	z.mutexZModemSZOO.Lock()
	z.ZModemSZOO = on
	z.mutexZModemSZOO.Unlock()
}
