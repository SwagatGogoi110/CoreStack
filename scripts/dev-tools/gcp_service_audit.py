import os
import re

def audit_gcp_services():
    services_path = "internal/services/gcp/"
    if not os.path.exists(services_path):
        print(f"Directory {services_path} not found.")
        return

    services = sorted([d for d in os.listdir(services_path) if os.path.isdir(os.path.join(services_path, d))])
    
    print(f"{'GCP Service':<25} | {'Status':<15} | {'Notes'}")
    print("-" * 80)
    
    for svc in services:
        handler_path = os.path.join(services_path, svc, "handler.go")
        if not os.path.exists(handler_path):
            print(f"{svc:<25} | No Handler     | -")
            continue
            
        with open(handler_path, "r") as f:
            content = f.read()
            
        # Heuristic for implementation status
        if "w.WriteHeader(http.StatusNotImplemented)" in content and "Not implemented" in content:
            status = "Stub"
        else:
            status = "Partial/Implemented"
            
        print(f"{svc:<25} | {status:<15} | Scaffolded")

if __name__ == "__main__":
    audit_gcp_services()
