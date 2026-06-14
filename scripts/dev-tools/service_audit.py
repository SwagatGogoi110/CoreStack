import os
import re

def audit_services():
    services_path = "internal/services/aws/"
    services = sorted([d for d in os.listdir(services_path) if os.path.isdir(os.path.join(services_path, d))])
    
    print(f"{'Service':<25} | {'Status':<15} | {'Implemented Actions'}")
    print("-" * 80)
    
    for svc in services:
        handler_path = os.path.join(services_path, svc, "handler.go")
        if not os.path.exists(handler_path):
            print(f"{svc:<25} | No Handler     | -")
            continue
            
        with open(handler_path, "r") as f:
            content = f.read()
            
        # Look for actions in switch cases
        actions = re.findall(r'case "(\w+)"', content)
        
        # Heuristic for implementation status
        is_stub = "UnknownOperationException" in content and len(actions) <= 1
        if "HandleJSON" in content and "HandleQuery" in content:
             # Check if both just return errors
             json_unsupported = "does not support JSON protocol" in content or "not supported" in content
             query_unsupported = "does not support Query protocol" in content or "not supported" in content
             if json_unsupported and query_unsupported and len(actions) == 0:
                 status = "Stub"
             elif len(actions) > 0:
                 status = "Partial" if is_stub else "Implemented"
             else:
                 status = "Stub"
        else:
            status = "Implemented" if len(actions) > 2 else "Partial"
            
        action_str = ", ".join(actions[:5]) + ("..." if len(actions) > 5 else "")
        print(f"{svc:<25} | {status:<15} | {action_str or 'None'}")

if __name__ == "__main__":
    audit_services()
