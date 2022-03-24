package models

// Generate these after models are done generating

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . AuthProvider
//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . StreamEventSource
//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . EntityEventSource
//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . StoreProvider
