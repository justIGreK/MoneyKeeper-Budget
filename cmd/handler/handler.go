package handler

import (
  "google.golang.org/grpc"
  budgetProto "github.com/justIGreK/MoneyKeeper-Budget/pkg/go/budget"
)

type Handler struct {
	server grpc.ServiceRegistrar
	budget BudgetService
}

func NewHandler(grpcServer grpc.ServiceRegistrar, budgetSRV BudgetService) *Handler {
	return &Handler{server: grpcServer, budget: budgetSRV}
}
func (h *Handler) RegisterServices() {
	h.registerBudgetService(h.server, h.budget)
}

func (h *Handler) registerBudgetService(server grpc.ServiceRegistrar, budget BudgetService) {
	budgetProto.RegisterBudgetServiceServer(server, &BudgetServiceServer{BudgetSRV: budget})
}
