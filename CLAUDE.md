# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

**roeid-reader** is a Go utility for reading Romanian electronic ID (CEI) smart cards via PC/SC interface.
It implements the PACE (Password Authenticated Connection Establishment) cryptographic protocol
following BSI TR-03110 specification.

## Build & Development Commands

```bash
# Build
go build -o roeid-reader ./cmd/roeid-reader

# Run (requires smart card reader + CEI card inserted)
./roeid-reader

# Run tests
go test ./...

# Run single test
go test -v -run TestName ./internal/smartcard/...

# Go linting & formatting
go fmt ./...
go vet ./...

# Markdown linting
npm run lint:md
npm run lint:md:fix
```

## Architecture

This project follows **Domain-Driven Design** with **hexagonal architecture** (ports & adapters),
organized into five bounded contexts:

```text
internal/
├── smartcard/     # BC1: PC/SC transport & APDU messaging (working)
├── pace/          # BC2: PACE cryptographic protocol (skeleton)
├── crypto/        # BC3: ECC, AES, KDF primitives (skeleton)
├── messaging/     # BC4: Encrypted APDU layer (skeleton)
└── carddata/      # BC5: Identity data reading & parsing (skeleton)
```

Each bounded context follows the same layered structure:

- `domain/` - Interfaces, value objects, entities (no external deps)
- `application/` - Service orchestration (depends on domain + infrastructure)
- `infrastructure/` - External library wrappers, adapters

**Key principle:** Contexts communicate via interfaces only, never direct imports.

## Critical Patterns

### APDU Handling

- APDU byte sequences follow **ISO/IEC 7816-4** strictly
- Format: `[CLA] [INS] [P1] [P2] [Lc] [Data...] [Le]`
- Always check status words (SW1/SW2): `0x9000` = success
- Use domain `StatusError` type for APDU errors, not generic Go errors

### Cryptographic Standards

- **BSI TR-03110 v2.13**: PACE protocol specification
- **NIST SP 800-38B**: AES-CMAC computation
- **BrainpoolP256r1**: Required elliptic curve (not in Go stdlib)

### Security

- Cryptographic keys (`AESKey`, `Password`) have `Clear()` methods for memory cleanup
- `K_pi` (password-derived key) never transmitted; stays local
- `SendSequenceCounter` must increment BEFORE each encryption operation

## Key Dependency

`github.com/ebfe/scard` - PC/SC smart card interface wrapper for reader communication.

## Contributing Guidelines

- Preserve APDU byte sequences exactly; reference ISO/IEC 7816-4 in comments
- Add new infrastructure functions to the appropriate context's `infrastructure/` layer
- For crypto implementations, validate against TR-03110 or NIST test vectors
- Update `PLAN.html` and `IMPLEMENTATION_ROADMAP.md` when modifying protocol flow
- Reference standards (e.g., `// BSI TR-03110 Section 4.2`) in comments
- follow Go 1.25 best practices
- logging should be done using zerolog instead of log stdlib or fmt

## Implementation Status

See `IMPLEMENTATION_ROADMAP.md` for the 13-phase roadmap. Currently:

- Phase 1 complete: SmartCard context working (PC/SC connection, APDU exchange)
- Phase 2-13 planned: PACE protocol, crypto primitives, secure messaging, identity reading

**Known blocker:** Brainpool256r1 ECC support requires third-party library evaluation.

## Development Workflow

**Before implementing any feature:**

1. Read `IMPLEMENTATION_ROADMAP.md` to identify the current active phase
2. Verify all prior phases are complete (all checkboxes checked)
3. Do not implement features from future phases - phases have dependencies
4. Check for blockers noted in the roadmap (e.g., Brainpool ECC for Phase 5+)

**Phase dependencies (do not skip):**

- Phase 2-3 (APDU layer, app selection) → required before any PACE work
- Phase 4 (password processing) → required before Phase 6 (mapping uses K_pi)
- Phase 5 (Brainpool ECC) → required before Phase 6-7 (mapping and key agreement)
- Phase 6-9 (PACE protocol) → required before Phase 10 (secure messaging)
- Phase 10 (secure messaging) → required before Phase 11 (reading identity data)

**When completing work:**

1. Mark completed tasks with `[x]` in `IMPLEMENTATION_ROADMAP.md`
2. Verify the phase deliverable is met before moving to the next phase
3. Add tests for completed functionality where possible
4. Update this file's "Implementation Status" section when a phase completes

## Documentation

- `docs/architecture/DDD_STRUCTURE.md` - Architecture overview
- `docs/architecture/BOUNDED_CONTEXTS.md` - Context specifications
- `docs/architecture/DOMAIN_MODELS.md` - Entity relationships
