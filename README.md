# PassVault
Minimalistic TUI (BubbleTea) password manager that uses SQLite database for encrypted password storage. KDF (Argon2) is used to derive KEK from user-supplied password, which is then used to decrypt DEK stored in database. DEK is stored in RAM and used for cryptographic operations performed using AES-GCM.

![](https://raw.githubusercontent.com/skiba-mateusz/PassVault/main/demo.gif)

## Features
* **Minimalist TUI Based on BubbleTea** - fast and smooth terminal-based service
* **Local SQLite DB** - stores only encrypted data
* **Strong AES-GCM Encryption** - protects both confidentiality and integrity
* **Secure Key Architecture** - Password -> Argon2 -> KEK -> DEK -> AES-GCM
* **List, Add, Delete and Edit** - from the command line interface
* **Lack of Dependece on External Services** - everything works offline

## License
MIT License
