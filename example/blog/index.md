---
name: 'index'
processor: 'blog'
title: 'News'
pagination: 2
---

{{if .prevPage}} <a href="{{.prevPage}}">zurück</a>{{end}} {{if .nextPage}} <a href="{{.nextPage}}">nächste Seite</a>{{end}} ({{.entryCount}} Einträge)

{{.blogentries}}

<hr/>
{{.actualPage}}/{{.pageCount}}