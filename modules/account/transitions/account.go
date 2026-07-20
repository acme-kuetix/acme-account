package transitions

import (
	"github.com/kuetix/engine/engine/domain/interfaces"
	"github.com/kuetix/engine/engine/workflow"
)

type accountTransitions struct {
	workflow.BaseServiceTransition
}

func NewAccountTransitions() interfaces.ServiceTransitions {
	return &accountTransitions{}
}

const accountsCollection = "accounts"
