# Domain Models & Relationships

## Entity & Value Object Diagrams

### Smart Card Domain

```text
Card (Aggregate Root)
  ├── APDU (Value Object) → Transmit
  ├── Response (Value Object) ← Receive
  │   ├── Data: []byte
  │   ├── SW1, SW2: byte → StatusCode()
  │   └── IsSuccess(): bool
  ├── CardStatus (Value Object)
  │   ├── ATR: []byte
  │   ├── ActiveProtocol: string
  │   └── Reader: string
  └── StatusError (Value Object)
      ├── Code: uint16
      ├── SW1, SW2: byte
      └── String(): string
```text

**APDU Serialization:** `[CLA] [INS] [P1] [P2] [Lc] [Data...] [Le]`

**Status Words:** Last 2 bytes of response

- `0x9000` = Success
- `0x6982` = Security status not satisfied (wrong PIN)
- `0x6984` = Reference key in use

---

### PACE Domain

```text
PACE Protocol Flow
  │
  ├─→ Phase 1: Password Processing
  │   ├── Input: Password (PIN)
  │   ├── Process: SHA-256 derivation
  │   └── Output: K_pi (session password material)
  │       └── Never transmitted; stays local
  │
  ├─→ Phase 2: Mapping Phase
  │   ├── Card sends: Encrypted Nonce
  │   ├── Terminal: AES-CBC decrypt(nonce, K_pi)
  │   ├── Terminal & Card: Map nonce → MappedDomain
  │   │   ├── Input: Nonce bytes
  │   │   ├── Algorithm: GM (Generic) or IM (Integrated)
  │   │   └── Output: Modified EC parameters (P, A, B, G, N, H)
  │   └── Exchange: Mapped public keys
  │
  ├─→ Phase 3: Key Agreement (ECDH)
  │   ├── Terminal: Generate ephemeral key Ks
  │   ├── Both: Exchange public keys (Qs, Qc)
  │   ├── Both: Compute shared secret Z = Ks × Qc
  │   └── Result: Identical shared secret on both sides
  │
  └─→ Phase 4: Mutual Authentication
      ├── Both compute: T_pi (terminal), T_c (card)
      ├── Comparison: If T_pi == T_c → authenticated
      └── Error: Mismatch → wrong password or attack
```text

**Key Agreement State:**

```text
KeyAgreement
  ├── EphemeralPrivateKey: Ks
  ├── EphemeralPublicKey: Qs (our public key)
  ├── CardPublicKey: Qc (their public key)
  └── SharedSecret: Z (both sides identical)
```text

---

### Cryptography Domain

```text
EllipticCurve (Interface)
  ├── Name(): "BrainpoolP256r1"
  ├── P(): Field prime
  ├── A(), B(): Curve coefficients
  ├── G(): Generator point
  ├── Order(): Point order
  ├── ScalarMult(k, p): k × p on curve
  └── Add(p1, p2): p1 + p2 on curve

Point (Value Object)
  ├── X: Big integer
  └── Y: Big integer
     └── IsPointAtInfinity(): bool

AESKey (Value Object)
  ├── key: []byte (secured)
  └── Bytes(): []byte

KDF (Interface)
  └── Derive(data, counter, length): []byte
      └── Counter mode: Iterate hash/AES for multi-block keys
```text

**Supported Curves:**

- BrainpoolP256r1 (primary for CEI)
  - Field: 256-bit prime
  - Generator order: 256-bit
  - Used for PACE mapping & key agreement

**Key Derivation (TR-03110):**

```text
K_enc, K_mac = KDF(Z, counter_enc, counter_mac, KDF_length)
├── Input: Shared secret Z
├── Counter mode: i = 1, 2, ... for multi-block keys
└── Output: K_enc (encryption), K_mac (authentication)
```text

---

### Messaging Domain

```text
SendSequenceCounter (State Object)
  ├── value: uint64 (initialized to 0)
  ├── Increment(): void
  └── Bytes(): [8]byte (big-endian)
      └── Used in CMAC computation

SecureMessage (Value Object)
  ├── EncryptedData: []byte (AES-128-CBC ciphertext)
  └── CMAC: []byte (authentication tag)

Secure APDU Format:
  [Encrypted APDU Body] + [CMAC(SSC + Encrypted)]
  │                       │
  │                       └─→ Ensures integrity & SSC ordering
  └──────────────────────────→ Confidentiality
```text

**Message Protection:**

```text
For each outgoing APDU:
  1. Increment SSC
  2. Encrypt APDU data with K_enc (AES-CBC)
  3. Compute CMAC(SSC || Encrypted) with K_mac
  4. Send [Encrypted || CMAC]

For incoming response:
  1. Verify CMAC(SSC || Encrypted) matches
  2. Decrypt with K_enc
  3. Increment SSC
  4. Return plaintext response
```text

---

### Card Data Domain

```text
Identity (Aggregate Root)
  ├── Surname: string
  ├── GivenNames: string
  ├── CNP: string (Cod Numeric Personal - national ID)
  ├── DateOfBirth: string
  ├── PlaceOfBirth: string
  ├── Citizenship: string
  ├── Gender: string
  ├── DocumentNumber: string
  ├── Series: string
  ├── IssueDate: string
  ├── IssuePlace: string
  └── ExpiryDate: string

IdentityRepository (Interface)
  ├── Read(): (*Identity, error)
  │   └── Via secure messaging (PACE required)
  └── Save(identity): error
      └── Optional local persistence
```text

**Card File Structure (CEI):**

```text
MF (Master File)
  ├── DF.CERT (Certificates)
  ├── DF.AUTH (Authentication)
  └── DF.ID (Identity Data) ← READ via Secure Messaging
      ├── EF.PersonalData (CNP, surname, names, DOB, etc.)
      ├── EF.DocumentData (document number, series, issue/expiry dates)
      └── EF.Address (residence address)

Parsing: TLV-encoded fields (Tag-Length-Value structure)
```text

---

## Aggregates & Boundaries

**SmartCard Aggregate:**

- Root: Card (interface)
- Members: APDU, Response, CardStatus
- Invariant: Status words always present (SW1, SW2)

**PACE Aggregate:**

- Root: KeyAgreement
- Members: Password, Nonce, MappedDomain, AuthenticationTag
- Invariant: Shared secret identical on both sides post-agreement

**Crypto Aggregate:**

- Root: EllipticCurve (stateless interface)
- Members: Point, AESKey, KDF
- Invariant: EC operations must preserve curve properties

**Messaging Aggregate:**

- Root: SendSequenceCounter (mutable state)
- Members: SecureMessage
- Invariant: SSC must increment before each encryption

**CardData Aggregate:**

- Root: Identity (immutable)
- Members: (only primitive fields)
- Invariant: All fields populated post-read; no partial records

---

## Event Sequence (Happy Path)

```text
1. SmartCard Context:
   Terminal connects → Card sends ATR
   Terminal: SELECT application (CEI AID)
   Terminal: VERIFY PIN → Card confirms

2. PACE Context:
   Terminal: Request PACE (GENERAL AUTHENTICATE)
   Card: Return encrypted nonce
   Terminal: Decrypt nonce with K_pi
   Both: Map nonce to EC curve
   Terminal: Send ephemeral public key
   Card: Send ephemeral public key
   Both: Compute shared secret Z
   Both: Exchange authentication tags (mutual auth)

3. Messaging Context:
   Derive K_enc, K_mac from Z
   Initialize SSC = 0

4. CardData Context:
   Terminal: SELECT identity EF (via secure messaging)
   Terminal: READ BINARY (via secure messaging)
   Terminal: Parse TLV → populate Identity record
   Return Identity to application

5. Output:
   Display: Name, CNP, Date of Birth, Citizenship, etc.
```text

---

## Error Cases & Invariant Violations

| Error | Context | Cause | Recovery |
| --- | --- | --- | --- |
| Wrong PIN | PACE | Auth tags don't match | Reject, allow retry |
| Nonce mapping failed | PACE | Corrupted nonce or invalid mapping | Reject & restart |
| CMAC verification fails | Messaging | Corrupted/tampered message | Disconnect, restart |
| Card removal | SmartCard | Physical removal during operation | Graceful error, reconnect |
| Invalid ATR | SmartCard | Card not recognized or malfunction | Log & reject |

---

## Testing Invariants

- **APDU:** Serialize/deserialize round-trip preserves all bytes
- **Response:** Status codes always in last 2 bytes
- **PACE:** K_pi never serialized; Z computed identically on both sides
- **Crypto:** EC point addition is associative & closed
- **Messaging:** SSC increments monotonically; CMAC prevents forgery
- **CardData:** All fields populated after Read(); no nulls in Identity

---

## References

- **ISO/IEC 7816-4:** Smart card APDU specification
- **BSI TR-03110:** PACE protocol specification
- **NIST SP 800-38B:** AES-CMAC specification
- **BrainpoolECC:** European EC standard
