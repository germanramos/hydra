package balancer

import (
	. "github.com/innotech/hydra/model/entity"
	. "github.com/innotech/hydra/model/repository"

	zmq "github.com/innotech/hydra/vendors/github.com/alecthomas/gozmq"
)

type Balancer struct {
	Plumbers      map[string]Plumber `map of: key => application_id, value => Plumber`
	AppRepository EtcdBaseRepository
}

func NewBalancer() *Balancer {
	b := new(Balancer)
	// TODO: Service Container with repositories?
	b.AppRepository = NewEctdRepository("applications")
	return b
}

func (b *Balancer) Configure() error {
	apps, err := b.AppRepository.GetAll()
	if err != nil {
		return err
	}

}

func (b *Balancer) RegisterPlumbers(apps []EtcdBaseModel) {
	for _, app := range apps {
		RegisterPlumber(app)
	}
}

func (b *Balancer) RegisterPlumber(app EtcdBaseModel) {
	var appId string
	var appProps map[string]interface{}
	// TODO:
	for appId, appProps = range app {
	}
	b.Plumbers[appId] = NewPlumber()
}

func (b *Balancer) Start(map[string][]map[string]interface{}) {
	context, _ := zmq.NewContext()
	defer context.Close()

	//  Frontend socket talks to clients over TCP
	// frontend, _ := context.NewSocket(zmq.ROUTER)
	// frontend.Bind("ipc://frontend.ipc")
	// defer frontend.Close()

	//  Backend socket talks to workers over inproc
	backend, _ := context.NewSocket(zmq.DEALER)
	// backend.Bind("ipc://backend.ipc")
	// TODO: Config address
	backend.Bind("tcp://127.0.0.1:4444")
	defer backend.Close()

	//  Launch pool of worker threads, precise number is not critical
	// for i := 0; i < 5; i++ {
	// 	go server_worker()
	// }

	//  Connect backend to frontend via a proxy
	// items := zmq.PollItems{
	// 	zmq.PollItem{Socket: frontend, Events: zmq.POLLIN},
	// 	zmq.PollItem{Socket: backend, Events: zmq.POLLIN},
	// }

	for {
		_, err := zmq.Poll(items, -1)
		if err != nil {
			fmt.Println("Server exited with error:", err)
			break
		}

		if items[0].REvents&zmq.POLLIN != 0 {

			parts, _ := frontend.RecvMultipart(0)
			backend.SendMultipart(parts, 0)

		}
		if items[1].REvents&zmq.POLLIN != 0 {

			parts, _ := backend.RecvMultipart(0)
			frontend.SendMultipart(parts, 0)
		}
	}
}

func (b *Balancer) Balance(app map[string]interface{}) {

}
