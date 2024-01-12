---
name: 'index'
processor: 'gallery'
title: 'Gallerie'
images: 'images'
thumbswidth: 200
imgproperties: 
  - description
  - tags
imageentry: '<div style="display: inline-block;overflow: hidden;width:200px;height:280px;padding: 5px 5px 5px 5px;"><a href="{{`{{.source}}`}}"><img loading="lazy" src="{{`{{.thumbnail}}`}}" alt="{{`{{.name}}`}}"><p style="margin-top: 8px;">{{`{{.name}}`}}<br/>Beschreibung: {{`{{.description}}`}}<br/>Größe: {{`{{.size}}`}}</p></a></div>'
---
This is a new gallery page with the title {{.title}}

{{.images}}