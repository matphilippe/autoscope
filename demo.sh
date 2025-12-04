#!/bin/bash

# Demo script for svscope

echo "=== SVScope Demo ==="
echo ""

# Create a test directory structure
echo "Creating test directory structure..."
mkdir -p demo_repo/{src/api,src/db,modules/auth,modules/cache}

cd demo_repo

# Initialize git repo
git init -q

# Create config
cat > .svscope.yaml << 'EOF'
modules:
  - name: api
    files: src/api/**/*
  - name: db
    files: src/db/**/*
  - filesRe: modules/(?P<scope>\w+)/.*
EOF

echo "Config created:"
cat .svscope.yaml
echo ""

# Create some files
echo "package api" > src/api/handler.go
echo "package db" > src/db/query.go
echo "# Auth module" > modules/auth/main.tf

# Stage files
git add src/api/handler.go

echo "Staged file: src/api/handler.go"
echo ""
echo "Original commit message: 'fix: fixed the bug'"
echo ""

# Test the tool
echo "fix: fixed the bug" > test_commit.txt
../svscope test_commit.txt

echo "Modified commit message:"
cat test_commit.txt
echo ""

# Test with multiple files
git add src/db/query.go modules/auth/main.tf

echo "Staged files: src/api/handler.go, src/db/query.go, modules/auth/main.tf"
echo ""
echo "Original commit message: 'feat: add new feature'"
echo ""

echo "feat: add new feature" > test_commit2.txt
../svscope test_commit2.txt

echo "Modified commit message:"
cat test_commit2.txt
echo ""

# Cleanup
cd ..
rm -rf demo_repo

echo "=== Demo Complete ==="
