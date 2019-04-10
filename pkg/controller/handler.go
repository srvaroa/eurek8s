package controller

import (
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/hudl/fargo"
	extensions_v1 "k8s.io/api/extensions/v1beta1"
)

// Handler interface contains the methods that are required
type Handler interface {
	Init() error
	ObjectCreated(keyRaw string, obj interface{})
	ObjectDeleted(keyRaw string, obj interface{})
	ObjectUpdated(keyRaw string, objOld, objNew interface{})
}

type EurekaSyncer struct {
	eureka   fargo.EurekaConnection
	liveChan chan *fargo.Instance
	deadChan chan string
}

// Init handles any handler initialization
func (e *EurekaSyncer) Init() error {
	log.Info("EurekaSyncer.Init")
	// TODO: pass the eureka endpoint via config
	e.eureka = fargo.NewConn("http://127.0.0.1:8080/eureka/v2")
	e.liveChan = make(chan *fargo.Instance)
	e.deadChan = make(chan string)
	go e.beat()
	return nil
}

func (e *EurekaSyncer) ObjectCreated(keyRaw string, obj interface{}) {
	log.Info("EurekaSyncer.ObjectCreated")
	e.reconcile(keyRaw, obj)
}

func (e *EurekaSyncer) ObjectDeleted(keyRaw string, obj interface{}) {
	log.Info("EurekaSyncer.ObjectDeleted")
	e.reconcile(keyRaw, obj)
}

func (e *EurekaSyncer) ObjectUpdated(keyRaw string, objOld, objNew interface{}) {
	log.Info("EurekaSyncer.ObjectUpdated")
	e.reconcile(keyRaw, objNew)
}

func (e *EurekaSyncer) reconcile(keyRaw string, obj interface{}) {
	if obj == nil {
		e.deadChan <- keyRaw
		return
	}
	ing := obj.(*extensions_v1.Ingress)
	appName := ing.Labels["app"]
	host := ing.Spec.Rules[0].Host
	if len(appName) == 0 {
		// TODO: error? e.deadChan <- keyRaw?
		return
	}
	log.Infof("Ingress found for app: %s, %s %+v", keyRaw, appName, host)

	e.liveChan <- &fargo.Instance{
		UniqueID: func(i fargo.Instance) string {
			return ing.Name
		},
		App:      appName,
		HostName: ing.Name,
		// TODO: do we want to set the LB's IP? (I think it can be
		// pulled from the ingress.Status)
		IPAddr:           host,
		VipAddress:       host,
		SecureVipAddress: host,
		Status:           fargo.UP,
		Port:             8080,
		DataCenterInfo:   fargo.DataCenterInfo{Name: fargo.MyOwn},
	}

	// log.Infof("-> Status: %+v", ing.Status.LoadBalancer)
	log.Infof("-> Host: %+v", ing.Spec.Rules[0].Host[0])
}

func (e *EurekaSyncer) beat() {
	instances := map[string]*fargo.Instance{}
	tickChan := time.NewTicker(time.Second * 10).C
	for {
		select {
		case _ = <-tickChan:
			log.Infof("Tick %d", len(instances))
			for _, i := range instances {
				log.Infof("Heartbeat for %s: %s", i.HostName, i.IPAddr)
				e.eureka.RegisterInstance(i)
			}
		case i := <-e.liveChan:
			key := i.UniqueID(*i)
			log.Infof("Live %s", key)
			instances[key] = i
			e.eureka.RegisterInstance(i)
		case key := <-e.deadChan:
			log.Infof("Dead %s", key)
			// The key is in the form $namespace/$pod_name, and it looks
			// like Eureka isn't able to handle / in the instance id.
			// It's either this, or s/\//_/ in the registration (which
			// maybe looks cleaner.  Reconsider sometime.)
			i := instances[strings.Split(key, "/")[1]]
			if i != nil {
				e.eureka.DeregisterInstance(i)
			}
			delete(instances, key)
		}
	}
}
