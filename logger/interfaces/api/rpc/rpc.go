package rpc

// import (
// 	"log"

// 	ml "github.com/yaroyan/ms/logger/domain/model/log"
// 	"github.com/yaroyan/ms/logger/infrastructure/persistence"
// )

// type RPCServer struct {
// 	Repository persistence.LogRepository
// }

// type RPCPayload struct {
// 	Name string
// 	Data string
// }

// func (r *RPCServer) LogInfo(payload RPCPayload, resp *string) error {
// 	e := ml.Log{
// 		Name: payload.Name,
// 		Data: payload.Data,
// 	}
// 	err := r.Repository.Insert(e)

// 	if err != nil {
// 		log.Println("error writing to mongo", err)
// 		return err
// 	}

// 	*resp = "Processed payload via RPC: " + payload.Name
// 	return nil
// }
