package model

type Ingress struct {
	NamespaceWithName
	RouteHost string      `json:"routeHost"`
	RoutePath []RoutePath `gorm:"foreignKey:RouteID" json:"routePath"`
}