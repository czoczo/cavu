#!/usr/bin/python3

# Separate script to generate static entries for CasaVue config based on "Awesome Selfhosted" content on GitHub

import re
import requests

sourceUrl = "https://raw.githubusercontent.com/awesome-selfhosted/awesome-selfhosted/master/README.md"

# Retrieve makdown and save it to list of lines
response = requests.get(sourceUrl)
lines = response.text.split('\n')

# Pattern for catching groups/namespaces
pattern_namespace = r"^### (.+)$"

# Pattern for catching appications with DEMO links
pattern_item = r"^- \[([^\]]+)\]\([^\)]+\) - (.+) \(\[Demo\]\(([^\)]+)\)"

for line in lines:
    match_namespace = re.search(pattern_namespace, line)
    match_item = re.search(pattern_item, line)

    if match_namespace:
        namespace = match_namespace.group(1)
        continue

    if not match_item:
        continue

    name = match_item.group(1)
    description = match_item.group(2).replace("\"","\\\"")
    url = match_item.group(3)
    print(
      f"""    - name: "{name}"
      namespace: "{namespace}"
      url: "{url}"
      description: "{description}\""""
    )
