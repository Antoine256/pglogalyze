# pglogalyze


produit simple => une commande à lancer en psql pour les droits ou sinon spécifié user mdp en option

analyze au présent en direct et/ou analyse des logs passés ?

lancement de la commande : pglogalyze --option1 --option2 parameter ....

Les différentes options :
- lecture en direct (par défaut) ou dans un fichier spécifié sinon 
- spécifié l'instance sinon variables globales ?
- niveau de log (voir si je peux récupérer le niveau de verbosité de postgres s'il y en a un et afficher une erreur si c'est incompatible)
- user mdp
- database spécifique
- type de log à surveiller (connexion, erreur de requête, timeout... voir en fonction des erreurs courantes postgres)
- recherche d'un mot spécifique
(METTRE DES COULEURS)
pour récupérer les logs : https://betterstack.com/community/guides/logging/how-to-start-logging-with-postgresql/

$> pglogalyze --option1 --option2 parameter
-- analyse de l'instance postgres "Nom de l'instance (en paramètre ou par défaut)
surveillance des logs (AVEC LES OPTIONS) des fichiers suivants :
- /data/pgsql/13/log/postgresql.log
- /data/pgsql/15/log/postgresql.log

Logs (LES 10 DERNIERS ?):

- 2023-07-30 08:31:50.628 UTC [2176] user@database listening on IPv4 address "127.0.0.1", port 5432
.
.
.
.
.


## V1

un fichier donné par l'utilisateur => récupérer le paramètre et parser les lignes, donner le nombre d'erreur.


En dev :

go run .\cmd\pglogalyze\main.go ....
exemple : go run .\cmd\pglogalyze\main.go  -f ../../../Desktop/psql/log/postgresql-2026-01-18_124734.log -st 2026-01-18 12:51:00 -et 2026-01-18 14:30:34

Pour BUILD :

//go build -o pglogalyze ./cmd/pglogalyze
//go build -o C:\Users\dupas\Documents\GitHub\pglogalyze ./cmd/pglogalyze

Lancer le script .\UbuntuSharedFolder\build.ps1 depuis la racine de ce projet pour build l'image en mode linux et que le fichier soit mis dans le dossier partager 

ensuite, il suffit de lancer la commande suivante depuis la VM ubuntu :
pglogalyze -f /var/log/postgresql/postgresql-16-main.log -l LOG -et 2026-03-07T11:46:33

pglogalyze -f ../../../Desktop/psql/log/postgresql-2026-01-18_124734.log -st 2026-01-18 12:51:00 -et 2026-01-18 14:30:34


V2 Ajout de fonctionnalitées d'annalyse...(applicatif vs infra ??), lecture en direct, database spécifiquement
type de recherche: statement (aussi une severity), connection, duration, checkpoint, starting PostgreSQL and shutting down

V3 définir l'os, dossier de log et le fichier le plus récent (ou en fonction de la date) et voir si l'user a le droit de lecture, voir si possible de récupérer via les infos du services postgres mais peut être compliqué en fonction des droits user.
 
mettre une limite de ligne par défaut et en paramètre.