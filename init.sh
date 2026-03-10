#!/bin/bash
set -e

echo "Waiting for Ollama..."
until ollama list > /dev/null 2>&1; do
  sleep 1
done

echo "Ollama ready. Checking model..."
if ! ollama list | grep -q "dolphin-unleashed"; then
  echo "Pulling base model..."
  ollama pull dolphin-llama3:8b
  echo "Creating dolphin-unleashed..."
  ollama create dolphin-unleashed -f /Modelfile
  echo "Model ready."
else
  echo "dolphin-unleashed already exists."
fi
