import os
import sys

def scaffold_gcp_service(key, name):
    os.makedirs(f"internal/services/gcp/{key}/model", exist_ok=True)
    
    with open(f"internal/services/gcp/{key}/model/models.go", "w") as f:
        f.write(f'package model\n\ntype {name}Resource struct {{\n\tId string `json:"name"`\n}}\n')
        
    with open(f"internal/services/gcp/{key}/service.go", "w") as f:
        f.write(f"""package {key}

import (
	"context"
	"github.com/hectorvent/cloudstack/internal/storage"
)

type {name}Service struct {{
}}

func New{name}Service(factory *storage.Factory) (*{name}Service, error) {{
	return &{name}Service{{}}, nil
}}
""")

    with open(f"internal/services/gcp/{key}/handler.go", "w") as f:
        f.write(f"""package {key}

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/hectorvent/cloudstack/internal/core/common"
)

type {name}Handler struct {{
	service *{name}Service
}}

func New{name}Handler(service *{name}Service) *{name}Handler {{
	return &{name}Handler{{
		service: service,
	}}
}}

func (h *{name}Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {{
	// Basic routing example for GCP REST API
	path := r.URL.Path
	
	if r.Method == "GET" && strings.HasSuffix(path, "/{key}s") {{
		// Example: List
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{{"items": []any{{}}}})
		return
	}}

	w.WriteHeader(http.StatusNotImplemented)
	fmt.Fprintf(w, `{{\"error\": {{\"code\": 501, \"message\": \"Not implemented: %s %s\"}}}}`, r.Method, path)
}}
""")

    os.makedirs("docs/gcp-services", exist_ok=True)
    with open(f"docs/gcp-services/{key}.md", "w") as f:
        f.write(f"# {name} (GCP)\n\nCurrently a stub for the {name} service.\n")

if __name__ == "__main__":
    if len(sys.argv) < 3:
        print("Usage: python3 scaffold_gcp_service.py <key> <Name>")
        sys.exit(1)
    scaffold_gcp_service(sys.argv[1], sys.argv[2])
