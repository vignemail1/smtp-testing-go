# smtp-testing-go

Outil SMTP léger écrit en Go, inspiré de Swaks, sans dépendance externe.

## Fonctionnalités

- Connexion SMTP avec timeout configurable
- `EHLO`/`HELO` personnalisé
- TLS implicite et `STARTTLS`
- Authentification `PLAIN` et `CRAM-MD5`
- Expéditeur et plusieurs destinataires
- En-têtes et corps de message personnalisés
- Mode verbeux
- Mot de passe via `SMTP_PASSWORD`
- Binaires multi-plateformes publiés avec GoReleaser

## Installation

Téléchargez un binaire depuis la page [Releases](https://github.com/vignemail1/smtp-testing-go/releases), ou compilez le projet :

```sh
go build -o smtp-testing-go .
```

## Exemples

Envoi simple :

```sh
./smtp-testing-go -server smtp.example.com:25 \
  -from sender@example.com \
  -to recipient@example.com \
  -subject "Test SMTP" \
  -body "Message de test"
```

Avec authentification et STARTTLS :

```sh
SMTP_PASSWORD='secret' ./smtp-testing-go \
  -server smtp.example.com:587 \
  -from sender@example.com \
  -to recipient@example.com \
  -username sender@example.com \
  -auth plain \
  -starttls \
  -subject "Test" \
  -body "Bonjour"
```

Mode verbeux et plusieurs destinataires :

```sh
./smtp-testing-go -server localhost:2525 \
  -from sender@example.com \
  -to one@example.com,two@example.com \
  -verbose
```

Consultez l'aide avec `-h` pour la liste complète des options.

## Licence

Ce projet est distribué sous licence MIT.
