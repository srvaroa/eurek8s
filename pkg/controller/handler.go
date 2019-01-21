package controller

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/hudl/fargo"
	apps_v1 "k8s.io/api/apps/v1"
)

// Handler interface contains the methods that are required
type Handler interface {
	Init() error
	ObjectCreated(obj interface{})
	ObjectDeleted(obj interface{})
	ObjectUpdated(objOld, objNew interface{})
}

type EurekaSyncer struct {
	eureka fargo.EurekaConnection
}

// Init handles any handler initialization
func (e *EurekaSyncer) Init() error {
	log.Info("EurekaSyncer.Init")
	e.eureka = fargo.NewConn("http://127.0.0.1:8080/eureka/v2")
	return nil
}

func (e *EurekaSyncer) ObjectCreated(obj interface{}) {
	log.Info("EurekaSyncer.ObjectCreated")
	// assert the type to a Deploymen  object to pull out relevant data
	deployment := obj.(*apps_v1.Deployment)
	log.Infof("    ResourceVersion: %s", deployment.ObjectMeta.ResourceVersion)
	log.Infof("    Replicas: %d", deployment.Status.ReadyReplicas)

	instance := &fargo.Instance{
		UniqueID: func(i fargo.Instance) string {
			return fmt.Sprintf("%s:%s", "192.168.1.1", "asdfas")
		},
		App:              deployment.Name,
		HostName:         "host",
		IPAddr:           "192.168.1.1",
		VipAddress:       "192.168.1.1",
		SecureVipAddress: "192.168.1.1",
		Status:           fargo.UP,
		Port:             8080,
		DataCenterInfo:   fargo.DataCenterInfo{Name: fargo.MyOwn},
	}
	// log.Infof("Registering instance", instance)
	e.eureka.RegisterInstance(instance)

}

func (e *EurekaSyncer) ObjectDeleted(obj interface{}) {
	log.Info("EurekaSyncer.ObjectDeleted")
}

func (e *EurekaSyncer) ObjectUpdated(objOld, objNew interface{}) {
	log.Info("EurekaSyncer.ObjectUpdated")
}
