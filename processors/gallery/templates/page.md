---
name: "{{.pagename}}"
title: "{{.pagename}}"
processor: "{{.processor}}"
images: 'images'
thumbswidth: 200
crop: true
fluid: false
imageproperties: [title]
imagecontainer: '<div style="display: flex;overflow: hidden;flex-wrap: wrap;justify-content: center;">{{`{{.images}}`}}</div>'
imageentry: '<div style="display: inline-block;overflow: hidden;width:200px;height:280px;padding: 5px 5px 5px 5px;"><a href="{{"{{`{{.source}}`}}"}}"><img loading="lazy" src="{{"{{`{{.thumbnail}}`}}"}}" alt="{{"{{`{{.name}}`}}"}}"><p style="margin-top: 8px;">{{"{{`{{.name}}`}}"}}<br/>Titel: {{"{{`{{.title}}`}}"}}<br/>Größe: {{"{{`{{.size}}`}}"}}</p></a></div><br/>'

---
This is a new gallery page with the title {{.title}}

{{.images}}