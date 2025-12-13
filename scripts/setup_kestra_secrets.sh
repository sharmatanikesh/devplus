#!/bin/bash

# Define the secrets file
INPUT_ENV=".kestra_secrets"
OUTPUT_ENV=".env_encoded"

# Check if input file exists
if [ ! -f "$INPUT_ENV" ]; then
    echo "Creating template $INPUT_ENV..."
    echo "GEMINI_API_KEY=your_gemini_api_key_here" > "$INPUT_ENV"
    echo "Please edit $INPUT_ENV with your actual secrets."
    exit 1
fi

echo "Encoding secrets from $INPUT_ENV to $OUTPUT_ENV..."

# Clear output file
> "$OUTPUT_ENV"

# Read and encode each line
while IFS='=' read -r key value || [ -n "$key" ]; do
    # Skip empty lines or comments
    if [[ -z "$key" || "$key" == \#* ]]; then continue; fi
    
    # Base64 encode the value
    # MacOS 'base64' default is fine, Linux might need -w 0
    if [[ "$OSTYPE" == "darwin"* ]]; then
        ENCODED_VALUE=$(echo -n "$value" | base64)
    else
        ENCODED_VALUE=$(echo -n "$value" | base64 -w 0)
    fi
    
    # Write to output with SECRET_ prefix
    echo "SECRET_$key=$ENCODED_VALUE" >> "$OUTPUT_ENV"
done < "$INPUT_ENV"

echo "Done! Secrets encoded in $OUTPUT_ENV."
echo "You can now run: docker-compose up -d"
