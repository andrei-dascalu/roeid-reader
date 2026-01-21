# roeid-reader PACE Implementation Roadmap

**Estimated effort:** 4–8 weeks for a minimal, working implementation

**Goal:** Implement full PACE (Password Authenticated Connection Establishment) protocol to securely unlock Romanian CEI cards and read protected identity data.

---

## Phase 1: Foundation & Preparation

### Task 1.1: Knowledge Base Review

- [x] Study ISO/IEC 7816-4 (Smart Card APDUs)
  - ISO/IEC 7816-4:2020 standard (official spec)
  - [OpenSC APDU documentation](https://github.com/OpenSC/OpenSC/wiki/APDU-Commands)
  - [JMRTD eMRTD specification](https://github.com/jmrtd/jmrtd/wiki)

- [x] Review ISO/IEC 7816-8 & 7816-9 (Cryptography & SM concepts)
  - ISO/IEC 7816-8:2016 (Cryptographic Information Security)
  - [GlobalPlatform Secure Messaging specification](https://www.globalplatform.org/)
  - [OpenSSL SM documentation and examples](https://www.openssl.org/)

- [x] Read BSI TR-03110 (PACE specification)
  - [BSI TR-03110 v2.13 PDF](https://www.bsi.bund.de/EN/Publications/TechnicalGuidelines/tr03110/index_htm.html)
  - [BoringSSL PACE implementation reference](https://boringssl.googlesource.com/boringssl/+/master/src/crypto/)
  - [JMRTD PACE source code](https://github.com/jmrtd/jmrtd/tree/master/jmrtd/src/main/java/net/sf/scuba/smartcards/)

- [x] Understand NIST SP 800-38B (AES-CMAC)
  - [NIST SP 800-38B PDF](https://nvlpubs.nist.gov/nistpubs/Legacy/SP/nistspecialpublication800-38b.pdf)
  - [golang.org/x/crypto/cmac package documentation](https://pkg.go.dev/golang.org/x/crypto/cmac)
  - [AES-CMAC test vectors and examples](https://csrc.nist.gov/projects/cryptographic-algorithm-validation-program/)

- [x] Document key concepts in project wiki/comments
  - Create `docs/PACE_GLOSSARY.md` with terms (nonce, KDF, mapping, SSC, etc.)
  - Add inline comments in code with RFC/standard references
  - Create `docs/PROTOCOL_FLOW.md` with sequence diagrams

**Deliverable:** Annotated references and local notes

---

## Phase 2: Raw APDU Layer (Step 1)

### Task 2.1: PC/SC Transport Enhancement

- [ ] Review existing `src/infrastructure/reader.go` code
- [ ] Create `src/apdu/transport.go` for low-level APDU handling
- [ ] Implement full APDU logging (all requests/responses with timestamps)
- [ ] Add error handling for PC/SC failures
- [ ] Test with current PIN verification flow

**Deliverable:** Clean APDU logging; all card interactions visible in logs

### Task 2.2: APDU Response Parsing

- [ ] Create `src/apdu/response.go` to parse status words
- [ ] Document common status word mappings (0x9000, 0x6982, etc.)
- [ ] Add validation helpers (e.g., `IsSuccess()`, `IsSecurityError()`)

**Deliverable:** Reusable response parsing utilities

---

## Phase 3: Application Selection (Step 2)

### Task 3.1: CEI Application Discovery

- [ ] Implement SELECT by AID (`D2 76 00 01 24 01`)
- [ ] Parse FCI (File Control Information) response
- [ ] Document expected FCI structure for CEI cards

**Deliverable:** Ability to SELECT and identify CEI card

### Task 3.2: PACE Algorithm Capability Discovery

- [ ] Create `src/pace/info.go` with `PaceInfo` struct
- [ ] Parse PACE OIDs from FCI
- [ ] Identify supported curves (BrainpoolP256r1, etc.)
- [ ] Identify supported ciphers (AES-128, AES-256)

**Deliverable:** Runtime detection of card's PACE capabilities

---

## Phase 4: Password Processing (Step 3)

### Task 4.1: PIN/CAN to Key Material Derivation

- [ ] Create `src/pace/password.go`
- [ ] Implement K_pi derivation from PIN using SHA-256
- [ ] Support both PIN and CAN password types
- [ ] Add test vectors from TR-03110

**Deliverable:** K_pi derivation tested against reference vectors

**Note:** K_pi is never sent to the card; stays local

---

## Phase 5: Brainpool Elliptic Curve Support

### Task 5.1: Brainpool Curve Implementation

- [ ] Evaluate options:
  - Port existing Go ECC libraries
  - Use `golang.org/x/crypto` (if sufficient)
  - Implement custom Brainpool256r1
- [ ] Create `src/crypto/brainpool.go`
- [ ] Implement point addition, scalar multiplication
- [ ] Test with known test vectors

**Deliverable:** Working BrainpoolP256r1 ECDH

### Task 5.2: Elliptic Curve Utilities

- [ ] Point encoding/decoding (compressed/uncompressed)
- [ ] Scalar generation (random, safe)
- [ ] Public key validation

**Deliverable:** Reusable EC cryptography module

---

## Phase 6: Mapping Phase (Step 4) — **Hardest Part**

### Task 6.1: Encrypted Nonce Exchange

- [ ] Implement GENERAL AUTHENTICATE command to request mapping nonce
- [ ] Decrypt nonce using K_pi (AES-128 in CBC mode)
- [ ] Parse nonce structure

**Deliverable:** Successfully decrypt card's mapping nonce

### Task 6.2: Nonce-to-Domain Parameters

- [ ] Implement mapping (GM = Generic Mapping or IM = Integrated Mapping)
- [ ] Map decrypted nonce to EC domain parameters
- [ ] Validate mapped curve is identical to original

**Deliverable:** Verified domain parameter mapping

### Task 6.3: Mapped Public Key Generation & Exchange

- [ ] Generate ephemeral key pair on mapped domain
- [ ] Encode mapped public key
- [ ] Send to card via GENERAL AUTHENTICATE
- [ ] Receive and validate card's mapped public key

**Deliverable:** Successful bidirectional mapped key exchange

---

## Phase 7: Key Agreement (Step 5)

### Task 7.1: ECDH on Mapped Curves

- [ ] Implement ECDH shared secret computation
- [ ] Compute Z = Kp_mapped × Qc_mapped
- [ ] Validate Z is not point-at-infinity

**Deliverable:** Shared secret Z established with card

### Task 7.2: Authentication Tags Exchange

- [ ] Receive `T_c` (card's authentication tag) from GENERAL AUTHENTICATE
- [ ] Compute local `T_pi` (terminal's authentication tag)
- [ ] Compare hashes to detect MITM or wrong password

**Deliverable:** Mutual authentication before session key derivation

---

## Phase 8: Session Key Derivation (Step 6)

### Task 8.1: Key Derivation Function (KDF)

- [ ] Implement TR-03110 KDF (counter mode with AES or SHA-256)
- [ ] Derive K_enc (encryption key) and K_mac (MAC key)
- [ ] Store SSC (Send Sequence Counter) = 0

**Deliverable:** Session keys K_enc and K_mac

### Task 8.2: Test with Known Vectors

- [ ] Validate KDF against TR-03110 test vectors
- [ ] Document key material for debugging

**Deliverable:** Tested KDF matching specification

---

## Phase 9: Mutual Authentication (Step 7)

### Task 9.1: Verify Card Authentication

- [ ] Receive and parse card's final authentication response
- [ ] Compute expected card tag
- [ ] Reject if mismatch (wrong PIN/CAN/crypto error)

**Deliverable:** Confirm card knows session keys

### Task 9.2: Send Terminal Authentication

- [ ] Compute and send terminal's authentication token
- [ ] Document authentication exchange format

**Deliverable:** Bidirectional proof of key possession

---

## Phase 10: Secure Messaging Layer (Step 8)

### Task 10.1: APDU Encryption

- [ ] Create `src/messaging/secure.go`
- [ ] Encrypt APDU command body using K_enc (AES-128-CBC)
- [ ] Increment SSC before each command

**Deliverable:** Encrypted APDU transmission

### Task 10.2: Message Authentication Code (CMAC)

- [ ] Implement AES-CMAC using `golang.org/x/crypto/cmac`
- [ ] Compute CMAC over SSC + encrypted APDU
- [ ] Append CMAC to ciphertext

**Deliverable:** Integrity-protected encrypted APDUs

### Task 10.3: Response Decryption

- [ ] Receive encrypted response
- [ ] Verify CMAC matches
- [ ] Decrypt response body
- [ ] Increment SSC

**Deliverable:** End-to-end Secure Messaging channel

---

## Phase 11: Reading Identity Data (Step 9)

### Task 11.1: File Navigation

- [ ] SELECT DFs/EFs for identity data
- [ ] Use Secure Messaging for all SELECT commands
- [ ] Document CEI file structure (MF, DFs, EFs)

**Deliverable:** Navigate card's file system securely

### Task 11.2: Data Extraction

- [ ] READ BINARY via Secure Messaging
- [ ] Parse TLV/ASN.1 structures
- [ ] Extract and display:
  - Surname, Given Names
  - CNP (Cod Numeric Personal)
  - Date of Birth
  - Citizenship
  - ID number expiry

**Deliverable:** Display CEI personal data

---

## Phase 12: Testing & Validation

### Task 12.1: APDU Logging Comparison

- [ ] Capture full APDU logs from successful run
- [ ] Compare with official middleware (or JMRTD)
- [ ] Verify protocol sequence matches spec

**Deliverable:** APDU flow validation against reference

### Task 12.2: Error Case Testing

- [ ] Wrong PIN → verify correct error
- [ ] Card removal during PACE → handle gracefully
- [ ] Corrupted nonce → verify rejection

**Deliverable:** Robust error handling

### Task 12.3: Integration Tests

- [ ] End-to-end test with real CEI card
- [ ] Verify all data reads are correct
- [ ] Document known limitations

**Deliverable:** Passing integration tests

---

## Phase 13: Documentation & Cleanup

### Task 13.1: Code Documentation

- [ ] Document PACE phases with references to TR-03110
- [ ] Add comments explaining each cryptographic step
- [ ] Create package-level documentation

**Deliverable:** Well-documented codebase

### Task 13.2: Project Structure

- [ ] Organize code by module (APDU, PACE, Crypto, Messaging)
- [ ] Add example usage
- [ ] Update README with PACE support

**Deliverable:** Production-ready code structure

---

## Decision Point: Hybrid vs. Full Native

**Before Phase 5:**

- **Option A (Recommended for Speed):** Stop after Phase 3, delegate PACE to JMRTD (Java), integrate results
- **Option B (Full Control):** Continue through Phase 13, native Go implementation

**Current choice:** [Choose A or B]

---

## Quick Reference: Key Files to Create

| Module | File | Purpose |
| --- | --- | --- |
| APDU | `src/apdu/transport.go` | APDU logging & transmission |
| APDU | `src/apdu/response.go` | Status word parsing |
| PACE | `src/pace/info.go` | Algorithm capability discovery |
| PACE | `src/pace/password.go` | K_pi derivation |
| PACE | `src/pace/mapping.go` | Mapping phase (Step 4) |
| PACE | `src/pace/agreement.go` | Key agreement (Step 5) |
| PACE | `src/pace/mutual_auth.go` | Authentication (Step 7) |
| Crypto | `src/crypto/brainpool.go` | Brainpool256r1 ECC |
| Crypto | `src/crypto/kdf.go` | Session key derivation |
| Messaging | `src/messaging/secure.go` | Encryption + CMAC |
| Main | `main.go` | Update to use PACE layer |

---

## Troubleshooting Checklist

- [ ] pcscd running? (`pcscd -h` or `systemctl status pcscd`)
- [ ] Card inserted and recognized? (`pcsc_scan`)
- [ ] APDU logs showing all exchanges?
- [ ] TR-03110 test vectors match implementation?
- [ ] Wrong PIN returns expected error code?

---

## Notes

- **Brainpool support is the biggest blocker.** Evaluate third-party libraries early.
- **APDU logging is critical.** If stuck, compare against official middleware.
- **Test each phase independently** before moving to the next.
- **Use reference implementations** (JMRTD, official CEI middleware) to validate protocol correctness.
