# Revaskel
Après une restauration, PostgreSQL démarre, le stockage répond et le site affiche sa page d’accueil. Pourtant, l’agent logistique ne peut plus télécharger le document associé à une expédition. La base et les fichiers ont été restaurés à deux instants différents. Le tableau de bord technique est vert, mais l’expédition reste bloquée.

Le CLI Go vérifie qu’une expédition d’un instantané retrouve son document et que son empreinte correspond à celle enregistrée. Le dernier point cohérent est le dernier élément du préfixe continu avant la première rupture.

## Démarrage

Depuis la racine du dépôt avec Go 1.22 ou ultérieur et Make :

```sh
make demo
make test
```

Commande directe :

```sh
go run . examples/complete
```

Le premier instantané est cohérent. Dans `examples/partial`, le document de `DEMO-002` manque. Le dernier point cohérent reste le 1er janvier 2026 à 08:00 UTC.

Le binaire renvoie 0 pour une reprise cohérente, 1 pour un document manquant ou divergent, 2 pour une erreur de lecture ou de format.

## Organisation

- `main.go` contient le moteur et son interface en ligne de commande.
- `examples/` contient les entrées fictives et reproductibles.
- `main_test.go` couvre le comportement nominal et les cas limites.

## Périmètre du socle

Ce premier adaptateur lit un répertoire local de confiance contenant un manifeste ordonné et des objets ordinaires. Les chemins symboliques et les archives externes ne font pas partie du format pris en charge. L’exercice ne restaure pas encore de service. Une opération absente du manifeste entier échappe à ce contrôle.

## Premier chantier à reprendre

Restaurer une petite base PostgreSQL et un répertoire objet dans un environnement isolé. Produire le manifeste depuis le registre des opérations attendues et tester la lecture d’un document avec les droits de l’opérateur.
