---
name: 'index'
processor: 'blog'
title: 'index'
pagination: 3
---
This is a new page with the title {{.title}}

{{if .prevPage}} <a href="{{.prevPage}}">zurück</a>{{end}} {{if .nextPage}} <a href="{{.nextPage}}">nächste</a>{{end}} 

{{.blogentries}}