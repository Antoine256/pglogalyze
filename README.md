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