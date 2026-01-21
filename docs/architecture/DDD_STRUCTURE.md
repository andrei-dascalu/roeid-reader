# DDD Architecture Summary

## Overview

The roeid-reader project is organized using Domain-Driven Design (DDD) principles. The codebase
is structured around five **bounded contexts**, each representing a distinct domain with its own
language, models, and responsibilities.

## Directory Structure

```bash
roeid-reader/
├── cmd/roeid-reader/
│   └── main.go                       # Entry point
├── internal/
│   ├── smartcard/                    # Bounded Context 1: Smart Card Operations
│   │   ├── domain/                   # Domain models (APDU, Response, Card, Status)
│   │   ├── infrastructure/           # PC/SC transport, logging
│   │   └── application/              # Service orchestration
│   ├── pace/                         # Bounded Context 2: PACE Protocol
│   │   ├── domain/                   # Password, Nonce, KeyAgreement, Auth
│   │   └── application/              # PACE orchestrator
│   ├── crypto/                       # Bounded Context 3: Cryptography
│   │   ├── domain/                   # EC, AES, KDF interfaces
│   │   └── infrastructure/           # Brainpool, CMAC, ECDH implementations
│   ├── messaging/                    # Bounded Context 4: Secure Messaging
│   │   ├── domain/                   # SecureMessage, SSC
│   │   └── application/              # Encryption/decryption service
│   └── carddata/                     # Bounded Context 5: Card Identity Data
│       ├── domain/                   # Identity, IdentityRepository
│       └── application/              # Data reading and parsing
└── docs/
    ├── architecture/
    │   ├── BOUNDED_CONTEXTS.md       # Detailed context descriptions
    │   └── DOMAIN_MODELS.md          # Entity relationships & invariants
    └── PACE_GLOSSARY.md              # Protocol terminology

```

## Five Bounded Contexts

| Context | Responsibility | Key Models | Status |
| --- | --- | --- | --- |
| **Smart Card** | PC/SC transport, APDU messaging | APDU, Response, Card, Status | ✅ Migrated |
| **PACE** | Cryptographic authentication protocol | Password, Nonce, KeyAgreement | 🔧 Skeleton |
| **Crypto** | ECC, AES, KDF primitives | EllipticCurve, AESKey, KDF | 🔧 Skeleton |
| **Messaging** | Encrypted APDU layer | SecureMessage, SSC | 🔧 Skeleton |
| **Card Data** | Personal identity records | Identity, IdentityRepository | 🔧 Skeleton |

## Architecture Principles

1. **Separation of Concerns:** Each context handles one domain responsibility
2. **Rich Domain Models:** Models contain both data and behavior (not anemic DTOs)
3. **Infrastructure Isolation:** External libs (PC/SC, crypto) live in infrastructure layer
4. **Application Services:** Orchestrate domain logic across aggregates
5. **No Direct Context Imports:** Contexts communicate via interfaces only
6. **Dependency Injection:** Explicit dependencies passed to services

## Development Roadmap

1. ✅ **Phase 0:** Directory structure & migration foundation
2. ✅ **Phase 1:** SmartCard context (working code from `main.go`)
3. 🔧 **Phase 2:** Implement remaining contexts (PACE, Crypto, Messaging, CardData)
4. 📝 **Phase 3:** Integration tests & documentation

## Key Files

- [BOUNDED_CONTEXTS.md](BOUNDED_CONTEXTS.md) - Detailed context specifications
- [DOMAIN_MODELS.md](DOMAIN_MODELS.md) - Entity diagrams and relationships
- [cmd/roeid-reader/main.go](../../cmd/roeid-reader/main.go) - Entry point
- [IMPLEMENTATION_ROADMAP.md](../../IMPLEMENTATION_ROADMAP.md) - Task tracking

## Running the Application

```bash
cd /Users/andrei/Projects/personal/roeid-reader
go build -o roeid-reader ./cmd/roeid-reader
./roeid-reader
# Interactive PIN prompt; card must be inserted
```

## Next Steps

1. Implement PACE protocol phases (Phase 3-9 of roadmap)
2. Resolve Brainpool ECC curve dependency
3. Establish secure messaging channel
4. Read protected identity data
5. Add integration tests with real CEI card

## Architecture Reference

- **Book:** "Domain-Driven Design" by Eric Evans
- **Pattern:** Hexagonal Architecture (Ports & Adapters)
- **Example:** Bounded contexts with clear domain language, infrastructure abstraction, application services
