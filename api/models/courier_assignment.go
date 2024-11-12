package models

type CourierAssignment struct {
	Id         string `json:"id"`
	OrderId    string `json:"order_id"`
	CourierId  string `json:"courier_id"`
	Status     string `json:"status"`
	AssignedAt string `json:"assigned_at"`
	UpdatedAt  string `json:"updated_at"`
}

type CreateCourierAssignment struct {
	OrderId   string `json:"order_id"`
	CourierId string `json:"courier_id"`
	Status    string `json:"status"`
}

type UpdateCourierAssignment struct {
	OrderId   string `json:"order_id"`
	CourierId string `json:"courier_id"`
	Status    string `json:"status"`
}

type GetCourierAssignment struct {
	Id         string `json:"id"`
	OrderId    string `json:"order_id"`
	CourierId  string `json:"courier_id"`
	Status     string `json:"status"`
	AssignedAt string `json:"assigned_at"`
	UpdatedAt  string `json:"updated_at"`
}

type GetAllCourierAssignmentsRequest struct {
	Search string `json:"search"`
	Page   uint64 `json:"page"`
	Limit  uint64 `json:"limit"`
}

type GetAllCourierAssignmentsResponse struct {
	CourierAssignments []CourierAssignment `json:"courier_assignments"`
	Count              int64               `json:"count"`
}
