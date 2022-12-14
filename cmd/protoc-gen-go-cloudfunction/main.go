package main

import (
	"strings"

	"google.golang.org/protobuf/compiler/protogen"
)

func main() {

	protogen.Options{}.Run(func(gen *protogen.Plugin) error {
		for _, f := range gen.Files {
			if !f.Generate {
				continue
			}
			generateFile(gen, f)
		}
		return nil
	})
}

// generateFile generates a _ascii.pb.go file containing gRPC service definitions.
func generateFile(gen *protogen.Plugin, file *protogen.File) {

	var g *protogen.GeneratedFile
	isFirst := true
	for _, msg := range file.Messages {
		if strings.HasSuffix(string(msg.Desc.Name()), "Request") {
			if isFirst {
				filename := file.GeneratedFilenamePrefix + "_cloudfunction.pb.go"
				g = gen.NewGeneratedFile(filename, file.GoImportPath)

				g.P("// Code generated by protoc-gen-go-cloudfunction. DO NOT EDIT.")
				g.P()
				g.P("package ", file.GoPackageName)
				g.P("import (")
				g.P("\"context\"")
				g.P("\"net/http\"")
				g.P("\"github.com/aiocean/cfutil\"")
				g.P(")")
				g.P()
				isFirst = false
			}

			var functionName = strings.TrimSuffix(string(msg.Desc.Name()), "Request")
			g.P("type ", functionName, "HandlerFunc = func(context.Context, *", functionName, "Request) (*", functionName, "Response, error)")
			g.P("func ", functionName, "Handler(w http.ResponseWriter, r *http.Request, do ", functionName, "HandlerFunc) {")
			g.P("if err := cfutil.ApplyCors(w, r); err != nil { cfutil.WriteError(w, r, http.StatusInternalServerError, err); return}")
			g.P("if err := cfutil.ApplyContentType(w, r); err != nil { cfutil.WriteError(w, r, http.StatusInternalServerError, err); return}")
			g.P("var request ", functionName, "Request")
			g.P("if err := cfutil.ReadRequest(r, &request); err != nil { cfutil.WriteError(w, r, http.StatusBadRequest, err); return}")
			g.P("response, err := do(r.Context(), &request)")
			g.P("if err != nil { cfutil.WriteError(w, r, http.StatusInternalServerError, err); return}")
			g.P("if err := cfutil.WriteResponse(w, r, response); err != nil {cfutil.WriteError(w, r, http.StatusInternalServerError, err)}")
			g.P("}")
		}
	}
}
