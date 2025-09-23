#!/usr/bin/env python3

import yaml
import os
import sys
from pathlib import Path

def validate_yaml_file(file_path):
    """Validate a single YAML file"""
    try:
        with open(file_path, 'r', encoding='utf-8') as file:
            yaml.safe_load_all(file)
        return True, None
    except yaml.YAMLError as e:
        return False, str(e)
    except Exception as e:
        return False, str(e)

def main():
    print("🔍 Validating Kubernetes YAML files...")
    
    # Get all YAML files
    yaml_files = []
    for root, dirs, files in os.walk('.'):
        for file in files:
            if file.endswith(('.yaml', '.yml')):
                yaml_files.append(os.path.join(root, file))
    
    if not yaml_files:
        print("❌ No YAML files found!")
        sys.exit(1)
    
    print(f"📋 Found {len(yaml_files)} YAML files")
    
    errors = []
    success_count = 0
    
    for yaml_file in sorted(yaml_files):
        print(f"  ✓ Validating {yaml_file}...")
        is_valid, error = validate_yaml_file(yaml_file)
        
        if is_valid:
            success_count += 1
        else:
            errors.append(f"{yaml_file}: {error}")
    
    print(f"\n📊 Results:")
    print(f"  ✅ Valid files: {success_count}")
    print(f"  ❌ Invalid files: {len(errors)}")
    
    if errors:
        print(f"\n❌ Errors found:")
        for error in errors:
            print(f"  - {error}")
        sys.exit(1)
    else:
        print(f"\n✅ All YAML files are valid!")
        print(f"\n🚀 Ready for Kubernetes deployment!")
        
        # Count resources
        resource_counts = {
            'Namespace': 0,
            'ConfigMap': 0,
            'Secret': 0,
            'StatefulSet': 0,
            'Deployment': 0,
            'Service': 0,
            'Ingress': 0,
            'NetworkPolicy': 0,
            'Job': 0
        }
        
        for yaml_file in yaml_files:
            with open(yaml_file, 'r', encoding='utf-8') as file:
                for doc in yaml.safe_load_all(file):
                    if doc and 'kind' in doc:
                        kind = doc['kind']
                        if kind in resource_counts:
                            resource_counts[kind] += 1
        
        print(f"\n📋 Resource summary:")
        for kind, count in resource_counts.items():
            if count > 0:
                print(f"  - {kind}: {count}")

if __name__ == "__main__":
    main()
