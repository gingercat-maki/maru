package modules

import (
	"github.com/temporalio/samples-go/saga"
	"go.temporal.io/sdk/worker"
)

func RegisterSagaTransfer(w worker.Worker) {
	w.RegisterActivity(saga.Withdraw)
	w.RegisterActivity(saga.WithdrawCompensation)
	w.RegisterActivity(saga.Deposit)
	w.RegisterActivity(saga.DepositCompensation)
	w.RegisterActivity(saga.StepWithError)
}

type SagaTransferRequest struct {
	EnabledSagaWithTransfer bool                 `json:"enabledSagaWithTransfer"`
	TransferDetailsRequest  saga.TransferDetails `json:"transferDetails"`
}
