---
name: 'index'
processor: 'gallery'
title: 'Gallerie'
images: 'images'
thumbswidth: 500
crop: false
fluid: true
imgproperties: 
  - description
  - tags
imageentry: '<div style="display: inline-block;overflow: hidden;width:{{`{{.thumbswidth}}`}}px;padding: 5px 5px 5px 5px;"><a href="{{`{{.source}}`}}" target="_blank"><img loading="lazy" src="{{`{{.thumbnail}}`}}" alt="{{`{{.name}}`}}"><span>{{`{{.name}}`}}<br/>Beschreibung: {{`{{.description}}`}}<br/>Größe: {{`{{.size}}`}}</span></a></div><br/>'

---
This is a new fluid gallery page with the title {{.title}}

{{.images}}