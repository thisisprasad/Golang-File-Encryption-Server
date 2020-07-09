package fesserver

import "fes/S_DES"

type FesServer struct {
	AppName   string `json:"name"`
	Port      int    `json:"port"`
	Protocol  string `protocol:"protocol"`
	Fesengine S_DES.DesEncryptor
}
