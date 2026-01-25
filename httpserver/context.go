package httpserver

import "context"

type contextKey string

const (
	RouteNameKey   contextKey = "route_name"
	RouteMethodKey contextKey = "route_method"
	RoutePathKey   contextKey = "route_path"
	AuthTypeKey    contextKey = "auth_type"
	RequestAuthKey contextKey = "request_auth"
)

// GetRouteName extracts the route name from context
func GetRouteName(ctx context.Context) string {
	if val := ctx.Value(RouteNameKey); val != nil {
		return val.(string)
	}
	return ""
}

// GetRouteMethod extracts the route method from context
func GetRouteMethod(ctx context.Context) string {
	if val := ctx.Value(RouteMethodKey); val != nil {
		return val.(string)
	}
	return ""
}

// GetRoutePath extracts the route path template from context
func GetRoutePath(ctx context.Context) string {
	if val := ctx.Value(RoutePathKey); val != nil {
		return val.(string)
	}
	return ""
}

// GetAuthType extracts the auth type from context
func GetAuthType(ctx context.Context) string {
	if val := ctx.Value(AuthTypeKey); val != nil {
		return val.(string)
	}
	return ""
}

// GetRequestAuth extracts the request authentication details from context
func GetRequestAuth(ctx context.Context) *RequestAuth {
	if val := ctx.Value(RequestAuthKey); val != nil {
		if auth, ok := val.(RequestAuth); ok {
			return &auth
		}
	}
	return nil
}