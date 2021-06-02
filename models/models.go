package models

import persistenceLayer "hellofresh/elastic"

type PersistanceType struct{}

func (e *PersistanceType) InitClient() {
	persistenceLayer.InitClient()
}

func (e *PersistanceType) IndexExists() (bool, error) {
	return persistenceLayer.IndexExists()
}

func (e *PersistanceType) CreateIndex() error {
	return persistenceLayer.CreateIndex()
}

func (e *PersistanceType) CreateMapping() error {
	return persistenceLayer.CreateMapping()
}
