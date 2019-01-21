package controller

import (
	log "github.com/Sirupsen/logrus"
	apps_v1 "k8s.io/api/apps/v1"
)

// Handler interface contains the methods that are required
type Handler interface {
	Init() error
	ObjectCreated(obj interface{})
	ObjectDeleted(obj interface{})
	ObjectUpdated(objOld, objNew interface{})
}

// TestHandler is a sample implementation of Handler
type TestHandler struct{}

// Init handles any handler initialization
func (t *TestHandler) Init() error {
	log.Info("TestHandler.Init")
	return nil
}

// ObjectCreated is called when an object is created
func (t *TestHandler) ObjectCreated(obj interface{}) {
	log.Info("TestHandler.ObjectCreated")
	// assert the type to a Deployment object to pull out relevant data
	deployment := obj.(*apps_v1.Deployment)
	log.Infof("    ResourceVersion: %s", deployment.ObjectMeta.ResourceVersion)
	log.Infof("    Replicas: %d", *deployment.Spec.Replicas)
}

// ObjectDeleted is called when an object is deleted
func (t *TestHandler) ObjectDeleted(obj interface{}) {
	log.Info("TestHandler.ObjectDeleted")
}

// ObjectUpdated is called when an object is updated
func (t *TestHandler) ObjectUpdated(objOld, objNew interface{}) {
	log.Info("TestHandler.ObjectUpdated")
}
