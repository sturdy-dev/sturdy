package routes

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/getsentry/sentry-go"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestStore_sentryPayload(t *testing.T) {
	nopClient, err := sentry.NewClient(sentry.ClientOptions{})
	assert.NoError(t, err)
	logger, _ := zap.NewDevelopment()
	route := Store(logger, nopClient)

	req, err := http.NewRequest("POST", "/", strings.NewReader(sentryPayload))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	route.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestStore_ravenPayload(t *testing.T) {
	nopClient, err := sentry.NewClient(sentry.ClientOptions{})
	assert.NoError(t, err)
	logger, _ := zap.NewDevelopment()
	route := Store(logger, nopClient)

	req, err := http.NewRequest("POST", "/", strings.NewReader(ravenPayload))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	route.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

// this is an example of a payload produced by github.com/getsentry/raven-go library
const ravenPayload = `
	{
  "culprit": "getsturdy.com/api/pkg/logger.New",
  "environment": "enterprise",
  "event_id": "c03b1664649b401fa1b32d1ce9a94bb6",
  "exception": {
    "stacktrace": {
      "frames": [
        {
          "abs_path": "/opt/homebrew/Cellar/go/1.18/libexec/src/runtime/proc.go",
          "context_line": "\tfn()",
          "filename": "runtime/proc.go",
          "function": "main",
          "in_app": false,
          "lineno": 250,
          "module": "runtime",
          "post_context": [
            "\tif raceenabled {",
            "\t\tracefini()",
            "\t}"
          ],
          "pre_context": [
            "\t\treturn",
            "\t}",
            "\tfn := main_main // make an indirect call, as the linker doesn't know the addressof the main package when laying down the runtime"
          ]
        },
        {
          "abs_path": "/Users/nikita.galaiko/src/sturdy/api/cmd/api/main.go",
          "context_line": "\tif err := di.Init(app).To(\u0026apiServer, \u0026ctx); err != nil {",
          "filename": "/Users/nikita.galaiko/src/sturdy/api/cmd/api/main.go",
          "function": "main",
          "in_app": true,
          "lineno": 24,
          "module": "main",
          "post_context": [
            "\t\tfmt.Printf(\"%+v\\n\", err)",
            "\t\tos.Exit(1)",
            "\t}"
          ],
          "pre_context": [
            "",
            "\tvar apiServer api.Starter",
            "\tvar ctx context.Context"
          ]
        },
        {
          "abs_path": "/Users/nikita.galaiko/src/sturdy/api/pkg/di/container.go",
          "context_line": "\tif err := container.Invoke(invokeFn.Interface()); err != nil {",
          "filename": "/Users/nikita.galaiko/src/sturdy/api/pkg/di/container.go",
          "function": "To",
          "in_app": true,
          "lineno": 360,
          "module": "getsturdy.com/api/pkg/di.(*Container)",
          "post_context": [
            "\t\treturn err",
            "\t}",
            ""
          ],
          "pre_context": [
            "\t\treturn []reflect.Value{}",
            "\t})",
            ""
          ]
        },
        {
          "abs_path": "/Users/nikita.galaiko/go/pkg/mod/go.uber.org/dig@v1.14.1/invoke.go",
          "context_line": "\targs, err := pl.BuildList(s)",
          "filename": "/Users/nikita.galaiko/go/pkg/mod/go.uber.org/dig@v1.14.1/invoke.go",
          "function": "Invoke",
          "in_app": false,
          "lineno": 85,
          "module": "go.uber.org/dig.(*Scope)",
          "post_context": [
            "\tif err != nil {",
            "\t\treturn errArgumentsFailed{",
            "\t\t\tFunc:   digreflect.InspectFunc(function),"
          ],
          "pre_context": [
            "\t\ts.isVerifiedAcyclic = true",
            "\t}",
            ""
          ]
        },
        {
          "abs_path": "/Users/nikita.galaiko/go/pkg/mod/go.uber.org/dig@v1.14.1/param.go",
          "context_line": "\t\targs[i], err = p.Build(c)",
          "filename": "/Users/nikita.galaiko/go/pkg/mod/go.uber.org/dig@v1.14.1/param.go",
          "function": "BuildList",
          "in_app": false,
          "lineno": 151,
          "module": "go.uber.org/dig.paramList",
          "post_context": [
            "\t\tif err != nil {",
            "\t\t\treturn nil, err",
            "\t\t}"
          ],
          "pre_context": [
            "\targs := make([]reflect.Value, len(pl.Params))",
            "\tfor i, p := range pl.Params {",
            "\t\tvar err error"
          ]
        },
        {
          "abs_path": "/Users/nikita.galaiko/go/pkg/mod/go.uber.org/dig@v1.14.1/param.go",
          "context_line": "\t\terr := n.Call(n.OrigScope())",
          "filename": "/Users/nikita.galaiko/go/pkg/mod/go.uber.org/dig@v1.14.1/param.go",
          "function": "Build",
          "in_app": false,
          "lineno": 296,
          "module": "go.uber.org/dig.paramSingle",
          "post_context": [
            "\t\tif err == nil {",
            "\t\t\tcontinue",
            "\t\t}"
          ],
          "pre_context": [
            "\t}",
            "",
            "\tfor _, n := range providers {"
          ]
        },
        {
          "abs_path": "/Users/nikita.galaiko/go/pkg/mod/go.uber.org/dig@v1.14.1/constructor.go",
          "context_line": "\targs, err := n.paramList.BuildList(c)",
          "filename": "/Users/nikita.galaiko/go/pkg/mod/go.uber.org/dig@v1.14.1/constructor.go",
          "function": "Call",
          "in_app": false,
          "lineno": 145,
          "module": "go.uber.org/dig.(*constructorNode)",
          "post_context": [
            "\tif err != nil {",
            "\t\treturn errArgumentsFailed{",
            "\t\t\tFunc:   n.location,"
          ],
          "pre_context": [
            "\t\t}",
            "\t}",
            ""
          ]
        },
        {
          "abs_path": "/Users/nikita.galaiko/go/pkg/mod/go.uber.org/dig@v1.14.1/param.go",
          "context_line": "\t\targs[i], err = p.Build(c)",
          "filename": "/Users/nikita.galaiko/go/pkg/mod/go.uber.org/dig@v1.14.1/param.go",
          "function": "BuildList",
          "in_app": false,
          "lineno": 151,
          "module": "go.uber.org/dig.paramList",
          "post_context": [
            "\t\tif err != nil {",
            "\t\t\treturn nil, err",
            "\t\t}"
          ],
          "pre_context": [
            "\targs := make([]reflect.Value, len(pl.Params))",
            "\tfor i, p := range pl.Params {",
            "\t\tvar err error"
          ]
        },
        {
          "abs_path": "/Users/nikita.galaiko/go/pkg/mod/go.uber.org/dig@v1.14.1/param.go",
          "context_line": "\t\terr := n.Call(n.OrigScope())",
          "filename": "/Users/nikita.galaiko/go/pkg/mod/go.uber.org/dig@v1.14.1/param.go",
          "function": "Build",
          "in_app": false,
          "lineno": 296,
          "module": "go.uber.org/dig.paramSingle",
          "post_context": [
            "\t\tif err == nil {",
            "\t\t\tcontinue",
            "\t\t}"
          ],
          "pre_context": [
            "\t}",
            "",
            "\tfor _, n := range providers {"
          ]
        },
        {
          "abs_path": "/Users/nikita.galaiko/go/pkg/mod/go.uber.org/dig@v1.14.1/constructor.go",
          "context_line": "\targs, err := n.paramList.BuildList(c)",
          "filename": "/Users/nikita.galaiko/go/pkg/mod/go.uber.org/dig@v1.14.1/constructor.go",
          "function": "Call",
          "in_app": false,
          "lineno": 145,
          "module": "go.uber.org/dig.(*constructorNode)",
          "post_context": [
            "\tif err != nil {",
            "\t\treturn errArgumentsFailed{",
            "\t\t\tFunc:   n.location,"
          ],
          "pre_context": [
            "\t\t}",
            "\t}",
            ""
          ]
        },
        {
          "abs_path": "/Users/nikita.galaiko/go/pkg/mod/go.uber.org/dig@v1.14.1/param.go",
          "context_line": "\t\targs[i], err = p.Build(c)",
          "filename": "/Users/nikita.galaiko/go/pkg/mod/go.uber.org/dig@v1.14.1/param.go",
          "function": "BuildList",
          "in_app": false,
          "lineno": 151,
          "module": "go.uber.org/dig.paramList",
          "post_context": [
            "\t\tif err != nil {",
            "\t\t\treturn nil, err",
            "\t\t}"
          ],
          "pre_context": [
            "\targs := make([]reflect.Value, len(pl.Params))",
            "\tfor i, p := range pl.Params {",
            "\t\tvar err error"
          ]
        },
        {
          "abs_path": "/Users/nikita.galaiko/go/pkg/mod/go.uber.org/dig@v1.14.1/param.go",
          "context_line": "\t\terr := n.Call(n.OrigScope())",
          "filename": "/Users/nikita.galaiko/go/pkg/mod/go.uber.org/dig@v1.14.1/param.go",
          "function": "Build",
          "in_app": false,
          "lineno": 296,
          "module": "go.uber.org/dig.paramSingle",
          "post_context": [
            "\t\tif err == nil {",
            "\t\t\tcontinue",
            "\t\t}"
          ],
          "pre_context": [
            "\t}",
            "",
            "\tfor _, n := range providers {"
          ]
        },
        {
          "abs_path": "/Users/nikita.galaiko/go/pkg/mod/go.uber.org/dig@v1.14.1/constructor.go",
          "context_line": "\targs, err := n.paramList.BuildList(c)",
          "filename": "/Users/nikita.galaiko/go/pkg/mod/go.uber.org/dig@v1.14.1/constructor.go",
          "function": "Call",
          "in_app": false,
          "lineno": 145,
          "module": "go.uber.org/dig.(*constructorNode)",
          "post_context": [
            "\tif err != nil {",
            "\t\treturn errArgumentsFailed{",
            "\t\t\tFunc:   n.location,"
          ],
          "pre_context": [
            "\t\t}",
            "\t}",
            ""
          ]
        },
        {
          "abs_path": "/Users/nikita.galaiko/go/pkg/mod/go.uber.org/dig@v1.14.1/param.go",
          "context_line": "\t\targs[i], err = p.Build(c)",
          "filename": "/Users/nikita.galaiko/go/pkg/mod/go.uber.org/dig@v1.14.1/param.go",
          "function": "BuildList",
          "in_app": false,
          "lineno": 151,
          "module": "go.uber.org/dig.paramList",
          "post_context": [
            "\t\tif err != nil {",
            "\t\t\treturn nil, err",
            "\t\t}"
          ],
          "pre_context": [
            "\targs := make([]reflect.Value, len(pl.Params))",
            "\tfor i, p := range pl.Params {",
            "\t\tvar err error"
          ]
        },
        {
          "abs_path": "/Users/nikita.galaiko/go/pkg/mod/go.uber.org/dig@v1.14.1/param.go",
          "context_line": "\t\terr := n.Call(n.OrigScope())",
          "filename": "/Users/nikita.galaiko/go/pkg/mod/go.uber.org/dig@v1.14.1/param.go",
          "function": "Build",
          "in_app": false,
          "lineno": 296,
          "module": "go.uber.org/dig.paramSingle",
          "post_context": [
            "\t\tif err == nil {",
            "\t\t\tcontinue",
            "\t\t}"
          ],
          "pre_context": [
            "\t}",
            "",
            "\tfor _, n := range providers {"
          ]
        },
        {
          "abs_path": "/Users/nikita.galaiko/go/pkg/mod/go.uber.org/dig@v1.14.1/constructor.go",
          "context_line": "\targs, err := n.paramList.BuildList(c)",
          "filename": "/Users/nikita.galaiko/go/pkg/mod/go.uber.org/dig@v1.14.1/constructor.go",
          "function": "Call",
          "in_app": false,
          "lineno": 145,
          "module": "go.uber.org/dig.(*constructorNode)",
          "post_context": [
            "\tif err != nil {",
            "\t\treturn errArgumentsFailed{",
            "\t\t\tFunc:   n.location,"
          ],
          "pre_context": [
            "\t\t}",
            "\t}",
            ""
          ]
        },
        {
          "abs_path": "/Users/nikita.galaiko/go/pkg/mod/go.uber.org/dig@v1.14.1/param.go",
          "context_line": "\t\targs[i], err = p.Build(c)",
          "filename": "/Users/nikita.galaiko/go/pkg/mod/go.uber.org/dig@v1.14.1/param.go",
          "function": "BuildList",
          "in_app": false,
          "lineno": 151,
          "module": "go.uber.org/dig.paramList",
          "post_context": [
            "\t\tif err != nil {",
            "\t\t\treturn nil, err",
            "\t\t}"
          ],
          "pre_context": [
            "\targs := make([]reflect.Value, len(pl.Params))",
            "\tfor i, p := range pl.Params {",
            "\t\tvar err error"
          ]
        },
        {
          "abs_path": "/Users/nikita.galaiko/go/pkg/mod/go.uber.org/dig@v1.14.1/param.go",
          "context_line": "\t\terr := n.Call(n.OrigScope())",
          "filename": "/Users/nikita.galaiko/go/pkg/mod/go.uber.org/dig@v1.14.1/param.go",
          "function": "Build",
          "in_app": false,
          "lineno": 296,
          "module": "go.uber.org/dig.paramSingle",
          "post_context": [
            "\t\tif err == nil {",
            "\t\t\tcontinue",
            "\t\t}"
          ],
          "pre_context": [
            "\t}",
            "",
            "\tfor _, n := range providers {"
          ]
        },
        {
          "abs_path": "/Users/nikita.galaiko/go/pkg/mod/go.uber.org/dig@v1.14.1/constructor.go",
          "context_line": "\targs, err := n.paramList.BuildList(c)",
          "filename": "/Users/nikita.galaiko/go/pkg/mod/go.uber.org/dig@v1.14.1/constructor.go",
          "function": "Call",
          "in_app": false,
          "lineno": 145,
          "module": "go.uber.org/dig.(*constructorNode)",
          "post_context": [
            "\tif err != nil {",
            "\t\treturn errArgumentsFailed{",
            "\t\t\tFunc:   n.location,"
          ],
          "pre_context": [
            "\t\t}",
            "\t}",
            ""
          ]
        },
        {
          "abs_path": "/Users/nikita.galaiko/go/pkg/mod/go.uber.org/dig@v1.14.1/param.go",
          "context_line": "\t\targs[i], err = p.Build(c)",
          "filename": "/Users/nikita.galaiko/go/pkg/mod/go.uber.org/dig@v1.14.1/param.go",
          "function": "BuildList",
          "in_app": false,
          "lineno": 151,
          "module": "go.uber.org/dig.paramList",
          "post_context": [
            "\t\tif err != nil {",
            "\t\t\treturn nil, err",
            "\t\t}"
          ],
          "pre_context": [
            "\targs := make([]reflect.Value, len(pl.Params))",
            "\tfor i, p := range pl.Params {",
            "\t\tvar err error"
          ]
        },
        {
          "abs_path": "/Users/nikita.galaiko/go/pkg/mod/go.uber.org/dig@v1.14.1/param.go",
          "context_line": "\t\terr := n.Call(n.OrigScope())",
          "filename": "/Users/nikita.galaiko/go/pkg/mod/go.uber.org/dig@v1.14.1/param.go",
          "function": "Build",
          "in_app": false,
          "lineno": 296,
          "module": "go.uber.org/dig.paramSingle",
          "post_context": [
            "\t\tif err == nil {",
            "\t\t\tcontinue",
            "\t\t}"
          ],
          "pre_context": [
            "\t}",
            "",
            "\tfor _, n := range providers {"
          ]
        },
        {
          "abs_path": "/Users/nikita.galaiko/go/pkg/mod/go.uber.org/dig@v1.14.1/constructor.go",
          "context_line": "\tresults := c.invoker()(reflect.ValueOf(n.ctor), args)",
          "filename": "/Users/nikita.galaiko/go/pkg/mod/go.uber.org/dig@v1.14.1/constructor.go",
          "function": "Call",
          "in_app": false,
          "lineno": 154,
          "module": "go.uber.org/dig.(*constructorNode)",
          "post_context": [
            "\tif err := n.resultList.ExtractList(receiver, false /* decorating */, results); err != nil {",
            "\t\treturn errConstructorFailed{Func: n.location, Reason: err}",
            "\t}"
          ],
          "pre_context": [
            "\t}",
            "",
            "\treceiver := newStagingContainerWriter()"
          ]
        },
        {
          "abs_path": "/Users/nikita.galaiko/go/pkg/mod/go.uber.org/dig@v1.14.1/container.go",
          "context_line": "\treturn fn.Call(args)",
          "filename": "/Users/nikita.galaiko/go/pkg/mod/go.uber.org/dig@v1.14.1/container.go",
          "function": "defaultInvoker",
          "in_app": false,
          "lineno": 220,
          "module": "go.uber.org/dig",
          "post_context": [
            "}",
            "",
            "// Generates zero values for results without calling the supplied function."
          ],
          "pre_context": [
            "type invokerFn func(fn reflect.Value, args []reflect.Value) (results []reflect.Value)",
            "",
            "func defaultInvoker(fn reflect.Value, args []reflect.Value) []reflect.Value {"
          ]
        },
        {
          "abs_path": "/opt/homebrew/Cellar/go/1.18/libexec/src/reflect/value.go",
          "context_line": "\treturn v.call(\"Call\", in)",
          "filename": "reflect/value.go",
          "function": "Call",
          "in_app": false,
          "lineno": 339,
          "module": "reflect.Value",
          "post_context": [
            "}",
            "",
            "// CallSlice calls the variadic function v with the input arguments in,"
          ],
          "pre_context": [
            "func (v Value) Call(in []Value) []Value {",
            "\tv.mustBe(Func)",
            "\tv.mustBeExported()"
          ]
        },
        {
          "abs_path": "/opt/homebrew/Cellar/go/1.18/libexec/src/reflect/value.go",
          "context_line": "\tcall(frametype, fn, stackArgs, uint32(frametype.size), uint32(abi.retOffset), uint32(frameSize), \u0026regArgs)",
          "filename": "reflect/value.go",
          "function": "call",
          "in_app": false,
          "lineno": 556,
          "module": "reflect.Value",
          "post_context": [
            "",
            "\t// For testing; see TestCallMethodJump.",
            "\tif callGC {"
          ],
          "pre_context": [
            "\t}",
            "",
            "\t// Call."
          ]
        },
        {
          "abs_path": "/Users/nikita.galaiko/src/sturdy/api/pkg/logger/zap.go",
          "context_line": "\tl.Error(\"message\", zap.Error(fmt.Errorf(\"error :(\")))",
          "filename": "/Users/nikita.galaiko/src/sturdy/api/pkg/logger/zap.go",
          "function": "New",
          "in_app": true,
          "lineno": 76,
          "module": "getsturdy.com/api/pkg/logger",
          "post_context": [
            "\treturn l, nil",
            "}",
            ""
          ],
          "pre_context": [
            "\t}",
            "",
            "\tl := zap.New(zapcore.NewTee(cores...), options...)"
          ]
        }
      ]
    },
    "type": "message",
    "value": "error :("
  },
  "level": "error",
  "logger": "root",
  "message": "message",
  "platform": "go",
  "project": "123",
  "release": "development",
  "server_name": "installation-2b5c592f-1f0f-45ba-a433-720b1f3eba8c",
  "timestamp": "2022-04-14T13:01:33.31"
}
`

// this is an example of a payload produced by github.com/getsentry/sentry-go library
const sentryPayload = `
{
  "contexts": {
    "device": {
      "arch": "arm64",
      "num_cpu": 10
    },
    "os": {
      "name": "darwin"
    },
    "runtime": {
      "go_maxprocs": 10,
      "go_numcgocalls": 8,
      "go_numroutines": 4,
      "name": "go",
      "version": "go1.18"
    }
  },
  "environment": "enterprise",
  "event_id": "ce58726b70f64c7c9cb596cc3daf0a07",
  "exception": [
    {
      "stacktrace": {
        "frames": [
          {
            "abs_path": "/Users/nikita.galaiko/src/sturdy/api/cmd/api/main.go",
            "context_line": "\tif err := di.Init(app).To(\u0026apiServer, \u0026ctx); err != nil {",
            "function": "main",
            "in_app": true,
            "lineno": 24,
            "module": "main",
            "post_context": [
              "\t\tfmt.Printf(\"%+v\\n\", err)",
              "\t\tos.Exit(1)",
              "\t}",
              "",
              "\tbanner.PrintBanner()"
            ],
            "pre_context": [
              "\t\tc.Import(xcontext.Module)",
              "\t}",
              "",
              "\tvar apiServer api.Starter",
              "\tvar ctx context.Context"
            ]
          },
          {
            "abs_path": "/Users/nikita.galaiko/src/sturdy/api/pkg/di/container.go",
            "context_line": "\tif err := container.Invoke(invokeFn.Interface()); err != nil {",
            "function": "(*Container).To",
            "in_app": true,
            "lineno": 360,
            "module": "getsturdy.com/api/pkg/di",
            "post_context": [
              "\t\treturn err",
              "\t}",
              "",
              "\treturn nil",
              ""
            ],
            "pre_context": [
              "\t\t\tdp.Elem().Set(args[i])",
              "\t\t}",
              "\t\treturn []reflect.Value{}",
              "\t})",
              ""
            ]
          },
          {
            "abs_path": "/Users/nikita.galaiko/go/pkg/mod/go.uber.org/dig@v1.14.1/invoke.go",
            "context_line": "\treturn c.scope.Invoke(function, opts...)",
            "function": "go.uber.org/dig.(*Container).Invoke",
            "in_app": true,
            "lineno": 46,
            "post_context": [
              "}",
              "",
              "// Invoke runs the given function after instantiating its dependencies.",
              "//",
              "// Any arguments that the function has are treated as its dependencies. The"
            ],
            "pre_context": [
              "// dependencies that they might have.",
              "//",
              "// The function may return an error to indicate failure. The error will be",
              "// returned to the caller as-is.",
              "func (c *Container) Invoke(function interface{}, opts ...InvokeOption) error {"
            ]
          },
          {
            "abs_path": "/Users/nikita.galaiko/go/pkg/mod/go.uber.org/dig@v1.14.1/invoke.go",
            "context_line": "\targs, err := pl.BuildList(s)",
            "function": "go.uber.org/dig.(*Scope).Invoke",
            "in_app": true,
            "lineno": 85,
            "post_context": [
              "\tif err != nil {",
              "\t\treturn errArgumentsFailed{",
              "\t\t\tFunc:   digreflect.InspectFunc(function),",
              "\t\t\tReason: err,",
              "\t\t}"
            ],
            "pre_context": [
              "\t\t\treturn errf(\"cycle detected in dependency graph\", s.cycleDetectedError(cycle))",
              "\t\t}",
              "\t\ts.isVerifiedAcyclic = true",
              "\t}",
              ""
            ]
          },
          {
            "abs_path": "/Users/nikita.galaiko/go/pkg/mod/go.uber.org/dig@v1.14.1/param.go",
            "context_line": "\t\targs[i], err = p.Build(c)",
            "function": "go.uber.org/dig.paramList.BuildList",
            "in_app": true,
            "lineno": 151,
            "post_context": [
              "\t\tif err != nil {",
              "\t\t\treturn nil, err",
              "\t\t}",
              "\t}",
              "\treturn args, nil"
            ],
            "pre_context": [
              "// to the underlying constructor.",
              "func (pl paramList) BuildList(c containerStore) ([]reflect.Value, error) {",
              "\targs := make([]reflect.Value, len(pl.Params))",
              "\tfor i, p := range pl.Params {",
              "\t\tvar err error"
            ]
          },
          {
            "abs_path": "/Users/nikita.galaiko/go/pkg/mod/go.uber.org/dig@v1.14.1/param.go",
            "context_line": "\t\terr := n.Call(n.OrigScope())",
            "function": "go.uber.org/dig.paramSingle.Build",
            "in_app": true,
            "lineno": 296,
            "post_context": [
              "\t\tif err == nil {",
              "\t\t\tcontinue",
              "\t\t}",
              "",
              "\t\t// If we're missing dependencies but the parameter itself is optional,"
            ],
            "pre_context": [
              "\t\t}",
              "\t\treturn _noValue, newErrMissingTypes(c, key{name: ps.Name, t: ps.Type})",
              "\t}",
              "",
              "\tfor _, n := range providers {"
            ]
          },
          {
            "abs_path": "/Users/nikita.galaiko/go/pkg/mod/go.uber.org/dig@v1.14.1/constructor.go",
            "context_line": "\targs, err := n.paramList.BuildList(c)",
            "function": "go.uber.org/dig.(*constructorNode).Call",
            "in_app": true,
            "lineno": 145,
            "post_context": [
              "\tif err != nil {",
              "\t\treturn errArgumentsFailed{",
              "\t\t\tFunc:   n.location,",
              "\t\t\tReason: err,",
              "\t\t}"
            ],
            "pre_context": [
              "\t\t\tFunc:   n.location,",
              "\t\t\tReason: err,",
              "\t\t}",
              "\t}",
              ""
            ]
          },
          {
            "abs_path": "/Users/nikita.galaiko/go/pkg/mod/go.uber.org/dig@v1.14.1/param.go",
            "context_line": "\t\targs[i], err = p.Build(c)",
            "function": "go.uber.org/dig.paramList.BuildList",
            "in_app": true,
            "lineno": 151,
            "post_context": [
              "\t\tif err != nil {",
              "\t\t\treturn nil, err",
              "\t\t}",
              "\t}",
              "\treturn args, nil"
            ],
            "pre_context": [
              "// to the underlying constructor.",
              "func (pl paramList) BuildList(c containerStore) ([]reflect.Value, error) {",
              "\targs := make([]reflect.Value, len(pl.Params))",
              "\tfor i, p := range pl.Params {",
              "\t\tvar err error"
            ]
          },
          {
            "abs_path": "/Users/nikita.galaiko/go/pkg/mod/go.uber.org/dig@v1.14.1/param.go",
            "context_line": "\t\terr := n.Call(n.OrigScope())",
            "function": "go.uber.org/dig.paramSingle.Build",
            "in_app": true,
            "lineno": 296,
            "post_context": [
              "\t\tif err == nil {",
              "\t\t\tcontinue",
              "\t\t}",
              "",
              "\t\t// If we're missing dependencies but the parameter itself is optional,"
            ],
            "pre_context": [
              "\t\t}",
              "\t\treturn _noValue, newErrMissingTypes(c, key{name: ps.Name, t: ps.Type})",
              "\t}",
              "",
              "\tfor _, n := range providers {"
            ]
          },
          {
            "abs_path": "/Users/nikita.galaiko/go/pkg/mod/go.uber.org/dig@v1.14.1/constructor.go",
            "context_line": "\targs, err := n.paramList.BuildList(c)",
            "function": "go.uber.org/dig.(*constructorNode).Call",
            "in_app": true,
            "lineno": 145,
            "post_context": [
              "\tif err != nil {",
              "\t\treturn errArgumentsFailed{",
              "\t\t\tFunc:   n.location,",
              "\t\t\tReason: err,",
              "\t\t}"
            ],
            "pre_context": [
              "\t\t\tFunc:   n.location,",
              "\t\t\tReason: err,",
              "\t\t}",
              "\t}",
              ""
            ]
          },
          {
            "abs_path": "/Users/nikita.galaiko/go/pkg/mod/go.uber.org/dig@v1.14.1/param.go",
            "context_line": "\t\targs[i], err = p.Build(c)",
            "function": "go.uber.org/dig.paramList.BuildList",
            "in_app": true,
            "lineno": 151,
            "post_context": [
              "\t\tif err != nil {",
              "\t\t\treturn nil, err",
              "\t\t}",
              "\t}",
              "\treturn args, nil"
            ],
            "pre_context": [
              "// to the underlying constructor.",
              "func (pl paramList) BuildList(c containerStore) ([]reflect.Value, error) {",
              "\targs := make([]reflect.Value, len(pl.Params))",
              "\tfor i, p := range pl.Params {",
              "\t\tvar err error"
            ]
          },
          {
            "abs_path": "/Users/nikita.galaiko/go/pkg/mod/go.uber.org/dig@v1.14.1/param.go",
            "context_line": "\t\terr := n.Call(n.OrigScope())",
            "function": "go.uber.org/dig.paramSingle.Build",
            "in_app": true,
            "lineno": 296,
            "post_context": [
              "\t\tif err == nil {",
              "\t\t\tcontinue",
              "\t\t}",
              "",
              "\t\t// If we're missing dependencies but the parameter itself is optional,"
            ],
            "pre_context": [
              "\t\t}",
              "\t\treturn _noValue, newErrMissingTypes(c, key{name: ps.Name, t: ps.Type})",
              "\t}",
              "",
              "\tfor _, n := range providers {"
            ]
          },
          {
            "abs_path": "/Users/nikita.galaiko/go/pkg/mod/go.uber.org/dig@v1.14.1/constructor.go",
            "context_line": "\targs, err := n.paramList.BuildList(c)",
            "function": "go.uber.org/dig.(*constructorNode).Call",
            "in_app": true,
            "lineno": 145,
            "post_context": [
              "\tif err != nil {",
              "\t\treturn errArgumentsFailed{",
              "\t\t\tFunc:   n.location,",
              "\t\t\tReason: err,",
              "\t\t}"
            ],
            "pre_context": [
              "\t\t\tFunc:   n.location,",
              "\t\t\tReason: err,",
              "\t\t}",
              "\t}",
              ""
            ]
          },
          {
            "abs_path": "/Users/nikita.galaiko/go/pkg/mod/go.uber.org/dig@v1.14.1/param.go",
            "context_line": "\t\targs[i], err = p.Build(c)",
            "function": "go.uber.org/dig.paramList.BuildList",
            "in_app": true,
            "lineno": 151,
            "post_context": [
              "\t\tif err != nil {",
              "\t\t\treturn nil, err",
              "\t\t}",
              "\t}",
              "\treturn args, nil"
            ],
            "pre_context": [
              "// to the underlying constructor.",
              "func (pl paramList) BuildList(c containerStore) ([]reflect.Value, error) {",
              "\targs := make([]reflect.Value, len(pl.Params))",
              "\tfor i, p := range pl.Params {",
              "\t\tvar err error"
            ]
          },
          {
            "abs_path": "/Users/nikita.galaiko/go/pkg/mod/go.uber.org/dig@v1.14.1/param.go",
            "context_line": "\t\terr := n.Call(n.OrigScope())",
            "function": "go.uber.org/dig.paramSingle.Build",
            "in_app": true,
            "lineno": 296,
            "post_context": [
              "\t\tif err == nil {",
              "\t\t\tcontinue",
              "\t\t}",
              "",
              "\t\t// If we're missing dependencies but the parameter itself is optional,"
            ],
            "pre_context": [
              "\t\t}",
              "\t\treturn _noValue, newErrMissingTypes(c, key{name: ps.Name, t: ps.Type})",
              "\t}",
              "",
              "\tfor _, n := range providers {"
            ]
          },
          {
            "abs_path": "/Users/nikita.galaiko/go/pkg/mod/go.uber.org/dig@v1.14.1/constructor.go",
            "context_line": "\targs, err := n.paramList.BuildList(c)",
            "function": "go.uber.org/dig.(*constructorNode).Call",
            "in_app": true,
            "lineno": 145,
            "post_context": [
              "\tif err != nil {",
              "\t\treturn errArgumentsFailed{",
              "\t\t\tFunc:   n.location,",
              "\t\t\tReason: err,",
              "\t\t}"
            ],
            "pre_context": [
              "\t\t\tFunc:   n.location,",
              "\t\t\tReason: err,",
              "\t\t}",
              "\t}",
              ""
            ]
          },
          {
            "abs_path": "/Users/nikita.galaiko/go/pkg/mod/go.uber.org/dig@v1.14.1/param.go",
            "context_line": "\t\targs[i], err = p.Build(c)",
            "function": "go.uber.org/dig.paramList.BuildList",
            "in_app": true,
            "lineno": 151,
            "post_context": [
              "\t\tif err != nil {",
              "\t\t\treturn nil, err",
              "\t\t}",
              "\t}",
              "\treturn args, nil"
            ],
            "pre_context": [
              "// to the underlying constructor.",
              "func (pl paramList) BuildList(c containerStore) ([]reflect.Value, error) {",
              "\targs := make([]reflect.Value, len(pl.Params))",
              "\tfor i, p := range pl.Params {",
              "\t\tvar err error"
            ]
          },
          {
            "abs_path": "/Users/nikita.galaiko/go/pkg/mod/go.uber.org/dig@v1.14.1/param.go",
            "context_line": "\t\terr := n.Call(n.OrigScope())",
            "function": "go.uber.org/dig.paramSingle.Build",
            "in_app": true,
            "lineno": 296,
            "post_context": [
              "\t\tif err == nil {",
              "\t\t\tcontinue",
              "\t\t}",
              "",
              "\t\t// If we're missing dependencies but the parameter itself is optional,"
            ],
            "pre_context": [
              "\t\t}",
              "\t\treturn _noValue, newErrMissingTypes(c, key{name: ps.Name, t: ps.Type})",
              "\t}",
              "",
              "\tfor _, n := range providers {"
            ]
          },
          {
            "abs_path": "/Users/nikita.galaiko/go/pkg/mod/go.uber.org/dig@v1.14.1/constructor.go",
            "context_line": "\targs, err := n.paramList.BuildList(c)",
            "function": "go.uber.org/dig.(*constructorNode).Call",
            "in_app": true,
            "lineno": 145,
            "post_context": [
              "\tif err != nil {",
              "\t\treturn errArgumentsFailed{",
              "\t\t\tFunc:   n.location,",
              "\t\t\tReason: err,",
              "\t\t}"
            ],
            "pre_context": [
              "\t\t\tFunc:   n.location,",
              "\t\t\tReason: err,",
              "\t\t}",
              "\t}",
              ""
            ]
          },
          {
            "abs_path": "/Users/nikita.galaiko/go/pkg/mod/go.uber.org/dig@v1.14.1/param.go",
            "context_line": "\t\targs[i], err = p.Build(c)",
            "function": "go.uber.org/dig.paramList.BuildList",
            "in_app": true,
            "lineno": 151,
            "post_context": [
              "\t\tif err != nil {",
              "\t\t\treturn nil, err",
              "\t\t}",
              "\t}",
              "\treturn args, nil"
            ],
            "pre_context": [
              "// to the underlying constructor.",
              "func (pl paramList) BuildList(c containerStore) ([]reflect.Value, error) {",
              "\targs := make([]reflect.Value, len(pl.Params))",
              "\tfor i, p := range pl.Params {",
              "\t\tvar err error"
            ]
          },
          {
            "abs_path": "/Users/nikita.galaiko/go/pkg/mod/go.uber.org/dig@v1.14.1/param.go",
            "context_line": "\t\terr := n.Call(n.OrigScope())",
            "function": "go.uber.org/dig.paramSingle.Build",
            "in_app": true,
            "lineno": 296,
            "post_context": [
              "\t\tif err == nil {",
              "\t\t\tcontinue",
              "\t\t}",
              "",
              "\t\t// If we're missing dependencies but the parameter itself is optional,"
            ],
            "pre_context": [
              "\t\t}",
              "\t\treturn _noValue, newErrMissingTypes(c, key{name: ps.Name, t: ps.Type})",
              "\t}",
              "",
              "\tfor _, n := range providers {"
            ]
          },
          {
            "abs_path": "/Users/nikita.galaiko/go/pkg/mod/go.uber.org/dig@v1.14.1/constructor.go",
            "context_line": "\tresults := c.invoker()(reflect.ValueOf(n.ctor), args)",
            "function": "go.uber.org/dig.(*constructorNode).Call",
            "in_app": true,
            "lineno": 154,
            "post_context": [
              "\tif err := n.resultList.ExtractList(receiver, false /* decorating */, results);err != nil {",
              "\t\treturn errConstructorFailed{Func: n.location, Reason: err}",
              "\t}",
              "",
              "\t// Commit the result to the original container that this constructor"
            ],
            "pre_context": [
              "\t\t\tReason: err,",
              "\t\t}",
              "\t}",
              "",
              "\treceiver := newStagingContainerWriter()"
            ]
          },
          {
            "abs_path": "/Users/nikita.galaiko/go/pkg/mod/go.uber.org/dig@v1.14.1/container.go",
            "context_line": "\treturn fn.Call(args)",
            "function": "go.uber.org/dig.defaultInvoker",
            "in_app": true,
            "lineno": 220,
            "post_context": [
              "}",
              "",
              "// Generates zero values for results without calling the supplied function.",
              "func dryInvoker(fn reflect.Value, _ []reflect.Value) []reflect.Value {",
              "\tft := fn.Type()"
            ],
            "pre_context": [
              "",
              "// invokerFn specifies how the container calls user-supplied functions.",
              "type invokerFn func(fn reflect.Value, args []reflect.Value) (results []reflect.Value)",
              "",
              "func defaultInvoker(fn reflect.Value, args []reflect.Value) []reflect.Value {"
            ]
          },
          {
            "abs_path": "/opt/homebrew/Cellar/go/1.18/libexec/src/reflect/value.go",
            "function": "Value.Call",
            "lineno": 339,
            "module": "reflect"
          },
          {
            "abs_path": "/opt/homebrew/Cellar/go/1.18/libexec/src/reflect/value.go",
            "function": "Value.call",
            "lineno": 556,
            "module": "reflect"
          },
          {
            "abs_path": "/Users/nikita.galaiko/src/sturdy/api/pkg/logger/zap.go",
            "context_line": "\tl.Error(\"message\", zap.Error(fmt.Errorf(\"error :(\")))",
            "function": "New",
            "in_app": true,
            "lineno": 85,
            "module": "getsturdy.com/api/pkg/logger",
            "post_context": [
              "\treturn l, nil",
              "}",
              ""
            ],
            "pre_context": [
              "\t\t}",
              "\t\tcores = append(cores, core)",
              "\t}",
              "",
              "\tl := zap.New(zapcore.NewTee(cores...), options...)"
            ]
          }
        ]
      },
      "type": "*errors.errorString",
      "value": "error :("
    }
  ],
  "extra": {
    "error": "error :("
  },
  "level": "error",
  "message": "message",
  "modules": {
    "getsturdy.com/api": "(devel)",
    "github.com/ProtonMail/go-crypto": "v0.0.0-20210428141323-04723f9f07d7",
    "github.com/ScaleFT/sshkeys": "v0.0.0-20200327173127-6142f742bca5",
    "github.com/TheZeroSlave/zapsentry": "v1.10.0",
    "github.com/aws/aws-sdk-go": "v1.38.47",
    "github.com/aymerick/douceur": "v0.2.0",
    "github.com/beorn7/perks": "v1.0.1",
    "github.com/bmatcuk/doublestar/v4": "v4.0.2",
    "github.com/bradleyfalzon/ghinstallation": "v1.1.1",
    "github.com/buildkite/go-buildkite/v3": "v3.0.0",
    "github.com/cenkalti/backoff": "v2.0.0+incompatible",
    "github.com/cenkalti/backoff/v4": "v4.1.2",
    "github.com/cespare/xxhash/v2": "v2.1.1",
    "github.com/dchest/bcrypt_pbkdf": "v0.0.0-20150205184540-83f37f9c154a",
    "github.com/dgrijalva/jwt-go": "v3.2.0+incompatible",
    "github.com/disintegration/imaging": "v1.6.2",
    "github.com/emirpasic/gods": "v1.12.0",
    "github.com/fatih/color": "v1.13.0",
    "github.com/getsentry/sentry-go": "v0.13.0",
    "github.com/gin-contrib/cors": "v1.3.1",
    "github.com/gin-contrib/gzip": "v0.0.5",
    "github.com/gin-contrib/sse": "v0.1.0",
    "github.com/gin-contrib/zap": "v0.0.2",
    "github.com/gin-gonic/gin": "v1.7.7",
    "github.com/go-git/gcfg": "v1.5.0",
    "github.com/go-git/go-billy/v5": "v5.3.1",
    "github.com/go-git/go-git/v5": "v5.4.3 =\u003e github.com/zegl/go-git/v5 v5.4.3-0.20220401122347-e4c6e92beccd",
    "github.com/go-playground/locales": "v0.14.0",
    "github.com/go-playground/universal-translator": "v0.18.0",
    "github.com/go-playground/validator/v10": "v10.10.0",
    "github.com/gofrs/flock": "v0.8.1",
    "github.com/golang-migrate/migrate/v4": "v4.15.1",
    "github.com/golang/protobuf": "v1.5.2",
    "github.com/google/go-github/v29": "v29.0.2",
    "github.com/google/go-github/v39": "v39.2.0",
    "github.com/google/go-querystring": "v1.1.0",
    "github.com/google/uuid": "v1.3.0",
    "github.com/gorilla/css": "v1.0.0",
    "github.com/gorilla/websocket": "v1.4.2",
    "github.com/gosimple/slug": "v1.9.0",
    "github.com/graph-gophers/dataloader/v6": "v6.0.0",
    "github.com/graph-gophers/graphql-go": "v1.3.0",
    "github.com/graph-gophers/graphql-transport-ws": "v0.0.1 =\u003e github.com/sturdy-dev/graphql-transport-ws v0.0.0-20211122094650-15c742155db6",
    "github.com/h2non/filetype": "v1.1.3",
    "github.com/hashicorp/errwrap": "v1.1.0",
    "github.com/hashicorp/go-multierror": "v1.1.1",
    "github.com/hashicorp/golang-lru": "v0.5.4",
    "github.com/imdario/mergo": "v0.3.12",
    "github.com/jbenet/go-context": "v0.0.0-20150711004518-d14ea06fba99",
    "github.com/jessevdk/go-flags": "v1.5.0 =\u003e github.com/sturdy-dev/go-flags v1.5.1-0.20220203104421-967e8bff1baf",
    "github.com/jmespath/go-jmespath": "v0.4.0",
    "github.com/jmoiron/sqlx": "v1.3.4",
    "github.com/jxskiss/base62": "v0.0.0-20191017122030-4f11678b909b",
    "github.com/keighl/postmark": "v0.0.0-20190821160221-28358b1a94e3 =\u003e github.com/sturdy-dev/postmark v0.0.0-20220413131856-fc6a9ecca126",
    "github.com/kevinburke/ssh_config": "v0.0.0-20201106050909-4977a11b4351",
    "github.com/leodido/go-urn": "v1.2.1",
    "github.com/lib/pq": "v1.10.4",
    "github.com/libgit2/git2go/v33": "v33.0.0",
    "github.com/mattn/go-colorable": "v0.1.11",
    "github.com/mattn/go-isatty": "v0.0.14",
    "github.com/matttproud/golang_protobuf_extensions": "v1.0.2-0.20181231171920-c182affec369",
    "github.com/microcosm-cc/bluemonday": "v1.0.16",
    "github.com/mitchellh/go-homedir": "v1.1.0",
    "github.com/opentracing/opentracing-go": "v1.2.0",
    "github.com/pkg/errors": "v0.9.1",
    "github.com/posthog/posthog-go": "v0.0.0-20211028072449-93c17c49e2b0",
    "github.com/prometheus/client_golang": "v1.11.0",
    "github.com/prometheus/client_model": "v0.2.0",
    "github.com/prometheus/common": "v0.26.0",
    "github.com/prometheus/procfs": "v0.6.0",
    "github.com/rainycape/unidecode": "v0.0.0-20150907023854-cb7f23ec59be",
    "github.com/sergi/go-diff": "v1.1.0",
    "github.com/sourcegraph/go-diff": "v0.6.2-0.20210526090523-35b24a7eb480 =\u003e github.com/ngalaiko/go-diff v0.6.2-0.20220224161118-fbc7fabee1d1",
    "github.com/tailscale/hujson": "v0.0.0-20210818175511-7360507a6e88",
    "github.com/tidwall/match": "v1.0.3",
    "github.com/ugorji/go/codec": "v1.1.7",
    "github.com/xanzy/ssh-agent": "v0.3.1",
    "github.com/xtgo/uuid": "v0.0.0-20140804021211-a0b114877d4c",
    "github.com/yuin/goldmark": "v1.4.4",
    "go.uber.org/atomic": "v1.7.0",
    "go.uber.org/dig": "v1.14.1",
    "go.uber.org/multierr": "v1.7.0",
    "go.uber.org/zap": "v1.21.0",
    "golang.org/x/crypto": "v0.0.0-20220112180741-5e0467b6c7ce",
    "golang.org/x/image": "v0.0.0-20210216034530-4410531fe030",
    "golang.org/x/net": "v0.0.0-20211112202133-69e39bad7dc2",
    "golang.org/x/oauth2": "v0.0.0-20210628180205-a41e5a781914",
    "golang.org/x/sync": "v0.0.0-20210220032951-036812b2e83c",
    "golang.org/x/sys": "v0.0.0-20211025201205-69cdffdb9359",
    "golang.org/x/term": "v0.0.0-20210927222741-03fcf44c2211",
    "golang.org/x/text": "v0.3.7",
    "google.golang.org/protobuf": "v1.27.1",
    "gopkg.in/square/go-jose.v2": "v2.6.0",
    "gopkg.in/warnings.v0": "v0.1.2",
    "gopkg.in/yaml.v2": "v2.4.0"
  },
  "platform": "go",
  "release": "development",
  "sdk": {
    "integrations": [
      "ContextifyFrames",
      "Environment",
      "IgnoreErrors",
      "Modules"
    ],
    "name": "sentry.go",
    "packages": [
      {
        "name": "sentry-go",
        "version": "0.13.0"
      }
    ],
    "version": "0.13.0"
  },
  "server_name": "installation-2b5c592f-1f0f-45ba-a433-720b1f3eba8c",
  "timestamp": "2022-04-14T14:57:00.065525+02:00",
  "user": {}
}`
