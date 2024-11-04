package eapp

import "context"

type appKey struct{}

func AppWithContext(ctx context.Context, app *App) context.Context {
	return context.WithValue(ctx, appKey{}, app)
}

func AppFromContext(ctx context.Context) (ReadApp, bool) {
	if app, ok := ctx.Value(appKey{}).(*App); ok {
		return ReadApp{
			app: app,
		}, true
	}
	return ReadApp{}, false
}
