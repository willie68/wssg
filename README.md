# WSSG

Willie's Static Site Generator

# Warum?

Darum. Ich finde jeder Programmierer sollte so ein paar Sachen mal programmiert haben. Hello World, Sortieralgorithmus, Datenbank, Website Generator. Letzterer fehlte mir noch.

# Basis

wssg ist ein kleiner feiner Websitegenerator mit folgenden Features (wenn er denn fertig ist, und was ist schon fertig...)

- Basis sind Markdown-Dateien
- automatische Indexerstellung
- eingebauter WebServer zum besseren Testen der Site
- austauschbare Template Engine

# Installation

wssg ist ein typisches Copy/Run Programm, deswegen auch in Golang geschrieben. Es kann aber auch mithilfe von Golang installiert werden. 

## Copy/Run

Nimm einfach die richtige Release Version und kopiere sie in ein Verzeichnis, dass bereits in deinem Pfad vorhanden ist. Oder kopiere das Binary in das root der neuen Site. 

Zum Testen ob alles funktioniert gibt einfach `wssg version` ein.

## via Golang

`go install github.com/willie68/wssg`

# Quickstart

Für den schnellen Start mit dem `wssg` wird zunächst eine Installation vorausgesetzt.

Als erstes erzeugen wir uns eine neue Site. 

`wssg init ./<sitename>`

Jetzt wird automatisch ein Verzeichnis mit dem Namen <sitename> erzeugt und dort dann alle wichtigen Konfiguration erzeugt. Auch wird direkt die erste Seite (index.md) mit erzeugt.

`wssg generate`  

generiert nun die Website. Das Ergebnis landet automatisch im Ordner `.wssg/output`

Wenn du den internen Webserver für den schnelleren Output verwenden möchtest, starte diesen mit folgendem Befehl:

`wssg serve`

 Jetzt wird zunächst die Webseite neu generiert und dann ein Webserver gestartet. Unter

http://localhost:8080/ 

kannst du dir das Ergebnis anschauen. Während der Webserver läuft, kannst du nun deine Webseite bearbeiten. Jede Änderung wird automatisch vom `wssg` registriert und die Seiten entsprechend upgedated. Ein einfaches F5 im Browser reicht, um dir das Ergebnis deiner Änderungen anzuschauen. Änderungen im Ordner `.wssg` werden **nicht** automatisch berücksichtigt. Dazu muss auch eine Änderung an einer anderen Datei (außerhalb des `.wssg` Ordners) erfolgen. 

Der Inhalt der `generate.yaml` wird nur beim Start des Server ausgewertet.

# Aufbau

Das Programm ist für folgende Struktur am besten geeignet. Die erste Ebene (Root) ist quais der Startpunkt. Dort muss für den Start eine index.md erstellt werden. Diese wird automatisch beim `wssg init` angelegt. Hier können dann weitere Seiten (pages) hinzugefügt werden. Zusätzliche Dateien, wie z.B. Bilder, Stylesheets, JS usw. können sowohl in den Rootordner wie auch in weiteren Unterordnern abgelegt werden. Unterordner können dann einfach per relativer Angabe referenziert werden. 

Möchte man einen neuen Bereich (section) mit verschiedene Seiten anlegen, kann man das mit `wssg new section <name>` machen. Dabei wird nun, ebenso wie im root Ordner, ein Unterverzeichnis .wssg mit den Einstellungen für diesen Bereich erstellt.  

# Programmparameter

# Seitenaufbau

## Frontmatter für Markdown

Die Markdown-Dateien sollten den Inhalt sollten mit einem Frontmatter Bereich starten. Dieser startet am Anfag der Datei mit `---` und endet ebenfalls mit `---`. Dazwischen steht ein Bereich mit Optionen für die aktuelle Seite im yaml Format:

```yaml
---
name: 'index'
processor: markdown
title: 'Index'
order: 10
---
```

`name`: technischer Name der Seite. Dieser wird u.a. für die Referenzierung und für die Html-Generierung verwendet. Erlaubt sind folgende Zeichen: `a-z,0-9,-,_` 

`processor`: Der Prozessor steht für den zu verwendenden Generierungsprozessor. Derzeit steht nur `markdown`, `gallery` und `plain` zur Verfügung. 

`title`: Der Seitentitel. Hier können auch Sonderzeichen verwendet werden.

`order`: steht für die Sortierungsfolge. Beim Abruf aller Seiten über {{ range .pages}} werden die Seiten nach dieser Reihenfolge aufsteigend sortiert. Der absolute Wert spielt keine Rolle, d.h. es muss nicht 0,1,2 verwendet werden. Um nachträglich Seiten einzufügen kann man auch 10 , 20, 30 für den Start benutzen. So kann man später neue Seite bei 15, 25 usw. einfügen.

Es können weitere Parameter angegeben werden, die von den jeweiligen Plugin/Prozessor definiert werden.  Oder auch nur von der eigenen Seite.

## Variablen für eine Seite

`{{.body}}` ergibt den konvertierten Text aus der Markdown Datei.

`{{.site.#}}` sind die Einstellungen für die gesamte Website. Hier stehen 1:1 alle Einstellungen aus der `siteconfig.yaml`. Beispielsweise  

`{{.site.Language}}` ergibt die Sprache

`{{.site.Title}}` den Webseitentitel. Ebenso funktionieren `{{.site.Description}}` und `{{.site.Keywords}}`

Unter `{{.site.Userproperties}}` stehen alle unbekannten Parameter zur Verfügung. Diese können von dem HTML Template definiert werden. Als Beispiel dient der `font` Parameter. Will man also den in der Seitenkonfiguration angegeben Font verwenden, gelingt das mit `{{.site.UserProperties.font}}`. Diese Userproperties werden auch als Defaults für Bereiche- und Seitenkonfigurationen verwendet. Weitere bereits definierte Userproperties: socialmedia oder webcontact

Für die aktuelle Seite sind folgende Variablen definiert:

`{{.page.URLPath}}` der relative Pfad der Seite

`{{.page.Name}}` der Name der Seite

`{{.page.Path}}` der Pfad auf die Ursprungsdatei 

`{{.page.Title}}` der Titel der Seite

 `{{range .pages}} ... {{end}}` kann über alle Seiten eines Bereiches iteriert werden. Innerhalb sind dann folgende Punkt Variablen definiert und verweisen auf die jeweilige Seite:

`{{.URLPath}}` der relative Pfad der Seite

`{{.Name}}` der Name der Seite

`{{.Path}}` der Pfad auf die Ursprungsdatei 

`{{.Title}}` der Titel der Seite

# Plugins

## markdown

Markdown ist ein Plugin oder besser Prozessor, der MD Dateien in HTML verwandelt. Dabei werden automatisch die o.g. Ersetzungen berücksichtigt. 

## plain

Beim plain Plugin wird der Seiteninhalt ohne Prozessor direkt als HTML interpretiert. Ersetzungen werden jedoch vorgenommen, die Seite aber nicht weiter verarbeitet. Dieses Plugin ist als Default gesetzt.

## gallery

Wird ein Prozessor gallery gesetzt, wird eine Bildgallery generiert. Folgende Frontmatter Parameter werden zusätzlich verwendet:

```yaml
---
name: 'index'
processor: 'gallery'
title: 'index'
images: 'images'
thumbswidth: 200
crop: true
imgproperties: 
  - description
  - tags
imagecontainer: '{{`{{.images}}`}}'
imageentry: '<div style="display: inline-block;overflow: hidden;width:200px;height:280px;padding: 5px 5px 5px 5px;"><a href="{{`{{.source}}`}}"><img src="{{`{{.thumbnail}}`}}" alt="{{`{{.name}}`}}"><p style="margin-top: 8px;">{{`{{.name}}`}}<br/>Beschreibung: {{`{{.description}}`}}<br/>Größe: {{`{{.size}}`}}</p></a></div>'
---
```

`images`: gibt das Verzeichnis an, wo die zu verarbeitenden Bilddaten liegen. Es kann nur ein Ordner angegeben werden. Alle Bilddaten darin werden dann verarbeitet. Als Bilder werden Dateien mit folgenden Endungen betrachtet: `*.jpeg, *.jpg, *.bmp, *.png` 

`thumbswidth`: ist die Breite der Thumbs, die von dem Plugin automatisch generiert werden.

`crop`: mit der boolschen Ausdruck crop kann man die Thumbnails entsprechend ihrer Breite abschneiden. Bei `false` bleibt bei den Thumbs das Seitenverhältnis erhalten, `true` erzeugt quadratische Thumbnails der Breite `thumbswidth`. 

`imgproperties`: Hier kann man optional eine Liste zusätzlicher Bildeigenschaften hinterlegen. Bei der Generierung wird dann im Bildordner eine Datei `_content.yaml` angelegt. Diese enthält pro Bild dann die entsprechenden Eigenschaften.

```yaml
balazs-ketyi:
    description: Monitor
    tags: Monitor, Arbeitsplatz
brad-neathery:
    description: Mann übt Schrift
    tags: Schreiben, Kamera, Laptop
budka-damdinsuren:
    description: Apple Monitor
    tags: Apple, Monitor
faizur-rehman:
    description: Arbeitsplatz
    tags: Arbeitsplatz
...
```

Wird ein neues Bild in den Ordner gelegt, wird bei der nächsten Generierung auch hier ein neuer Eintrag erzeugt. Vorhandene Einträge bleiben immer erhalten. Einträge für gelöschte Bilder werden derzeit nicht entfernt, aber auch nicht weiter berücksichtigt. Diese Eigenschaften können dann im `imageentry` entsprechend benutzt werden. 

`imagecontainer`: Hier kann man seinen eigenen Container um die Bilderliste definieren. An die Position des Macros `.images` kommt dann die Bilderliste.

`imageentry`: Hier steht das HTML Template, welches ein einzelnes Bild darstellt. Bitte beachte: für jedes Bild wird dieses Template einmal generiert. Dabei werden auch wieder die typischen Ersetzungen gemacht. {{.title}} würde also auch hier den Seitentitel einfügen. Um auf die Eigenschaften der Datei zugreifen zu können, müssen diese zus. gekennzeichnet werden. Dazu dient der folgende Ausdruck:

```
{{`{{.<keyname>}}`}}
```

keyname kann folgende Eigenschaften verwenden

`.source`: ergibt den relativen Pfad/Namen der Quelldatei (inkl. Unterordner)
`.thumbnail`: ergibt den relativen Pfad/Namen der Thumbnaildatei (inkl. Unterordner)
`.name`: ergibt den Dateinamen ohne Pfad und Endung
`.size`: ist die vom Menschen lesbare Dateigröße: 10KiB oder 12MiB...
`.sizebytes`: ist die Dateigröße in Bytes

Die aufbereitete Bilderliste wird dann an die Stelle `{{.images}}` der MD Datei eingefügt.

# Beispiel

Ein Beispiel für die Vielseitigkeit des `wssg` befindet sich im Verzeichnis `example`. Dieses kann man `wssg generate` oder `wssg serve` verwenden.
