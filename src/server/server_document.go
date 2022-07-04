package server

import (
	"context"
	"github.com/36625090/involution/logical"
	"github.com/36625090/involution/ui"
	assetfs "github.com/elazarl/go-bindata-assetfs"
	"github.com/gin-gonic/gin"
	"net/http"
	"path/filepath"
)

func (m *Server) addDocumentSchema() {
	documentResponse := make([]logical.DocumentResponse, 0)
	for name, backend := range m.backends {
		reply, _ := backend.Documents(context.Background())
		resp := logical.DocumentResponse{
			Name:      backend.BackendDescription(),
			Backend:   name,
			Documents: reply.Documents,
		}
		documentResponse = append(documentResponse, resp)
	}

	path := filepath.Join(m.opts.Http.Path, "_m", "schemas")
	m.httpTransport.GET(path, func(c *gin.Context) {
		c.SecureJSON(200, map[string]interface{}{
			"code":   0,
			"result": documentResponse,
		})
	})

}

func (m *Server) addDocumentUI() {
	fs := assetfs.AssetFS{Asset: ui.Asset, AssetDir: ui.AssetDir,
		Prefix: "docs", Fallback: "index.html"}
	path := filepath.Join(m.opts.Http.Path, "docs")
	m.logger.Info("initialize handle", "path", path)
	urlPattern := filepath.Join(path, "/*filepath")
	handle := createStaticHandler(path, &fs)

	m.httpTransport.GET(urlPattern, handle)

	m.httpTransport.HEAD(urlPattern, handle)
}

func createStaticHandler(relativePath string, fs http.FileSystem) gin.HandlerFunc {
	fileServer := http.StripPrefix(relativePath, http.FileServer(fs))

	return func(c *gin.Context) {
		if _, nolisting := fs.(http.FileSystem); nolisting {
			c.Writer.WriteHeader(http.StatusNotFound)
		}

		file := c.Param("filepath")

		// Check if file exists and/or if we have permission to access it
		f, err := fs.Open(file)
		if err != nil {
			if file != "/" {
				c.Redirect(302, "index.html?api_url="+c.Query("api_url"))
				return
			}
			c.Writer.WriteHeader(http.StatusNotFound)
			return
		}
		f.Close()
		fileServer.ServeHTTP(c.Writer, c.Request)
	}
}
