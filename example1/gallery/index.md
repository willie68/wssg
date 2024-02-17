---
crop: true
fluid: false
imagecontainer: '<div style="display: flex;overflow: hidden;flex-wrap: wrap;justify-content: center;">{{.images}}</div>'
imageentry: '<div style="display: inline-block;overflow: hidden;width:200px;height:280px;padding: 5px 5px 5px 5px;"><a href="{{`{{.source}}`}}"><img loading="lazy" src="{{`{{.thumbnail}}`}}" alt="{{`{{.name}}`}}"><p style="margin-top: 8px;">{{`{{.name}}`}}<br/>Titel: {{`{{.title}}`}}<br/>Größe: {{`{{.size}}`}}</p></a></div><br/>'
imageproperties:
    - title
images: images
name: 'index'
processor: 'gallery'
thumbswidth: 200
title: 'index'
---
This is a new gallery page with the title {{.title}}

{{.images}}