## Description

Ascii-art-web est un site internet qui vas prendre une chaine de caractere et avas l'afficher en art en fonction de divers polisse . 

## author

realiser entierement par Savary Swann .

## utilisation

dans le dossier ascii-art-web lancer un "go run . "
ensuite aller sur http://localhost:8080 

## detail

Le programme lit un fichier de police ASCII (ex: standard.txt), où chaque
caractère imprimable (de l'espace à '~') est représenté par 8 lignes de texte 
plus une vide pour la separation soit 9 ligne au total.

Pour chaque caractère du texte saisi par l'utilisateur, le programme calcule 
sa position dans le fichier de police grâce à son code ASCII, puis extrait 
les 8 lignes correspondantes. Ces lignes sont ensuite concaténées 
horizontalement pour former le dessin final, caractère par caractère.

Le serveur HTTP expose deux routes : GET / affiche la page d'accueil 
et POST /ascii-art traite le texte et la police choisis 
pour générer et afficher le résultat.

## Docker

# pour build 
docker build -t ascii-art-web:2.0 .
# pour le lancer 
docker run -d -p 8080:8080 --rm --name ascii-web \
  --label fr.zone01.titre="ascii-art-web container" \
  --label fr.zone01.version="2.0" \
  ascii-art-web:2.0

# pour verifier

docker ps
docker logs ascii-web
docker inspect --format '{{json .Config.Labels}}' ascii-web
docker inspect --format '{{json .Config.Labels}}' ascii-art-web:2.0

# l'arreter

docker stop ascii-web

# nettoyer 
docker rmi ascii-art-web:1.0
docker image prune
docker system df