---
name: "{{.pagename}}"
title: "{{.pagename}}"
processor: "{{.processor}}"
---
**{{.title}}**

This is a new blog entry with the title {{.title}} 

created at {{ dtFormat .created "Monday, 2.01.06" "en_US" }}

