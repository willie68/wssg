---
name: 'index'
processor: 'gallery'
title: 'Gallerie'
images: 'images'
thumbswidth: 400
crop: false
fluid: true
imageproperties: [description, tags]
imagecontainer: '<div style="display: flex;overflow: hidden;flex-wrap: wrap;justify-content: center;">{{`{{.images}}`}}</div>'
imageentry: 
style: 
imagelist: [kumpan-electric, balazs-ketyi, ryan-ancill, theme-photos, brad-neathery, budka-damdinsuren, faizur-rehman, glenn-carstens-peters, kelly-sikkema ]
listonly: true
---
This is a new fluid gallery page with the title "{{.title}}"

{{.images}}