#!/bin/sh
files=$(git diff --cached --name-only)
msg=$(cat "${1}")
autoscope "${msg}" $files
