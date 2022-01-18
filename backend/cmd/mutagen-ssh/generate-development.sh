#!/bin/bash

mkdir -p keys || true

if [ ! -f "keys/id_ed25519" ]; then
    echo "⭐ Could not find existing keypair, generating..."
    ssh-keygen -o -a 100 -t ed25519 -f keys/id_ed25519 -C "local-development-${USER}-$(hostname)@getsturdy.com" -P ""
    echo "✅ Generated SSH keypair (for development)"
else
    echo "✅ Using existing SSH keypair (for development)"
fi
