package db

import "sync"

// ServiceManager holds services
type ServiceManager struct {
	services map[string]*SparrowService
	Active   bool
}

// SparrowService interface for services
type SparrowService interface {
	Start()
	Stop()
}

// AddService add service
func (bge *ServiceManager) AddService(name string, v SparrowService) {
	bge.services[name] = &v
}

// StartAll starts all services
func (bge *ServiceManager) StartAll() {
	bge.Active = true

	var wg sync.WaitGroup

	wg.Add(len(bge.services))

	for _, v := range bge.services {
		go func(service *SparrowService) {
			defer wg.Done()
			(*service).Start()
		}(v)
	}

	wg.Wait()
}

// StopAll stops all services
func (bge *ServiceManager) StopAll() {
	bge.Active = false

	for _, v := range bge.services {
		go func(service *SparrowService) {
			(*service).Stop()
		}(v)
	}
}

// NewServiceManager returns new ServiceManager
func NewServiceManager() ServiceManager {
	return ServiceManager{
		services: make(map[string]*SparrowService),
		Active:   false,
	}
}
