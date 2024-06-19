---
title: "{{ replaceRE `^\d\d-\d\d__` "" .File.ContentBaseName | humanize | title }}"
date: {{ .Date }}
image: images/catan-close-up.png
categories:
  - board games
draft: false
---
