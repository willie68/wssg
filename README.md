# WSSG

Willie's Static Site Generator

# Warum?

Darum. Ich finde jeder Programmierer sollte so ein paar Sachen mal programmiert haben. Hello World, Sortieralgorithmus, Datenbank, Website Generator. Letzterer fehlte mir noch.

Und da ich gerade eine einfache Seite für einen Freund bauen sollte, die er später auch selber pflegen kann, kam ich auf diese Idee. Natürlich habe ich erst andere Generatoren ausprobiert. Ganz viele waren vor allem für ein Thema Blogs, Blogs und nochmal Blogs. Aber ich wollte keinen Blog, ich wollte einfach ein paar einfache Seiten am besten als markdown mit ein paar Bildern, mehr nicht. Ich wollte keine neue Verzeichnisstruktur lernen, einfach ein Verzeichnis erzeugen, Dateien reinlegen, oder noch besser nur editieren, weil die Dateien schon da sind. Evtl. noch das Layout anpassen und schon hat man seine statische Seite.  

# Basis

wssg ist ein kleiner feiner Websitegenerator mit folgenden Features (wenn er denn fertig ist, und was ist schon fertig...)

- Basis sind Markdown-Dateien, evtl. auch HTML
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

Für den schnellen Start mit dem `wssg` wird zunächst eine [Installation](#Installation) vorausgesetzt.

Als erstes erzeugen wir uns eine neue Site. 

`wssg init ./<sitename>`

Jetzt wird automatisch ein Verzeichnis mit dem Namen <sitename> erzeugt und dort dann alle wichtigen Konfiguration erzeugt. Auch die erste Seite (index.md) wird direkt mit erzeugt.

Wechsele in das neu erzeugt Verzeichnis.

`cd <sitename>`

`wssg generate`  generiert nun die Website. Das Ergebnis landet automatisch im Ordner `.wssg/output`

MIt `wssg serve` startest du den internen Webserver. Dieser generiert zunächst deine Seite und startet dann automatisch deinen Browser. Jetzt kannst du das Ergebnis direkt im Browser betrachten. Oder du rufst selber die generierte Webseite auf 

http://localhost:8080/ 

Während der Webserver läuft, kannst du nun deine Webseite bearbeiten. Jede Änderung wird automatisch vom `wssg` registriert und die Seiten entsprechend upgedated. Der Browser aktualisiert sich im Sekundentakt selber. Im Ordner `.wssg` werden Änderungen an `siteconfig.yaml` und `generate.yaml` **nicht** automatisch berücksichtigt. Änderungen an der `layout.html` werden jedoch berücksichtigt. Das Löschen des `output` Ordners, wo die generierten Daten abgelegt werden, triggert eine komplette neue Generierung.

# Aufbau

Das Programm ist für folgende Struktur am besten geeignet. Die erste Ebene (Root) ist quasi der Startpunkt. Dort muss für den Start eine index.md erstellt werden. Diese wird automatisch beim `wssg init` angelegt. Hier können dann weitere Seiten (pages) hinzugefügt werden. Zusätzliche Dateien, wie z.B. Bilder, Stylesheets, JS usw. können sowohl in den Rootordner wie auch in weiteren Unterordnern abgelegt werden. Unterordner können dann einfach per relativer Angabe referenziert werden. 

Möchte man einen neuen Bereich (section) mit verschiedene Seiten anlegen, kann man das mit `wssg new section <name>` machen. Dabei wird nun, ebenso wie im root Ordner, ein Unterverzeichnis .wssg mit den Einstellungen für diesen Bereich erstellt.  

# Programmparameter

# Site Einstellungen

In der Datei .wssg/siteconfig.yaml werden die Einstellungen der gesamten Website verwaltet. 

```yaml
baseurl: design_sauber.com
title: Design Sauber
description: Design Sauber, Advertising for everything and everyone
keywords: tutorial basic static website
language: de
font: Tahoma, Verdana, sans-serif
webcontact:
 url: mailto:info@example.com
 title: info@example.com
socialmedia:
 facebook:
  title: FB
  icon: /images/social_fb.png
  url: https://www.facebook.com/wilfried.klaas/
 youtube:
  title: YT
  icon: /images/social_yt.png
  url: https://www.youtube.com/channel/UCg5ZpZJGuLgz4maETfUc9EA
cookiebanner:
 enabled: true
 text: ''
```

Die Eigenschaften sind eigentlich selbsterklärend. Alle Eigenschaften stehen auf jeder Seite zur Verfügung und können auch von jeder Seite überschrieben werden. Zusätzlich sind diese auch unter dem Bereich "site" (nicht überschreibbar) zugreifbar. 

### Besonderheiten

`cookiebanner:` Mit dem cookiebanner kann eine Cookiebanner aktiviert werden. Der angegebene Text wird dann automatisch beim 1. Aufruf der Startseite eingeblendet.

# Seitenaufbau

## Frontmatter für Markdown

Die Markdown-Dateien sollten vor dem Inhalt mit einem Frontmatter Bereich starten. Dieser startet am Anfang der Datei mit `---` und endet ebenfalls mit `---`. Dazwischen steht ein Bereich mit Optionen für die aktuelle Seite im yaml Format:

```yaml
---
name: 'index'
processor: markdown
title: 'Index'
order: 10
---
```

`name`: technischer Name der Seite. Dieser wird u.a. für die Referenzierung und für die Html-Generierung verwendet. Erlaubt sind folgende Zeichen: `a-z,0-9,-,_` Die Namen in einem Bereich sollten alle unterschiedlich sein. Am einfachsten kann man das machen, indem man als Namen einfach den Dateinamen ohne .md verwendet.

`processor`: Der Prozessor steht für den zu verwendenden Generierungsprozessor.  

`title`: Der Seitentitel für die Anzeige z.B. in Menüs. Hier können auch Sonderzeichen verwendet werden.

`order`: steht für die Sortierungsfolge. Beim Abruf aller Seiten über {{ range .pages}} werden die Seiten nach dieser Reihenfolge aufsteigend sortiert. Der absolute Wert spielt keine Rolle, d.h. es muss nicht 0,1,2 verwendet werden. Um nachträglich Seiten einzufügen kann man auch 10 , 20, 30 für den Start benutzen. So kann man später neue Seite bei 15, 25 usw. einfügen.

Es können weitere Parameter angegeben werden, die von den jeweiligen layout/Plugin/Prozessor definiert werden.  Oder auch nur von der eigenen Seite als Makros verwendet werden.

## Variablen/Makros für eine Seite

`{{.body}}` ergibt den konvertierten Text aus der Markdown Datei.

`{{.site.#}}` sind die Einstellungen für die gesamte Website. Hier stehen 1:1 alle Einstellungen aus der `siteconfig.yaml`. Beispielsweise  

`{{.site.Language}}` ergibt z.B. die Sprache oder 

`{{.site.Title}}` den Webseitentitel. Ebenso funktionieren `{{.site.Description}}` und `{{.site.Keywords}}`

Unter `{{.site.Userproperties}}` stehen alle unbekannten Parameter zur Verfügung. Diese können von dem HTML Template verwendet werden. Als Beispiel dient der `font` Parameter. Will man also den in der Seitenkonfiguration angegeben Font verwenden, gelingt das mit `{{.site.UserProperties.font}}`. Diese Userproperties werden auch als Defaults für Bereiche- und Seitenkonfigurationen verwendet. Weitere bereits definierte Userproperties: `socialmedia` oder `webcontact`

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

`{{range $id, $el := .sections}} ... {{end}}` kann über alle Bereiche iteriert werden. Auch hier sind die gleichen Makros wie bei den Seiten möglich.

# Prozesoren

## Markdown

Markdown ist ein Prozessor, der MD Dateien in HTML verwandelt. Dabei werden automatisch die o.g. Makros berücksichtigt. 

## Plain

Beim plain Prozessor wird der Seiteninhalt ohne Prozessor direkt als HTML interpretiert. Makros werden jedoch ausgewertet, die Seite aber nicht weiter verarbeitet. Dieser Prozessor ist als Default gesetzt.

## Gallery

Wird ein Prozessor `gallery` gesetzt, wird eine Bildgallery generiert. Folgende Frontmatter Parameter werden zusätzlich verwendet:

```yaml
---
name: 'index'
processor: 'gallery'
title: 'index'
images: 'images'
thumbswidth: 200
crop: true
fluid: false
imageproperties: [description, tags]
imagecontainer: '{{`{{.images}}`}}'
imageentry: '<div style="display: inline-block;overflow: hidden;width:200px;height:280px;padding: 5px 5px 5px 5px;"><a href="{{`{{.source}}`}}"><img src="{{`{{.thumbnail}}`}}" alt="{{`{{.name}}`}}"><p style="margin-top: 8px;">{{`{{.name}}`}}<br/>Beschreibung: {{`{{.description}}`}}<br/>Größe: {{`{{.size}}`}}</p></a></div><br/>'
style: 'ownstyle'
---
imagelist: [kumpan-electric, balazs-ketyi, ryan-ancill, theme-photos, brad-neathery, budka-damdinsuren, faizur-rehman, glenn-carstens-peters, kelly-sikkema ]
listonly: true
---
```

`images`: gibt das Verzeichnis an, wo die zu verarbeitenden Bilddaten liegen. Es kann nur ein Ordner angegeben werden. Alle Bilddaten darin werden dann verarbeitet. Als Bilder werden Dateien mit folgenden Endungen betrachtet: `*.jpeg, *.jpg, *.bmp, *.png` Der Bildordner kann irgendwo auf deiner Festplatte liegen. Die verwendeten Bilder werden dann automatisch kopiert.

`thumbswidth`: ist die Breite der Thumbs, die von dem Prozessor automatisch generiert werden.

`crop`: mit der boolschen Ausdruck crop kann man die Thumbnails entsprechend ihrer Breite abschneiden. Bei `false` bleibt bei den Thumbs das Seitenverhältnis erhalten, `true` erzeugt quadratische Thumbnails der Breite `thumbswidth`. 

`fluid`: Mit fluid wird eine fluide Gallery erzeugt. Diese hat 3 Spalten. Zusätzlich wird bei kleineren Displays auf eine einspaltige Anzeige umgeschaltet. Ein optimaler Imageentry für die fluide Gallery wäre folgender:

```yaml
imageentry: '<div style="display: inline-block;overflow: hidden;width:{{`{{.thumbswidth}}`}}px;padding: 5px 5px 5px 5px;"><a href="{{`{{.source}}`}}" target="_blank"><img loading="lazy" src="{{`{{.thumbnail}}`}}" alt="{{`{{.name}}`}}"><span>{{{{.name}}}}<br/>Beschreibung: {{{{.description}}}}<br/>Größe: {{{{.size}}}}</span></a></div><br/>'
```

`imageproperties`: Hier kann man optional eine Liste zusätzlicher Bildeigenschaften hinterlegen. Bei der Generierung wird dann im Seitenordner eine Datei `_<seitenname>.props` angelegt. Diese enthält pro Bild dann die entsprechenden Eigenschaften. Die Datei wird automatisch generiert, wenn diese noch nicht vorhanden ist. Werden neue Bilder dem Ordner hinzugefügt, müssen diese per Hand in die Datei eingefügt werden. 

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

Die Eigenschaften können dann im `imageentry` als Makros benutzt werden. 

`imagecontainer`: (Optional) Hier kann man seinen eigenen Container um die Bilderliste definieren. An die Position des Macros `.images` kommt dann die Bilderliste.

`imageentry`: (Optional) Hier steht das HTML Template, welches ein einzelnes Bild darstellt. Ist das Feld leer oder nicht vorhanden wird der Default benutzt. D.h. für jedes Bild wird zusätzlich zum Thumbnail auch der Name und die Liste von Properties generiert. Anschauen kann man sich den Default im Beispiel auf der Fluidgallerie.

Bitte beachte: für jedes Bild wird dieses Template einmal generiert. Dabei werden auch wieder die typischen Ersetzungen gemacht. `{{.title}}` würde also auch hier den Seitentitel einfügen. Um auf die Eigenschaften der Datei zugreifen zu können, müssen diese zus. gekennzeichnet werden. Dazu dient der folgende Ausdruck:

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

`style`: Hier kann man seine eigenen CSS Style definieren. Dieser wird anstatt des Defaults geladen. 

`imagelist:` enthält eine Liste der Dateinamen (ohne Endung). In dieser Reihenfolge werden die Bilder dann angeordnet. Zusätzliche Bilder werden dahinter hinzugefügt, falls `listonly` nicht vorhanden oder auf false gesetzt wurde.

`listonly:` es werden nur die Images, die in `imagelist` stehen, verwendet. Weitere Images in dem Quellordner werden für diese Galerie ignoriert.

## Blog

Um einen Blog zu erzeugen, wird der Blog Prozessor verwendet. Zunächst muss ein Bereich mit dem Blog Prozessor erzeugt werden.

`wssg new section -p blog MyBlog`

erzeugt einen neuen Bereich `MyBlog`. Die neue `index.md` in dem Verzeichnis MyBlog ist dann der Container für die verschiedenen Einträge. Folgende Einstellungen findet man in der Datei (Frontmatter yaml)

```yaml
name: 'index'
processor: 'blog'
title: 'index'
pagination: 3
order: 'asc' #oder 'desc'
entrytemplate: '.content'
```

`name, processor` und `title` kennen wir zwar schon, aber hier werden die Einträge nur für die interne verarbeitung genutzt. D.h. konkret bei diesem Blogprozessor muss die Einstiegsseite `index` heißen. Also den Namen hier nicht ändern. Den Namen der Sektion kann man in der `.wssg/section.yaml` ändern. `pagination` gibt an wie viele Blogeinträge es pro Seite gibt. Jeder Blogeintrag ist eine eigene Markdown-Datei. Die Reihgenfolge ist automatisch nach Erzeugungsdatum absteigend sortiert, wenn `order` nicht angegeben wird. Bei `order: 'asc'` wird automatisch nach Erzeugungsdatum aufsteigend sortiert, also das älteste zuerst. mit entrytemplate kann man zusätzlich das Format jedes einzelnen Eintrages modifizieren. 

Beispielsweise würde:

```yaml
entrytemplate: '<div><table><tr><td>{{`{{if .entryeven}}`}}<img src="{{`{{.image}}`}}" height="300"/></td><td>{{`{{.content}}`}}{{`{{else}}`}}{{`{{.content}}`}}</td><td><img src="{{`{{.image}}`}}" height="300"/>{{`{{end}}`}}</td></tr></table></div>'
```

einen Wechselblog erzeugen. D.h. es gibt pro Blog 2 Einträge, `.content` ist der tatsächliche Text und `.image` wurde einfach im FM des Blobeintrages hinzugefügt und enthält einen Link zu einem Bild.
`.entryeven` wird vom Blogprozessor erzeugt und ist `true`, wenn eine gerade Seite gerendert werden soll. 
`.entrynumber` ist die laufende Nummer des Blogs und kann ebenfalls verwendet werden.

Nun wird bei geraden Einträgen links das Bild und rechts der Text erzeugt. Bei ungeraden Einträgen genau anders herum.

Und nochmal, wundert euch nicht über die doppelte `{{ }}` Syntax. Dieses Template wird 2x verarbeitet. Einmal wenn die Seite startet, um z.B. Titel, Description oder andere Dinge aus der zentralen Konfiguration oder der Seitenkonfiguration zu benutzen. Pro Blogeintrag wird das Template dann nochmal verarbeitet, um die eintragspezifischen Dinge einzusetzen.  

Einen neuen Eintrag erzeugt man mit 

`wssg new page -p blog MyBlog/Blog1`

Der neue Blog Eintrag hat nun den Namen Blog1 und kann mit der Datei blog1.md editiert werden. 

Auch hier gibt es etwas Frontmatter yaml:

```yaml
created: 2024-02-18T11:35:12.5003594+01:00
name: 'blog8'
processor: 'blog'
title: 'Neuer Kunde'
```

`created` gibt nur an, wann die Seite erzeugt worden ist. Änderungen wirken sich nicht aus. Zusätzlich wird eine List von Blogeinträgen in der Datei `_content.yaml` gepflegt. Diese enthält den Namen des Eintrags und das Erzeugungsdatum, welches dann auch im Blogeintrag als Makro zur Verfügung steht. 

Neue Einträge kann man hier belibig machen und z.B. im `entrytemplate` verwenden.

**Hinweis**: will man den Namen des Bereiches ändern, reicht es nicht, den Namen (title) in der index.md zu ändern. Der Name des Bereiches (title) steht in der `.wssg/section.yaml`. 

# Beispiel

Ein Beispiel für die Vielseitigkeit des `wssg` befindet sich im Verzeichnis `example`. Dieses kann man `wssg generate` oder `wssg serve` verwenden.

# FAQ

## Kann ich die verschiedenen Bereich in dem Hauptmenü nach eigenen Kriterien sortieren?

Üblicherweise werden die Bereiche automatisch nach dem Bereichsnamen sortiert. Will man eine andere Sortierung haben, kann man in jedem Bereich in der .wssg/section.yaml den Eintrag order entsprechend setzen. Auch hier spielt der absolute Wert keine Rolle, d.h. es muss nicht 0,1,2 verwendet werden. Um nachträglich Bereiche einzufügen, kann man auch 10 , 20, 30 für den Start benutzen. So kann man später neuer Bereich bei 15, 25 usw. eingefügt werden, ohne das man alle Bereichskonfigurationen ändern muss.

## Gibt es Cookies und einen Cookiebanner?

Wenn du keine speziellen Seitenscripte oder eigenes HTML verwendest, werden keine Cookies auf der Seite verwendet. Nur ein aktivierter Cookiebanner hinterlässt einen Cookie. Im Standard ist ein Cookiebanner integriert, allerdings abgeschaltet (und wird somit nicht im Ergebnis verwendet). Dieser wird aktiviert, wenn man in der .wssg/siteconfig.yaml folgenden Bereich verwendet: (ändern falls schon vorhanden)

```yaml
cookiebanner:
 enabled: true
 text: 'Diese Seite verwendet Cookies.'
```

In Text sollte dann ein emhr oder weniger sinnvoller Text stehen. Dieser wird in einem Bereich am Ende der Seite angezeigt. Mehr könnt ihr [hier](https://www.conductor.com/de/academy/glossar/cookies/) dazu finden. Ich verwende im Standard diese kleine Bibliothek zur Verwaltung und Anzeige: https://github.com/dobarkod/cookie-banner Der Banner wir nur beim ersten Besuch der Seite angezeigt.
Wie das genau aussieht, kann im Example angesehen werden. 
In jeder Seite wird dabei folgendes Template verwendet:

```html
<script type="text/javascript" id="cookiebanner"
  src="https://cdn.jsdelivr.net/gh/dobarkod/cookie-banner@1.2.2/dist/cookiebanner.min.js"
  data-message="{{.cookiebanner.text}}"></script>
```

Möchtest du deinen eigenen Banner verwenden, dann kannst du in der `.wssg/layout.html`  nach dem Makro `{{.cbanner}}`  suchen (am ende der Datei) und dieses mit deinem eigenen Code ersetzen. 

## Kann ich Bilder in den MD Dateien verwenden?

Ja aber sicher. Es gibt aber ein paar Regeln zu den Bildern. 

1. Die Bilder müssen im Seitenordner, da wo alle deine Markdown Dateien liegen. Am besten wäre ein eigener Imageordner für die Bilder.

2. Der Pfad zu den Bildern muss relativ vom Root-Ordner angegeben werden.
   Beispiel (aus dem Example):  

   ```markdown
   ![licht](images/licht.jpg)
   ```

## Kann ich eigene Seiten mit HTML verwenden?

Ja, du kannst HTML Dateien verwenden. Du kannst sogar Makros wie `{{.page.URLPath}}` oder `{{.site.Description}}`.  Im Example unter html (`html/index.html`) findest du mehr dazu.

## Gibt es auch eine Kontaktseite?

Jein, Im Beispiel gibt es eine Kontaktseite basierend auf https://web3forms.com/. Für eine Kontaktseite werden ja verschiedene Dinge benötigt. Das geht von ganz einfachen mailto- Links: 

```html
<a href="mailto:name@bla.de">Sende eine E-Mail</a>
```

, die tatschlich nur einen Mailclient auf dem Rechner des Benutzer benötigt, bis zu komplexen Formularen. Diese brauchen natürlich ein Backend, also einen Server, der die Kontaktanfragen weiter verarbeitet. Es gibt dafür viele Anbieter und dort kann man sich dann meistens eine entsprechende HTML Seite generieren lassen. 

## Serv meldet Fehler in Galerie, was kann ich tun?

Serv meldet folgender Fehler und die Seite wird nicht aktualisiert.

```text
2024/02/14 11:26:23 Error: generator: error processing site: &{%!V(string=yaml: line 14: mapping values are not allowed in this context)}
2024/02/14 11:26:23 Error: server: error generate site: yaml: line 14: mapping values are not allowed in this context
```

Meistens wird das Problem durch ein Sonderzeichen in einem der Properties ausgelöst.

Hier steht z.B. ein Duppelpunkt in der description.

```yaml
upl-40:
    description: Größen: 7x10cm, 12x12cm, 13x15cm, 17x17cm, 27x30cm
    title: Stoffbeutelchen
```

Lösung ist einfach den ganzen Eintrag in Doppelte Hochkommata setzen.

```yaml
upl-40:
    description: "Größen: 7x10cm, 12x12cm, 13x15cm, 17x17cm, 27x30cm"
    title: Stoffbeutelchen
```

## Kann ich den Output Ordner löschen (lassen)?

Ja, wenn du gerne mal den Output Ordner löschen und neu generieren möchtest, kannst du natürlich vor dem `wssg generate` den Ordner löschen. Oder du benutzt das `-c (--clear)` Flag. Damit wird sowohl vor dem `generate` wie auch beim `serve` das Ausgabeverzeichnis (Default: `.wssg/output`) einmal gelöscht. 

Du kannst auch das `output` Verzeichnis während des `wssg serve` löschen. `wssg` wird dann automatisch das `generate` ausführen.

## Kann ich einen eigenen Prozessor bauen?

Mit ein bisschen Programmiererfahrung in Golang kannst du natürlich einen eigenen Prozessor bauen. 

Als erstes solltest du das Projekt von Github forken. ([wssg](http://github.com/willie68/wssg))

Dann kannst du dir deinen Fork auschecken und bearbeiten. Prozessoren liegen im Prozessor Ordner. Jeder Prozessor hat hier einen eigenen Unterordner.  Am einfachsten ist es den Plain Prozessor zu kopieren und entsprechend umzubenennen.  Im `init()` des Package wird dann der Prozessor in dem Dependencies Injection framework registriert. (`do.ProvideNamedValue`) Damit diese Init Funktion auch angesprochen wird, muss in der `main.go` (Root Ordner) noch ein anonymer Import auf dein neues Package gemacht werden.  Ab jetzt steht dein Prozessor mit dem Namen zu Verfügung, den du ihm in der `Name()` Methode mit gegeben hast. Wenn du dann fertig bist, kannst du gerne einen PR stellen, dann kann ich evtl. deinen Prozessor mit in den wssg Standard übernehmen.

## Wie formatiere ich ein Datum in meiner Sprache in einem Blogeintrag? 

Um das Erzeugungsdatum eine Blogeintrages als Makro zu benutzen dient das `.created` Makro.  Die normale Form ist aber für die Anzeige kaum zu gebrauchen. (  `2024-02-18 11:35:06.7435846 +0100 CET` ) 
Im Blog Prozessor steht deswegen eine eigene Formatierungsfunktion zur Verfügung. Diese wird folgendermassen benutzt:

```
{{ dtFormat .created "Monday, 2.01.06" "de_DE" }}
```

`dtFormat` ist der Name der Funktion, dann folgen 3 Parameter.  
`.created` ist das Feld, was den Datumstempel enthält. 
`"Monday, 2.01.06"` ist das Layout, wie das Datum als Text dargestellt werden soll. Im Anhang findet man eine kleine Hilfe dazu.
`"de_DE"` ist das Gebietsschema. Die wichtigsten sind: `"en_US", "en_GB", "de_DE", "fr_FR", "es_ES"`. Weitere Schemata auf Anfrage.

## Kann ich aus den Output ein Dockerimage mit HTTP-Server bauen?

Aber klar, das ist sogar sehr einfach. Zunächst musst du im .wssg Ordner eine Datei `Dockerfile` mit folgenden Inhalt erstellen. 

```dockerfile
FROM httpd:2.4
COPY ./output/ /usr/local/apache2/htdocs/
```

Ich verwende sehr gerne den Apache httpd Server.

Zu, Erzeugen und starten des Images rufst du dann auf: (mein-image kannst du gerne gegen was anderes tauschen)

```shell
docker build -t mein-image-httpd .
docker run --name mein-image -p 9080:80 -d mein-image-httpd
```

Nach dem Build und Startprozess hast du nun einen neuen Container mit deiner Webseite laufen. Diesen kannst du jetzt unter `http://localhost:9080` erreichen.

# Anhang

## Golang Datumsformat

| Go layout | Format | Example  | Description     |
| --------- | ------ | -------- | --------------- |
| 2006      | YYYY   | `"2022"` | Four-digit year |
| 06        | YY     | `"22"`   | Two-digit year  |

| Go layout | Format | Example  | Description                                     |
| --------- | ------ | -------- | ----------------------------------------------- |
| January   | MMMM   | `"July"` | Full month name                                 |
| Jan       | MMM    | `"Jul"`  | Three-letter abbreviation of the month          |
| 01        | MM     | `"07"`   | Two-digit month (with a leading 0 if necessary) |
| 1         | M      | `"7"`    | At most two-digit month (without a leading 0)   |

| Go layout | Format | Example     | Description                                                  |
| --------- | ------ | ----------- | ------------------------------------------------------------ |
| Monday    | DDDD   | `"Tuesday"` | Full weekday name                                            |
| Mon       | DDD    | `"Tue"`     | Three-letter abbreviation of the weekday                     |
| 02        | DD     | `"08"`      | Two-digit month day (with a leading 0 if necessary)          |
| _2        | _D     | `" 8"`      | Two-character month day with a leading space if necessary    |
| 2         | D      | `"8"`       | At most two-digit month day (without a leading 0)            |
| 002       | ddd    | `"074"`     | Three-digit day of the year (with a leading 0 if necessary)  |
| __2       | __d    | `" 74"`     | Three-character day of the year with a leading spaces if necessary |

### Time format

| Go layout | Format | Example | Description                                               |
| --------- | ------ | ------- | --------------------------------------------------------- |
| 15        | hh     | `"17"`  | Two-digit 24h format hour                                 |
| 03        | hh     | `"05"`  | Two digit 12h format hour (with a leading 0 if necessary) |
| 3         | h      | `"5"`   | At most two-digit 12h format hour (without a leading 0)   |
| PM        | am/pm  | `"AM"`  | AM/PM mark (uppercase)                                    |
| pm        | am/pm  | `"am"`  | AM/PM mark (lowercase)                                    |

| Go layout | Format | Example | Description                                      |
| --------- | ------ | ------- | ------------------------------------------------ |
| 04        | mm     | `"07"`  | Two-digit minute (with a leading 0 if necessary) |
| 4         | m      | `"7"`   | At most two-digit minute (without a leading 0)   |

| Go layout              | Format | Example        | Description                                      |
| ---------------------- | ------ | -------------- | ------------------------------------------------ |
| 05                     | ss     | `"09"`         | Two-digit second (with a leading 0 if necessary) |
| 5                      | s      | `"9"`          | At most two-digit second (without a leading 0)   |
| .0, .00, …, .000000000 | .s     | `".126284000"` | A fractional second (trailing zeros included)    |
| .9, .99, …, .999999999 | .s     | `".126284"`    | A fractional second (trailing zeros omitted)     |

### Time zone format

| Go layout | Format         | Example       | Description                                                  |
| --------- | -------------- | ------------- | ------------------------------------------------------------ |
| MST       | TTT            | `"CEST"`      | Abbreviation of the time zone                                |
| -070000   | ±hhmmss        | `"+010000"`   | Numeric time zone offset with hours, minutes, and seconds    |
| -07:00:00 | ±hh:mm:ss      | `"+01:00:00"` | Numeric time zone offset with hours, minutes, and seconds separated by colons |
| -0700     | ±hhmm          | `"+0100"`     | Numeric time zone offset with hours and minutes              |
| -07:00    | ±hh:mm         | `"+01:00"`    | Numeric time zone offset with hours and minutes separated by colons |
| -07       | ±hh            | `"+01"`       | Numeric time zone offset with hours                          |
| Z070000   | Z or ±hhmmss   | `"+010000"`   | Like -070000 but prints `"Z"` instead of `"+000000"` for the UTC zone (ISO 8601 behavior) |
| Z07:00:00 | Z or ±hh:mm:ss | `"+01:00:00"` | Like -07:00:00 but prints `"Z"` instead of `"+00:00:00"` for the UTC zone (ISO 8601 behavior) |
| Z0700     | Z or ±hhmm     | `"+0100"`     | Like -0700 but prints `"Z"` instead of `"+0000"` for the UTC zone (ISO 8601 behavior) |
| Z07:00    | Z or ±hh:mm    | `"+01:00"`    | Like -07:00 but prints `"Z"` instead of `"+00:00"` for the UTC zone (ISO 8601 behavior) |
| Z07       | Z or ±hh       | `"+01"`       | Like -07 but prints `"Z"` instead of `"+00"` for the UTC zone (ISO 8601 behavior) |
