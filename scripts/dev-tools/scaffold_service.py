import os

def scaffold_service(key, name):
    os.makedirs(f"internal/services/{key}/model", exist_ok=True)
    
    with open(f"internal/services/{key}/model/models.go", "w") as f:
        f.write('package model\n\ntype Resource struct {\n\tId string `json:"Id"`\n}\n')
        
    with open(f"internal/services/{key}/service.go", "w") as f:
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

    with open(f"internal/services/{key}/handler.go", "w") as f:
        f.write(f"""package {key}

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/hectorvent/cloudstack/internal/core/common"
)

type {name}JsonHandler struct {{
	service *{name}Service
}}

func New{name}JsonHandler(service *{name}Service) *{name}JsonHandler {{
	return &{name}JsonHandler{{
		service: service,
	}}
}}

func (h *{name}JsonHandler) HandleJSON(action string, request json.RawMessage, rc *common.RequestContext) (any, error) {{
	return nil, fmt.Errorf("UnknownOperationException: Operation %s is not supported", action)
}}

func (h *{name}JsonHandler) HandleQuery(action string, params url.Values, rc *common.RequestContext) (string, error) {{
	return "", fmt.Errorf("{name} does not support Query protocol")
}}
""")

    with open(f"docs/services/{key}.md", "w") as f:
        f.write(f"# {name}\n\nCurrently a stub for the {name} service.\n")

if __name__ == "__main__":
    import sys
    if len(sys.argv) < 3:
        print("Usage: python3 scaffold_service.py <key> <Name>")
        sys.exit(1)
    scaffold_service(sys.argv[1], sys.argv[2])
