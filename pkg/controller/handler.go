package controller

import (
	log "github.com/Sirupsen/logrus"
	"github.com/hudl/fargo"
	core_v1 "k8s.io/api/core/v1"
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
	e.reconcile(obj.(*core_v1.Pod))
}

func (e *EurekaSyncer) ObjectDeleted(obj interface{}) {
	log.Info("EurekaSyncer.ObjectDeleted")
	// e.reconcile(obj.(*core_v1.Pod))
}

func (e *EurekaSyncer) ObjectUpdated(objOld, objNew interface{}) {
	log.Info("EurekaSyncer.ObjectUpdated")
	if objNew != nil {
		e.reconcile(objNew.(*core_v1.Pod))
	}
}

func (e *EurekaSyncer) reconcile(pod *core_v1.Pod) {
	if pod.Status.Phase == "Running" {
		instance := &fargo.Instance{
			UniqueID: func(i fargo.Instance) string {
				return pod.Name
			},
			App:      pod.Labels["app"],
			HostName: "host",
			// TODO: set the service DNS here (can we? Or a VIP?)
			IPAddr:           "192.168.1.1",
			VipAddress:       "192.168.1.1",
			SecureVipAddress: "192.168.1.1",
			Status:           fargo.UP,
			Port:             8080,
			DataCenterInfo:   fargo.DataCenterInfo{Name: fargo.MyOwn},
		}
		log.Infof("Registering instance %s", pod.Name)
		e.eureka.RegisterInstance(instance)
	} else {
		log.Infof("-> Phase: %+v", pod.Status.Phase)
	}
}
