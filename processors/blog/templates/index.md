---
name: "{{.pagename}}"
title: "{{.pagename}}"
processor: "{{.processor}}"
pagination: 3
---
This is a new blog with the title {{.title}}

{{if .prevPage}} <a href="{{.prevPage}}">back</a>{{end}} {{if .nextPage}} <a href="{{.nextPage}}">next</a>{{end}}  ({{.entryCount}} Entries)

{{.blogentries}}
<hr/>
{{.actualPage}}/{{.pageCount}}