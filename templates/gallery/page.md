---
name: "{{.pagename}}"
title: "{{.pagename}}"
processor: "{{.processor}}"
images: 'images'
thumbswidth: 200
imageentry: '<div style="display: inline-block;overflow: hidden;width:200px;height:250px;padding: 5px 5px 5px 5px;"><a href="{{"{{`{{.source}}`}}"}}"><img src="{{"{{`{{.thumbnail}}`}}"}}" alt="{{"{{`{{.name}}`}}"}}"><p style="margin-top: 8px;">{{"{{`{{.name}}`}}"}}<br/>size: {{"{{`{{.size}}`}}"}}</p></a></div>'

---
This is a new gallery page with the title {{.title}}

{{.images}}