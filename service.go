package common

import (
	"sync/atomic"

	"github.com/neatio-network/neatio/chain/log"
)

type Service interface {
	Start() (bool, error)
	OnStart() error

	Stop() bool
	OnStop()

	Reset() (bool, error)
	OnReset() error

	IsRunning() bool

	String() string
}

type BaseService struct {
	logger  log.Logger
	name    string
	started uint32
	stopped uint32
	Quit    chan struct{}

	impl Service
}

func NewBaseService(logger log.Logger, name string, impl Service) *BaseService {
	return &BaseService{
		logger: logger,
		name:   name,
		Quit:   make(chan struct{}),
		impl:   impl,
	}
}

func (bs *BaseService) Start() (bool, error) {
	if atomic.CompareAndSwapUint32(&bs.started, 0, 1) {
		if bs.logger != nil {
			bs.logger.Infof("Starting %v (%v)", bs.name, bs.impl)
		}

		if atomic.LoadUint32(&bs.stopped) == 1 {
			atomic.StoreUint32(&bs.stopped, 0)
		}

		bs.Quit = make(chan struct{})
		err := bs.impl.OnStart()
		if err != nil {

			atomic.StoreUint32(&bs.started, 0)
			if bs.logger != nil {
				bs.logger.Errorf("Starting failed with %v", err)
			}
			return false, err
		}

		if bs.logger != nil {
			bs.logger.Infof("Started %v (%v)", bs.name, bs.impl)
		}
		return true, err
	} else {
		if bs.logger != nil {
			bs.logger.Debugf("Not starting %v -- already started, impl: %v", bs.name, bs.impl)
		}
		return false, nil
	}
}

func (bs *BaseService) OnStart() error { return nil }

func (bs *BaseService) Stop() bool {
	if atomic.CompareAndSwapUint32(&bs.stopped, 0, 1) {
		if bs.logger != nil {
			bs.logger.Infof("Stopping %v (%v)", bs.name, bs.impl)
		}
		bs.impl.OnStop()
		close(bs.Quit)
		if atomic.LoadUint32(&bs.started) == 1 {
			atomic.StoreUint32(&bs.started, 0)
		}
		if bs.logger != nil {
			bs.logger.Infof("Stopped %v (%v)", bs.name, bs.impl)
		}
		return true
	} else {
		if bs.logger != nil {
			bs.logger.Debugf("Stopping %v (ignoring: already stopped) , impl: %v", bs.name, bs.impl)
		}
		return false
	}
}

func (bs *BaseService) OnStop() {}

func (bs *BaseService) Reset() (bool, error) {
	if atomic.CompareAndSwapUint32(&bs.stopped, 1, 0) {

		atomic.CompareAndSwapUint32(&bs.started, 1, 0)

		bs.Quit = make(chan struct{})
		return true, bs.impl.OnReset()
	} else {
		if bs.logger != nil {
			bs.logger.Debug("Can't reset ", bs.name, ". Not stopped, impl:", bs.impl)
		}
		return false, nil
	}

	return false, nil
}

func (bs *BaseService) OnReset() error {
	PanicSanity("The service cannot be reset")
	return nil
}

func (bs *BaseService) IsRunning() bool {
	return atomic.LoadUint32(&bs.started) == 1 && atomic.LoadUint32(&bs.stopped) == 0
}

func (bs *BaseService) Wait() {
	<-bs.Quit
}

func (bs *BaseService) String() string {
	return bs.name
}

type QuitService struct {
	BaseService
}

func NewQuitService(logger log.Logger, name string, impl Service) *QuitService {
	if logger != nil {
		logger.Warn("QuitService is deprecated, use BaseService instead")
	}
	return &QuitService{
		BaseService: *NewBaseService(logger, name, impl),
	}
}
