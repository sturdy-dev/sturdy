package graphql

/*
It's currently not possible to disable the server-generation in gqlgen.
And the generated.go file must always be created (the location is configurable however)
We don't need the file, so the second generation step will delete the file if it exists.
*/

//go:generate go run github.com/99designs/gqlgen
//go:generate /bin/bash -c "rm generated.go || true"
