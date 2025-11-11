#!/usr/bin/env python3
"""
Helper script to convert real mushroom images to base64 for k6 load testing.

Usage:
  1. Place mushroom images in a folder (e.g., test-images/)
  2. Run: python3 generate-test-images.py test-images/
  3. Copy the output array into load-test.js to replace sampleImages
"""

import base64
import sys
import os
from pathlib import Path

def image_to_base64(image_path):
    """Convert image file to base64 string."""
    with open(image_path, 'rb') as img_file:
        return base64.b64encode(img_file.read()).decode('utf-8')

def main():
    if len(sys.argv) < 2:
        print("Usage: python3 generate-test-images.py <image-folder>")
        print("Example: python3 generate-test-images.py test-images/")
        sys.exit(1)

    folder = sys.argv[1]

    if not os.path.exists(folder):
        print(f"Error: Folder '{folder}' does not exist")
        sys.exit(1)

    # Supported image formats
    extensions = ['.jpg', '.jpeg', '.png', '.gif', '.webp']

    # Find all image files
    image_files = []
    for ext in extensions:
        image_files.extend(Path(folder).glob(f'*{ext}'))
        image_files.extend(Path(folder).glob(f'*{ext.upper()}'))

    if not image_files:
        print(f"No image files found in '{folder}'")
        print(f"Supported formats: {', '.join(extensions)}")
        sys.exit(1)

    print(f"Found {len(image_files)} images. Converting to base64...\n")

    # Generate JavaScript array
    print("const sampleImages = [")
    for i, img_path in enumerate(image_files):
        try:
            base64_str = image_to_base64(img_path)
            # Truncate display for readability
            display_str = base64_str[:60] + '...' if len(base64_str) > 60 else base64_str
            print(f"  '{base64_str}',  // {img_path.name} ({len(base64_str)} chars)")
        except Exception as e:
            print(f"  // Error reading {img_path.name}: {e}", file=sys.stderr)
    print("];")

    print(f"\n✓ Generated {len(image_files)} base64 encoded images")
    print("Copy the array above into load-test.js to replace sampleImages")

if __name__ == "__main__":
    main()
