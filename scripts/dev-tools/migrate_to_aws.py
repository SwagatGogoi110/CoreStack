import os
import shutil
import glob

# Create multi-cloud dirs
os.makedirs("internal/services/aws", exist_ok=True)
os.makedirs("internal/services/gcp", exist_ok=True)
os.makedirs("internal/services/azure", exist_ok=True)

# Move AWS services
services_dir = "internal/services"
for item in os.listdir(services_dir):
    if item in ["aws", "gcp", "azure"]:
        continue
    src = os.path.join(services_dir, item)
    if os.path.isdir(src):
        dst = os.path.join("internal/services/aws", item)
        shutil.move(src, dst)

print("Moved existing services to internal/services/aws/")

# Update imports in Go files
old_import = "github.com/hectorvent/cloudstack/internal/services/"
new_import = "github.com/hectorvent/cloudstack/internal/services/aws/"

def replace_in_file(filepath):
    with open(filepath, 'r') as f:
        content = f.read()
    if old_import in content:
        # Don't accidentally replace if it's already updated, though Python replace is dumb.
        # It's better to make sure it's not replacing `internal/services/aws/` with `internal/services/aws/aws/`
        content = content.replace(new_import, old_import) # Revert any accidental double applying
        content = content.replace(old_import, new_import)
        with open(filepath, 'w') as f:
            f.write(content)
        print(f"Updated {filepath}")

# Update all .go files in the project
for root, _, files in os.walk("."):
    # skip .git and vendor
    if ".git" in root or "vendor" in root or "node_modules" in root:
        continue
    for file in files:
        if file.endswith(".go"):
            replace_in_file(os.path.join(root, file))

print("Import paths updated successfully.")
