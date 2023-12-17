# WSSG

Willie's Static Site Generator

# Warum?

Darum. Ich finde jeder Programmierer sollte so ein paar Sachen mal programmiert haben. Hello Wlord, Sortieralgorythmus, Datenbank, Website Generator. Letzerer fehlte mir noch.

# Basis

wssg ist ein kleiner feiner Website generator mit folgenden Features (wenn er denn fertig ist, und was ist schon fertig...)

- Basis sind markdown Dateien
- automatische Indexerstellung
- eingebauter WebServer zum besseren Testen der Site
- austauschbare Template Engine

# Installation

wssg ist ein typisches Copy/Run Programm, deswegen auch in Golang geschrieben. Es kann aber auch mithilfe von Golang installiert werden. 

## Copy/Run

Nimm einfach die richtige Releaseversion und kopiere sie in ein Verzeichniss, daß bereits in deinem Pfad vorhanden ist. Oder kopiere das Binary in das root der neuen Site. 

Zum Testen ob alles funktioniert gibt einfach `wssg version` ein.

## via Golang

`go install github.com/willie68/wssg`

# Quickstart

Für den schnellen Start mit dem wssg wird zunächst einen Installation vorausgesetzt.

Als erstes erzeugen wir uns eine neue Side. 

`wssg init ./<sitename>`

Jetzt wird automatisch ein Verzeichniss mit dem Namen <sitename> erzeugt und dort dann alle wichtigen Konfiguration erzeugt. Auch wird direkt eine erste Seiten (index.md) mit erzeugt.

`wssg generate`  

generiert nun die Website. Das Ergebniss landet automatisch m Ordner `.wssg/output`

# Aufbau

Das Programm ist für folgende Struktur am besten geeignet. Die erste Ebene (Root) ist quais der Startpunkt. Dort muss für den Start eine index.md erstellt werden. Diese wird automatisch beim `wssg init` angelegt. Hier können dann weitere Seiten (pages) hinzugefügt werden. Zusätzliche Dateien, wie z.B. Bilder, Stylesheets, JS usw. können sowohl in den Rootordner wie auch in weiteren Unterordnern abgelegt werden. Unterordner können dann einfach per relativer Angabe referenziert werden. 

Möchte man einen neuen Bereich (section) mit verschiedene Seiten anlegen, kann man das mit `wssg new section <name>` machen. Dabei wird nun, ebenso wie im root Ordner, ein Unterverzeichnis .wssg mit den Einstellungen für diesen Bereich erstellt.  

# Programmparameter

## Variablen für eine Seite

`{{.body}}` ergibt den konvertierten Text aus der Markdown Datei.

`{{.site.#}}` sind die Einstellungen für die gesamte Website. Hier stehen 1:1 alle Einstellungen aus der `siteconfig.yaml`. Beispielsweise  

`{{.site.language}}` ergibt die Sprache

`{{.site.title}}` den Webseitentitel. Ebenso funktionieren `{{.site.description}}` und `{{.site.keywords}}`

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
