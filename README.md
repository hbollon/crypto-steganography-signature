# crypto-steganography-signature 
### TP 2 BOLLON Hugo

Ce TP à été réalisé entièrement et est fonctionnel.
Un petit CLI à été préparé pour un test simplifié et le code commenté entièrement (en anglais) pour une compréhension aisé.
- Utilisation du CLI: 
```bash
Usage: go run . <COMMAND> <ARGS>
Commands:
        generate-custom-diplome <NAME> <GRADE> <IMG_PATH> <KEYPAIR_BIT_SIZE>
        extract-lsb-from-diplome <IMG_PATH>
        verify-signature <ORIGINAL_DATA> <ENCODED_DATA>
```

Il faut bien entendu avoir au préalable installé Go.
Pour vérifier le message encodé, il faut obligatoirement le connaitre et ensuite utiliser la commande verify-signature avec ce dernier et le message encodé.
Il aurait, sinon, fallut utiliser la clé privé.
