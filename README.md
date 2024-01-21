# Distributed architecture project

[lien du google doc pour s'organiser](https://docs.google.com/document/d/1eTxQAxZicfh38igHVdozKWNw_PwN3aOlUZEojc4ndGA/edit#heading=h.i4yrdggs5m7f)

[lien kanban/jira](https://fil2026.atlassian.net/jira/software/projects/KAN/boards/1)

lancer le docker influxDB :
docker compose --env-file  configs/influxdb.env up

generer le swagger :
- check si on bien le go/bin dans le $PATH (export PATH=$PATH:$GOPATH/bin)
- swag init -g cmd/api/main.go à la racine

Groupe members :
- Tom Freret
- Louis Painter
- Matthis Bleuet
- Antoine Otegui
