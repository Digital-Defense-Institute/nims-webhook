# NIMS (Notion Incident Management System) Webhook

This is an all-in-one binary that will catch detections sent via webhook to `/hooks/alert:9000` and create alerts in your [NIMS alerts database](https://ecapuano.notion.site/Alpha-Notion-Incident-Management-System-16ad5339b31680679a91e41ac722dfcc).

## To build 
Install Go
```bash
sudo apt update
sudo apt -y install golang-go
```
Build the binary
```bash
cd nims-webhook
go mod init nims-webhook
go mod tidy
go build nims-webhook.go
```

## To run
Either build the binary (steps above) if you wish to make modifications, or download it from [the releases page](https://github.com/shortstack/nims-webhook/releases) and run it.  

First, replace Notion auth token and database IDs with yours in `.env`
```bash
NIMS_ASSETS_DATABASE_ID=
NIMS_ALERTS_DATABASE_ID=
NOTION_AUTH_TOKEN=
```
Run the binary
```bash
chmod +x nims-webhook
./nims-webhook
```

## Fields
The following fields are currently utilized:
* `routing.hostname` - the hostname of the affected host
* `routing.int_ip` - the internal IP address of the affected host
* `routing.event_time` - the timestamp of the detection event
* `detect` - the full event details captured during detection
* `detect_mtd` - metadata associated with the detection
* `link` - the URL linking directly to the alert within LimaCharlie
* `cat` - the name or category of the alert  

To customize these fields or replace them with others from your JSON objects, you can edit the `nims-webhook.go` file, specifically in the `webhookHandler` function.

Similarly, if you wish to modify fields in your Notion template and integrate those changes into the script, updates can be applied in both the `webhookHandler` and `addAlert` functions.

## Example request 
```bash
curl -X POST http://0.0.0.0:9000/hooks/alert \
-H "Content-Type: application/json" \
-d '{
  "author": "_soteria-rules-edr-123abc45-678d-901e-fghi-234567jklmno[bulk][segment]",
  "cat": "00456-WIN-mshta_Network_Connection_to_External_IP",
  "detect": {
    "event": {
      "COMMAND_LINE": "C:\\Windows\\System32\\mshta.exe",
      "CREATION_TIME": 1736019135814,
      "FILE_IS_SIGNED": 1,
      "FILE_PATH": "C:\\Windows\\System32\\mshta.exe",
      "HASH": "a1234567b89cd012ef34gh567ijklmn8901234567abcdef890123456789abcdef",
      "NETWORK_ACTIVITY": [
        {
          "DESTINATION": {
            "IP_ADDRESS": "192.168.15.23",
            "PORT": 443
          },
          "IS_OUTGOING": 1,
          "PROTOCOL": "tcp4",
          "SOURCE": {
            "IP_ADDRESS": "172.16.10.45",
            "PORT": 60432
          },
          "STATE": 8,
          "TIMESTAMP": 1736019210633
        }
      ],
      "PARENT_PROCESS_ID": 1052,
      "PROCESS_ID": 2032,
      "USER_NAME": "CORP\\AdminUser"
    },
    "routing": {
      "arch": 2,
      "did": "",
      "event_id": "4a3b2c1d-e5f6-47g8-h9ij-k123lmnopq45",
      "event_time": 1736019224486,
      "event_type": "NETWORK_CONNECTIONS",
      "ext_ip": "203.0.113.45",
      "hostname": "corporate-webserver.corp.internal",
      "iid": "f12g34h5-i6jk-78lm-90no-pq12rstuv345",
      "int_ip": "172.16.10.45",
      "moduleid": 3,
      "oid": "5678abcd-910e-11fg-hijk-123456lmnopq",
      "parent": "abcd1234ef567gh8910ijklm234nopqr",
      "plat": 268435456,
      "sid": "789abcd1-2345-6789-0efg-hijklm123nop",
      "tags": [
        "windows",
        "suspicious_execution"
      ],
      "this": "123abc456def789ghi012jkl345mnop678"
    }
  },
  "detect_id": "abcdef12-3456-789a-0bc1-defghijklmno",
  "detect_mtd": {
    "description": "MSHTA is a legitimate tool used to execute HTML applications. It can be abused by attackers to download and execute malicious scripts. This detector identifies mshta.exe making external network connections, which is indicative of potential malicious activity.",
    "falsepositives": [
      "Legitimate administrative use of mshta.exe in secure environments."
    ],
    "references": [
      "https://lolbas-project.github.io/lolbas/Binaries/Mshta/",
      "https://attack.mitre.org/techniques/T1218/005/",
      "https://redcanary.com/threat-detection-report/techniques/mshta/"
    ],
    "tags": [
      "attack.t1218.005",
      "attack.t1071.001",
      "attack.t1105"
    ]
  },
  "gen_time": 1736019224489,
  "link": "https://app.limacharlie.io/orgs/5678abcd-910e-11fg-hijk-123456lmnopq/sensors/789abcd1-2345-6789-0efg-hijklm123nop/timeline?time=1736019224&selected=123abc456def789ghi012jkl345mnop678",
  "namespace": "general",
  "priority": 2,
  "routing": {
    "arch": 2,
    "did": "",
    "event_id": "4a3b2c1d-e5f6-47g8-h9ij-k123lmnopq45",
    "event_time": 1736019224486,
    "event_type": "NETWORK_CONNECTIONS",
    "ext_ip": "203.0.113.45",
    "hostname": "corporate-webserver.corp.internal",
    "iid": "f12g34h5-i6jk-78lm-90no-pq12rstuv345",
    "int_ip": "172.16.10.45",
    "moduleid": 3,
    "oid": "5678abcd-910e-11fg-hijk-123456lmnopq",
    "parent": "abcd1234ef567gh8910ijklm234nopqr",
    "plat": 268435456,
    "sid": "789abcd1-2345-6789-0efg-hijklm123nop",
    "tags": [
      "windows",
      "suspicious_execution"
    ],
    "this": "123abc456def789ghi012jkl345mnop678"
  },
  "source": "5678abcd-910e-11fg-hijk-123456lmnopq.f12g34h5-i6jk-78lm-90no-pq12rstuv345.789abcd1-2345-6789-0efg-hijklm123nop.10000000.3",
  "source_rule": "service.WIN-mshta_Network_Connection_to_External_IP",
  "ts": 1736019224000
}'
```
