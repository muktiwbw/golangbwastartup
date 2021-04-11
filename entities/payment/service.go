package payment

import (
	"os"

	"github.com/veritrans/go-midtrans"
)

type Service interface {
	GenerateSnapLink(map[string]interface{}) (midtrans.SnapResponse, error)
}

type service struct{}

func NewService() Service {
	return &service{}
}

func (s *service) GenerateSnapLink(snapReqData map[string]interface{}) (midtrans.SnapResponse, error) {
	midclient := midtrans.NewClient()
	midclient.ServerKey = os.Getenv("MIDTRANS_SERVER_KEY")
	midclient.ClientKey = os.Getenv("MIDTRANS_CLIENT_KEY")
	midclient.APIEnvType = midtrans.Sandbox

	snapGateway := midtrans.SnapGateway{
		Client: midclient,
	}

	snapReq := &midtrans.SnapReq{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  snapReqData["transaction"].(map[string]interface{})["orderID"].(string),
			GrossAmt: snapReqData["transaction"].(map[string]interface{})["grossAmt"].(int64),
		},
		CustomerDetail: &midtrans.CustDetail{
			FName: snapReqData["customer"].(map[string]interface{})["fName"].(string),
			LName: snapReqData["customer"].(map[string]interface{})["lName"].(string),
			Email: snapReqData["customer"].(map[string]interface{})["email"].(string),
		},
		Items: &[]midtrans.ItemDetail{
			midtrans.ItemDetail{
				ID:    snapReqData["project"].(map[string]interface{})["id"].(string),
				Price: snapReqData["transaction"].(map[string]interface{})["grossAmt"].(int64),
				Qty:   1,
				Name:  snapReqData["project"].(map[string]interface{})["name"].(string),
			},
		},
	}

	snapTokenResp, err := snapGateway.GetToken(snapReq)

	if err != nil {
		return snapTokenResp, err
	}

	return snapTokenResp, nil
}
