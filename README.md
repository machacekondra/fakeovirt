[![Fakeovirt Repository on Quay](https://quay.io/repository/kubev2v/fakeovirt/status "Fakeovirt Repository on Quay")](https://quay.io/repository/kubev2v/fakeovirt)

This repository is used for testing imports from oVirt using Forklift.

# Deployment
Deployment example on k8s:

```yaml
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: fakeovirt
  namespace: default
spec:
  replicas: 1
  selector:
    matchLabels:
      app: fakeovirt
  template:
    metadata:
      labels:
        app: fakeovirt
    spec:
      containers:
      - name: fakeovirt
        image: machacekondra/fakeovirt
        ports:
        - containerPort: 12346
        env:
          - name: NAMESPACE
            value: default
          - name: PORT
            value: 12346
---
apiVersion: v1
kind: Service
metadata:
  name: fakeovirt
  namespace: default
spec:
  selector:
    app: fakeovirt
  type: NodePort
  ports:
  - name: fakeovirt
    port: 12346
      targetPort: 12346
```

# Stubbing 

## Recording
To stub response to a request for a given URL, send POST request to the `/stub` endpoint of the service providing desired behavior in its body.
```
curl -X POST https://fakeovirt:12346/stub -d @stubbing.json
```

where `stubbing.json` contains:
```json
[
  {
    "path": "/ovirt-engine/api/vms/123",
    "method": "GET",
    "responses": [
          {
            "responseBody": "<vm href=\"/ovirt-engine/api/vms/123\" id=\"123\">\n<name>cirrosvm</name>\n<description/>\n<comment/>\n<link href=\"/ovirt-engine/api/vms/123/graphicsconsoles\" rel=\"graphicsconsoles\"/>\n<link href=\"/ovirt-engine/api/vms/123/diskattachments\" rel=\"diskattachments\"/>\n<link href=\"/ovirt-engine/api/vms/123/nics\" rel=\"nics\"/>\n<bios>\n<boot_menu>\n<enabled>false</enabled>\n</boot_menu>\n<type>q35_sea_bios</type>\n</bios>\n<cpu>\n<architecture>x86_64</architecture>\n<topology>\n<cores>1</cores>\n<sockets>1</sockets>\n<threads>1</threads>\n</topology>\n</cpu>\n<cpu_shares>0</cpu_shares>\n<creation_time>2020-03-06T09:46:41.294+01:00</creation_time>\n<delete_protected>false</delete_protected>\n<display>\n<allow_override>false</allow_override>\n<copy_paste_enabled>true</copy_paste_enabled>\n<disconnect_action>LOCK_SCREEN</disconnect_action>\n<file_transfer_enabled>true</file_transfer_enabled>\n<monitors>1</monitors>\n<single_qxl_pci>false</single_qxl_pci>\n<smartcard_enabled>false</smartcard_enabled>\n<type>spice</type>\n</display>\n<high_availability>\n<enabled>false</enabled>\n<priority>1</priority>\n</high_availability>\n<io>\n<threads>0</threads>\n</io>\n<memory>134217728</memory>\n<memory_policy>\n<ballooning>false</ballooning>\n<guaranteed>134217728</guaranteed>\n<max>536870912</max>\n</memory_policy>\n<migration>\n<auto_converge>inherit</auto_converge>\n<compressed>inherit</compressed>\n<encrypted>inherit</encrypted>\n</migration>\n<migration_downtime>-1</migration_downtime>\n<multi_queues_enabled>true</multi_queues_enabled>\n<origin>ovirt</origin>\n<os>\n<boot>\n<devices>\n<device>hd</device>\n</devices>\n</boot>\n<type>other</type>\n</os>\n<placement_policy>\n<affinity>pinned</affinity>\n<hosts>\n</hosts>\n</placement_policy>\n<serial_number>\n<policy>vm</policy>\n</serial_number>\n<sso>\n<methods/>\n</sso>\n<start_paused>false</start_paused>\n<stateless>false</stateless>\n<storage_error_resume_behaviour>auto_resume</storage_error_resume_behaviour>\n<time_zone>\n<name>Etc/GMT</name>\n</time_zone>\n<type>desktop</type>\n<usb>\n<enabled>false</enabled>\n</usb>\n<next_run_configuration_exists>false</next_run_configuration_exists>\n<numa_tune_mode>interleave</numa_tune_mode>\n<status>down</status>\n<stop_reason/>\n</vm>\n",
            "responseCode": 200,
            "times": 1
          }
        ]
  },
  {
    "path": "/ovirt-engine/api/vms/xyz",
    "method": "GET",
    "responses": [
         {
           "responseCode": 404
         }
       ]
  }
]
```

## Resetting 

## Defaults
To reset fakeovirt behavior to default one, send POST request to the `/reset` endpoint

```
curl -X POST https://fakeovirt:12346/reset
```

## Only some static paths
To enable only chosen statically-defined stubbings, send  POST request to the `reset` endpoint with `configurators` parameter.
The parameter can have multiple values from available below:
- `static-vms` - enables pre-defined VMs and related resources   
- `static-sso` - enables pre-defined SSO endpoint
- `static-namespace` - enables pre-defined namespace endpoint
- `static-transfers` - enables pre-defined image transfers endpoint

For example:
```
curl -X POST https://fakeovirt:12346/reset?configurators=static-sso,static-transfers
```
