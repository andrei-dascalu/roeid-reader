# DDD Architecture - Bounded Contexts

This document details the five bounded contexts in roeid-reader's Domain-Driven Design architecture.

---

## 1. Smart Card Context (`internal/smartcard/`)

### Purpose

Manage all direct interaction with physical smart cards via the PC/SC interface.

### Domain Models

- **Card:** Interface defining transmit/disconnect operations
- **APDU:** ISO/IEC 7816-4 command capsule (CLA, INS, P1, P2, data)
- **Response:** APDU response with data and status words (SW1/SW2)
- **CardStatus:** ATR, active protocol, reader information
- **Reader:** Smart card reader abstraction
- **StatusError:** Domain-specific error type with status code interpretation

### Key Behaviors

- APDU serialization to ISO/IEC 7816-4 format
- Status word (0x9000 = success) validation
- Resource lifecycle management (connect/disconnect)

### Infrastructure

- **PCSCTransport:** PC/SC binding using `github.com/ebfe/scard`
- **APDULogger:** Logs all APDU exchanges with timestamps

### Application Service

- **SmartCardService:** Orchestrates connection, application selection, PIN verification
  - Methods: `Connect()`, `Disconnect()`, `SelectApplication()`, `VerifyPIN()`, `Transmit()`

### Dependencies

- External: `github.com/ebfe/scard`
- Internal: None (lowest-level context)

---

## 2. PACE Context (`internal/pace/`)

### Purpose

Implement the PACE (Password Authenticated Connection Establishment) cryptographic protocol.

### Domain Models

- **Password:** User PIN or CAN (sensitive, cleared after use)
- **Nonce:** Card-provided random value for mapping
- **MappedDomain:** EC domain parameters after nonce mapping
- **Point:** Elliptic curve point (X, Y coordinates)
- **KeyAgreement:** Ephemeral ECDH state and shared secret
- **AuthenticationTag:** Proof of session key possession
- **MappingType:** Enum (GM = Generic, IM = Integrated)

### PACE Protocol Phases (To Implement)

1. **Password Processing:** PIN → K_pi (one-way via SHA-256)
2. **Mapping Phase:** Nonce decryption + domain parameter mapping
3. **Key Agreement:** Ephemeral ECDH on mapped curve → shared secret Z
4. **Authentication:** Mutual proof using authentication tags

### Application Service

- **PACEService:** Orchestrates all PACE phases
  - Method: `Execute()` → returns established session or error

### Dependencies

- Internal: smartcard, crypto contexts
- External: None (crypto delegated to crypto context)

### Notes

- K_pi never leaves this context
- Nonce decryption uses K_pi + AES-128-CBC
- Both terminal and card compute authentication tags independently

---

## 3. Crypto Context (`internal/crypto/`)

### Purpose

Provide cryptographic primitives (isolated from business logic).

### Domain Models (Interfaces)

- **EllipticCurve:** Abstractshape scalar multiplication, point addition
- **AESKey:** Symmetric key wrapper
- **KDF:** Key derivation function interface

### Infrastructure

- **Brainpool:** BrainpoolP256r1 elliptic curve implementation
  - Status: TODO - implement ECC point operations
  - Reference: BSI TR-03110
- **CMAC:** AES-CMAC wrapper (uses `golang.org/x/crypto/cmac`)
- **ECDH:** Elliptic Curve Diffie-Hellman operations

### Key Derivation

- TR-03110 KDF using counter mode (AES or SHA-256)
- Derives K_enc (encryption) and K_mac (MAC) from shared secret Z

### Dependencies

- External: `golang.org/x/crypto`
- Internal: None

### Notes

- This context is infrastructure-heavy; domain logic is minimal
- Brainpool support is critical blocker - evaluate third-party Go libraries early

---

## 4. Messaging Context (`internal/messaging/`)

### Purpose

Encrypt/decrypt APDUs after successful PACE authentication.

### Domain Models

- **SecureMessage:** Encrypted APDU + CMAC (authentication tag)
- **SendSequenceCounter (SSC):** Tracks message number (must increment before encryption)

### Domain Rules

- SSC incremented before each encryption
- CMAC verified before decryption (prevents tampering)
- Encryption is bidirectional (both commands and responses)

### Application Service

- **SecureMessagingService:** Wraps/unwraps APDUs
  - Methods: `Encrypt(apdu)` → SecureMessage
  - Methods: `Decrypt(response)` → plaintext APDU

### Dependencies

- Internal: crypto context (K_enc, K_mac, CMAC)
- Requires: Successful PACE authentication first

### Notes

- Only active after PACE mutual authentication succeeds
- SSC state is mutable per session; never reset

---

## 5. Card Data Context (`internal/carddata/`)

### Purpose

Read and parse protected personal identity data from CEI.

### Domain Models

- **Identity:** Personal data record (surname, CNP, DOB, citizenship, etc.)
- **IdentityRepository:** Interface for accessing identity (read/save)

### Domain Rules

- Only readable after successful PACE authentication
- Data is immutable once read
- Parsing respects TLV/ASN.1 structure from card file system

### Application Service

- **CardDataService:** Orchestrates file navigation and parsing
  - Methods: `ReadIdentity()` → Identity record

### Dependencies

- Internal: smartcard context (SELECT, READ BINARY), messaging context (encryption)
- External: TLV/ASN.1 parsing library (optional)

### Notes

- Must use Secure Messaging for all file access
- CEI file structure: MF → DF → EF (specific paths for identity data)
- Last context to activate; depends on everything else

---

## Context Interaction Flow

```text
User Input
    ↓
SmartCard Context (PIN verification)
    ↓
PACE Context (establish authenticated channel)
    ↓
Messaging Context (encryption setup)
    ↓
CardData Context (read protected data)
    ↓
Output Identity
```

### Communication Patterns

- **No direct imports between contexts** - use interfaces only
- **Dependency Injection:** Each context receives its dependencies (e.g., PACEService gets SmartCardService)
- **Errors bubble up:** Domain errors are propagated to application layer

---

## Testing Strategy

1. **SmartCard Context:** Mock PC/SC, test APDU serialization/parsing
2. **PACE Context:** Use TR-03110 test vectors, mock nonce/curve operations
3. **Crypto Context:** Validate against NIST/BSI reference vectors
4. **Messaging Context:** Test encryption/decryption with known keys
5. **CardData Context:** Mock card responses, test TLV parsing

---

## Migration Checklist

- [ ] All domain models defined with clear invariants
- [ ] Application services orchestrate logic correctly
- [ ] Infrastructure isolated (scard, crypto libraries)
- [ ] Interfaces defined for cross-context communication
- [ ] No circular dependencies between contexts
- [ ] Error types use domain-specific status codes
- [ ] DI container (if needed) wires contexts correctly
