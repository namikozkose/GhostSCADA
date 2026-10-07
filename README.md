# GhostSCADA 👻⚡

**Advanced ICS/SCADA Man-in-the-Middle (MitM) Data Poisoning Proxy**

GhostSCADA is a proof-of-concept (PoC) MitM proxy written in Go, designed to simulate the "Screen Blinding" techniques used by advanced persistent threats (APTs) like Stuxnet. 

This module sits between an edge industrial sensor (ModbusTrap) and a central SOC server (SCADA-C2). When the edge sensor detects a physical sabotage attempt and sends a critical alert, GhostSCADA intercepts the HTTP/JSON payload in transit, modifies the telemetry to indicate "System Stable", and forwards the poisoned data to the SOC. The SOC operators remain completely blind to the ongoing physical destruction.

## ⚙️ Core Architecture & Attack Mechanics

*   **Traffic Interception:** Built using Go's `net/http` to act as a seamless proxy, capturing API traffic destined for the central C2 server.
*   **On-the-fly Data Poisoning:** Parses incoming JSON telemetry from edge nodes. If a `Modbus Function 0x05` (Sabotage) alert is detected, it rewrites the `attack_type` to a benign status and zeroes out the raw hexadecimal footprint.
*   **HMI Screen Blinding:** Simulates a critical flaw in ICS networks where edge-to-center telemetry lacks Mutual TLS (mTLS) or cryptographic payload signing (HMAC), allowing undetected manipulation of the Human-Machine Interface (HMI).

## 🚀 Defense & Mitigation

To prevent this class of attack, the architecture must implement:
1.  **mTLS (Mutual TLS):** Cryptographically authenticating both the edge node and the central server to prevent rogue proxies from terminating the connection.
2.  **Payload Signing (HMAC):** Edge nodes must cryptographically sign the JSON payload. Any modification in transit by a MitM proxy will invalidate the signature at the SOC level.