package rpc

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/comm"
	"github.com/davyxu/cellnet/msglog"
)

type RemoteCallMsg interface {
	GetMsgID() uint32
	GetMsgData() []byte
	GetCallID() int64
}

type msgEvent interface {
	Session() cellnet.Session
	Message() interface{}
}

func ProcRPC(userFunc cellnet.EventProc) cellnet.EventProc {

	return func(raw cellnet.EventParam) cellnet.EventResult {

		if ev, ok := raw.(msgEvent); ok {
			rpcMsg, ok := ev.Message().(RemoteCallMsg)
			if ok {
				msg, meta, err := cellnet.DecodeMessage(rpcMsg.GetMsgID(), rpcMsg.GetMsgData())

				if err == nil {
					switch raw.(type) {
					case *cellnet.RecvMsgEvent:

						log.Debugf("#rpc recv(%s)@%d %s(%d) | %s",
							ev.Session().Peer().Name(),
							ev.Session().ID(),
							meta.Name,
							meta.ID,
							cellnet.MessageToString(msg))

						switch ev.Message().(type) {
						case *comm.RemoteCallREQ:

							userFunc(&RecvMsgEvent{ev.Session(), msg, rpcMsg.GetCallID()})

						case *comm.RemoteCallACK:
							request := getRequest(rpcMsg.GetCallID())
							if request != nil {
								request.RecvFeedback(msg)
							}
						}

					case *cellnet.SendMsgEvent:

						log.Debugf("#rpc send(%s)@%d %s(%d) | %s",
							ev.Session().Peer().Name(),
							ev.Session().ID(),
							meta.Name,
							meta.ID,
							cellnet.MessageToString(msg))

					}
				}
			}

		}

		return userFunc(raw)
	}
}

func init() {
	msglog.BlockMessageLog("comm.RemoteCallREQ")
	msglog.BlockMessageLog("comm.RemoteCallACK")
}
