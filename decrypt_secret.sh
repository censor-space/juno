#!/bin/sh

# ファイルを復号化
# --batchでインタラクティブなコマンドを防ぎ、
# --yes で質問に対して "はい" が返るようにする
gpg --quiet --batch --yes --decrypt --passphrase="$LARGE_SECRET_PASSPHRASE" \
--output eguchi-wedding-firebase-adminsdk.json eguchi-wedding-firebase-adminsdk.json.gpg
gpg --quiet --batch --yes --decrypt --passphrase="$LARGE_SECRET_PASSPHRASE" \
--output eguchi-wedding-google-app-engine-service-account.json eguchi-wedding-google-app-engine-service-account.json.gpg
