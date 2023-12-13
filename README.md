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



# Programmparameter
