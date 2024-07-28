package config

import (
	"net/http"
	"game/pb"
)

type BaseConfReq struct {
}

type BaseConfRes struct {
	pb.ProBaseResponse
}

func BaseConfigHandler(w http.ResponseWriter, r *http.Request)  {

}
