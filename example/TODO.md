# mit Michael besprechen

SIcherheitsloch: Anmeldung ... "Benutzer nicht gefunden!"

# mit Gemini besprechen

# mit DM Team besprechen

Eigentlich war der Sinn dahinter, daß man eben keine Func mehr schreiben muss. 
Der Service implementiert ein Check Interface (Checkable, oder sowas) und wird beim wireup einfach beim Healthservice registriert. (Wo das wireup stattfindet, ist nochmal eine andere Diskussion, main finde ich da nicht optimal, weil schlecht testbar. Und wir müssen eh den Injector mit geben, denn der NIL Injector macht beim Testen Probleme.)
Dabei entfallen dann einfach genau diese "Brücken" Funktionen. 
btw.: Eigentlich würde ich gerne das HealthCheckable Interface aus samber.Do verwenden. 
Dann entfiele das ganze selber registrieren. Einfach nur per do.Provide zur Verfügung stellen, fertig

# sonstiges

- Scan Prozess Schritt 2 in Tasks umwandeln.
  https://intranet.easy.de/display/~w.klaas/Scan+Prozess

## JSON

https://jqlang.github.io/jq/manual/

https://github.com/itchyny/gojq




## Architektur

### todo

### CA:

- eigene CA, services sollten dann nur CA Signierte Zertifikate verwenden.
- Hashicorp Vault
- mcs microvault
- eigene CA mit openssl: https://deliciousbrains.com/ssl-certificate-authority-for-local-https-development/#private-key-root-certificate-macos-linux

#### Key Management

- https://square.github.io/keywhiz/
  - MySQL as database
  - simple JSON /REST interface
  - wird leider seit 2019 nicht mehr gepflegt
- https://www.vaultproject.io/
  - HashiCorp Vault

### Lizenz- und Nutzungsdatenerfassung:

Eigenbau:

MongoDB Timeseries Collection

Graphite: https://graphiteapp.org

https://github.com/VictoriaMetrics/VictoriaMetrics

https://thanos.io/

https://m3db.io/

Fertige Lösungen:
https://www.racknap.com/billing-and-pricing-management/

https://www.financialforce.com/learn/billing/cloud-billing-software/  (Salesforce basierend)

https://www.nitrobox.com

https://rev.io/

### Suchengines

OpenSearch: https://opensearch.org/

Solr: https://solr.apache.org/



# IT

