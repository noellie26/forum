# forum-Communication de base

Ce projet consiste à créer un forum web qui permet :
- la communication entre les utilisateurs (messages, commentaires, préférences/aversions);
- associer des catégories aux publications (pour les utilisateurs connectés lors de la création d’une nouvelle publication);
- aimer et détester les messages et les commentaires (utilisateurs connectés);
- filtrage des messages (utilisateurs connectés).

### Stockage des données

Afin de stocker les données dans ce forum (comme les utilisateurs, les messages, les commentaires, etc.), la bibliothèque de base de données 
[SQLite](https://www.sqlite.org/index.html) est utilisé.

Les requêtes SELECT, CREATE et INSERT sont utilisées.

### Authentification

Le client peut enregistrer un nouvel utilisateur sur le forum en saisissant ses identifiants. 
Une session de connexion est créée pour accéder au forum et pouvoir ajouter des messages et des commentaires.

Les cookies sont utilisés pour permettre à chaque utilisateur d’avoir une seule session ouverte. Chacune de ces sessions contient un 
date d’expiration (24h). C’est à vous de décider combien de temps le cookie reste "vivant". UUID est utilisé comme ID de session.

Instructions pour l’enregistrement de l’utilisateur :
- Un courriel est requis :
  - Lorsque l’e-mail est déjà pris, une réponse d’erreur est renvoyée;
- Nom d’utilisateur requis :
  - lorsque le nom d’utilisateur est déjà pris, une réponse d’erreur est renvoyée;
- Mot de passe requis :
  - Le mot de passe est chiffré lorsqu’il est stocké.

### Communication

Afin que les utilisateurs puissent communiquer entre eux, ils peuvent créer des messages et des commentaires.

- seuls les utilisateurs enregistrés peuvent créer des publications et des commentaires;
- lorsque les utilisateurs enregistrés créent une publication, ils peuvent y associer une ou plusieurs catégories;
- la mise en œuvre et le choix des catégories (tags) étaient à la charge des développeuses
- les publications et les commentaires sont visibles par tous les utilisateurs (enregistrés ou non)
- Les utilisateurs non inscrits ne peuvent voir que les publications et les commentaires.

### Aime et n’aime pas

Seuls les utilisateurs enregistrés peuvent aimer ou détester les publications et les commentaires.

Le nombre de likes et d’aversions est visible par tous les utilisateurs (Enregistrés ou non).

### Filtrer

Un mécanisme de filtrage a été mis en place, qui permettra aux utilisateurs de filtrer les publications affichées par les catégories ;

### Téléchargement d’image et de video

Dans le téléchargement d’images/video de forum, les utilisateurs enregistrés ont la possibilité de créer un message contenant une image/video ainsi que du texte.

- Lorsque vous consultez le message, les utilisateurs et les invités peuvent voir l’image/video qui lui est associée.
Dans ce projet, les types JPG, JPEG, PNG et GIF pour l'image et WEPM, MP4, MOV, AVI et MKV pour la video sont gérés.

La taille maximale des images à charger est de 20 Mo et celle des videos est de 100Mo. S’il y a une tentative de chargement d’une image supérieure à 20 Mo ou d'une videos supérieure à 100Mo, 
un message d’erreur indique à l’utilisateur que l’image ou la video est trop grande.

### Fonctions avancées

Les utilisateurs sont informés (page « Notifications ») lorsque leurs messages sont :
- aimé/détesté (commentaires également)
- commenté

À la page « Mon activité », l’utilisateur peut suivre sa propre activité :
- Affiche les publications créées par l’utilisateur
- Indique où l’utilisateur a laissé un like ou un dislike
- Indique où et ce que l’utilisateur a commenté. Pour cela, le commentaire devra être affiché, ainsi que le post commenté

L’utilisateur peut modifier/supprimer ses messages et commentaires.

### Modération

La modération du forum est basée sur un système de modération. Il présente un modérateur qui, selon le niveau d’accès d’un utilisateur, 
approuve les messages affichés avant qu’ils ne deviennent visibles publiquement.

Le filtrage peut se faire en fonction des catégories du post triées par non pertinent, obscène, illégal ou insultant.

Il y a au total 4 types d’utilisateurs :
- Invités
  - Il s’agit d’utilisateurs non enregistrés qui ne peuvent ni publier, ni commenter, ni aimer ou détester un message. 
  Ils ont seulement la permission de voir ces messages, commentaires, goûts ou aversions.
- Utilisateurs
  - Ce sont les utilisateurs qui pourront créer, commenter, aimer ou détester les messages.
- Modérateurs (utilisateur : « Mod», mot de passe : « psw »)
  - Les modérateurs, comme expliqué ci-dessus, sont des utilisateurs qui ont accès à des fonctions spéciales :
    - Ils peuvent surveiller le contenu du forum en supprimant ou en signalant une publication à l’administrateur.
  - Pour créer un modérateur, l’utilisateur doit demander un administrateur pour ce rôle.
- Administrateurs (utilisateur : « admin », mot de passe : « psw »)
  - Utilisateurs qui gèrent les détails techniques nécessaires à l’exécution du forum. Cet utilisateur peut :
    - Promouvoir ou rétrograder un utilisateur normal vers ou depuis un utilisateur modérateur.
    - Recevoir les rapports des modérateurs. Si l’administrateur reçoit un rapport d’un modérateur, il peut y répondre.
    - Supprimer les publications et les commentaires
    - Gérer les catégories, en étant capable de les créer et de les supprimer.

### Docker

Pour le projet forum Docker est utilisé.

Comment :
- Créez l’image Docker en exécutant la commande suivante : “./docker/build.sh”.
- Une fois l’image construite, vous pouvez exécuter un conteneur basé sur l’image en utilisant la commande suivante : “./docker/run.sh”.
- Le conteneur commencera et votre application Go sera accessible à 
[http://localhost:8181](http://localhost:8181) dans votre navigateur Web.
- Pour arrêter et supprimer l’image, exécutez la commande suivante : “./docker/stop.sh”.

Assurez-vous que Docker est installé et exécuté sur votre machine avant de créer et d’exécuter l’image Docker.

### Paquets autorisés

- tous les colis Go standard sont autorisés
- sqlite3;
- bcrypt;
- UUID;

## Développeuses
- kdorgemor
- aaye
- aichsow
- kintuo
- vaberte