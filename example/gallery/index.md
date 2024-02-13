---
name: 'index'
processor: 'gallery'
title: 'Gallerie'
images: 'images'
thumbswidth: 200
crop: true
fluid: false
imageproperties: [description, tags] 
imagecontainer: '<div style="display: flex;overflow: hidden;flex-wrap: wrap;justify-content: center;">{{`{{.images}}`}}</div>'
imageentry: '<div style="display: inline-block;overflow: hidden;width:{{`{{.thumbswidth}}`}}px;padding: 5px 5px 5px 5px;"><a href="{{`{{.source}}`}}" target="_blank"><img loading="lazy" src="{{`{{.thumbnail}}`}}" alt="{{`{{.name}}`}}"><span>{{`{{.name}}`}}<br/>Beschreibung: {{`{{.description}}`}}<br/>Größe: {{`{{.size}}`}}</span></a></div><br/>'
style: ''
---
This is a gallery page with the title {{.title}}

{{.images}}