---
name: 'index'
processor: gallery
title: 'index'
---
This is a new page with the title {{.title}} {{ if .section }} in section {{.section.title}} {{ end }}