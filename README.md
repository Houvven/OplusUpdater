# Oplus Updater

Use Oplus official api to query OPlus,OPPO and Realme Mobile 's OS version update.

## Install

```shell
$ go install github.com/Houvven/OplusUpdater/cmd/updater@latest
```

## How to use?
```shell
$ updater -h                                              
Use Oplus official api to query OPlus,OPPO and Realme Mobile 's OS version update.

Usage:
  updater [flags]

Flags:
      --carrier my_manifest/build.prop   Found in my_manifest/build.prop file, under the `NV_ID` reference, e.g., --carrier=01000100
  -h, --help                             help for oplus-updater
      --imei string                      IMEI, e.g., --imei=86429XXXXXXXX98
      --mode int                         Mode: 0 (stable, default) or 1 (testing), e.g., --mode=0
      --model string                     Device model, e.g., --model=RMX3820
  -o, --ota-version string               OTA version (required), e.g., --ota-version=RMX3820_11.A.00_0000_000000000000
  -p, --proxy string                     Proxy server, e.g., --proxy=type://@host:port or --proxy=type://user:password@host:port, support http and socks proxy
      --region string                    Server zone: CN (default), EU or IN (optional), e.g., --region=CN (default "CN")
```

## Update request headers

| header         | getter                                              | example                           | Required |
|----------------|-----------------------------------------------------|-----------------------------------|----------|
| language       | `Local.getDefault().getLanguageTag()`               | zh-CN                             | ✅        |
| newLanguage    | `Local.getDefault().getLanguageTag()`               | zh-CN                             | ⭕        |
| romVersion     | _ro.build.display.id_                               | RMX3350_13.1.0.400(CN01)          | ⭕        |
| androidVersion | `"Android" + Build.VERSION.RELEASE`                 | Android 13                        | ✅        |
| colorOSVersion | "ColorOS" + _ro.build.version.oplusrom_             | ColorOS13.1.0                     | ✅        |
| infVersion     | magic value: 1                                      | 1                                 | ⭕        |
| otaVersion     | _ro.build.version.ota_                              | RMX3350_11.F.25_3250_202403011232 | ✅        |
| model          | `Build.MODEL `                                      | rmx3350                           | ✅        |
| mode           | `client_auto` or `server_auto` or `manual` or `???` | manual                            | ⭕        |
| nvCarrier      | _ro.build.oplus_nv_id_                              |                                   | ✅        |
| pipelineKey    | _ro.oplus.pipeline_key_                             |                                   | ⭕        |
| companyId      | _ro.oplus.company_id_                               | ~~is empty?~~                     | ⭕        |
| prjNum         | _ro.separate.soft_                                  |                                   | ⭕        |
| brand          | `Build.BRAND`                                       | realme                            | ⭕        |
| brandSota      | _ro.product.brand_                                  | realme                            | ⭕        |
| osType         | _ro.oplus.image.my_stock.type_                      | domestic_realme                   | ⭕        |
|                |                                                     |                                   |          |
| deviceId       | OPPO Framework provide GUID and SHA                 |                                   | ✅        |
| protectedKey   |                                                     |                                   | ✅        |

### Simple

```http
# Headers
language: zh-CN
newLanguage: zh-CN
romVersion: RMX3350_13.1.0.400(CN01)
androidVersion: Android13
colorOSVersion: ColorOS13.1.0
infVersion: 1
otaVersion: RMX3350_11.F.00_0000_000000000000
model: RMX3350
mode: manual
nvCarrier: 10010111
pipelineKey: ALLNET
companyId: 
prjNum: 21609
brand: realme
brandSota: realme
osType: domestic_realme
version: 2
deviceId: 628A6A292C3274921C497UW27DH201ODU17484CA5630345B24D20DAE1C792B90
protectedKey: {"SCENE_1":{"protectedKey":"VU0P4dUKr9gPQdiDdLI1QpV2e2K39D\/rGge4BCfQ5ZgreUJr\/QyvUx5wZrZKBff3NvMfeU9f3TIiPBxESELLW+\/GdSjuuHYr2yg??????????\/xyAH8K74HKj+OWSh8eRhi\/U3pti8T1sqPVLAXlvnrh\/3QBzaYCXuaXn823TNE+5Scnl0VaXjWwf2siNsKKAZUpeef3xCwn\/u8ILhYfTOCzOVNX+ZVXF8OW+RbhFZX9cff4y6RG943gHZyYI+H67UWnY2TjW8VP1\/FPHdx4bFLMRphE6psXYXY\/HAWLRqTZVilT\/BHWYM7HpD26lTmbb4oyfzcEy+vVo+YsGQXCZg==","version":"1717246406204","negotiationVersion":"1615879139745"}}
Content-Type: application/json; charset=utf-8
Content-Length: 1928
Host: component-ota-cn.allawntech.com
Connection: Keep-Alive
Accept-Encoding: gzip
User-Agent: okhttp/3.12.2
```

### protectedKey

type: `CryptoConfig`

- version: `timestamp`
- negotiationVersion: `pref`, name: app_feature

```json
{
  "SCENE_1": {
    "protectedKey": "VU0P4dUKr9gPQdiDdLI1QpV2e2K39D\/rGge4BCfQ5ZgreUJr\/QyvUx5wZrZKBff3NvMfeU9f3TIiPBxESELLW+\/GdSjuuHYr2yg??????????\/xyAH8K74HKj+OWSh8eRhi\/U3pti8T1sqPVLAXlvnrh\/3QBzaYCXuaXn823TNE+5Scnl0VaXjWwf2siNsKKAZUpeef3xCwn\/u8ILhYfTOCzOVNX+ZVXF8OW+RbhFZX9cff4y6RG943gHZyYI+H67UWnY2TjW8VP1\/FPHdx4bFLMRphE6psXYXY\/HAWLRqTZVilT\/BHWYM7HpD26lTmbb4oyfzcEy+vVo+YsGQXCZg==",
    "version": "1717246406204",
    "negotiationVersion": "1615879139745"
  }
}
```

## Update Request body raw data

```json
{
  "params": "{\"cipher\":\"dGnvtjHkz8HQzG0o74PECl1JAfXSdk\\\/gx0ykJqqqGh4rJSO3cpPskGASfCAxwgMiCfB7I9pLKXq2kjzH82ItvQAfvwvzAb9VFKyrYsYTDV5TrdrGgmZCzjCcBba6qkXfAhJrEPXii4mLtl+YZDe47ff04abgEoi4nU33KXxw87xH+opJZTR3aQasUtG5xci+QDelsSYqslAnZ7CSbMTv5oEWfuUKrNMYH8Tjayn5tUWuswRv8H7CjR+anw8\\\/dHrRxyFOpYWNhwcorp7NAq5R99Zn5ZO+0ETDS3tcjLjc66YyiSyq0eXHw5pX4+HdQAUcs9PMo8HYVehh4IPHDpH8npcnKuVbOry2O4gKsUyOfnNNq9ua4MD5eV7qplRRI7kJ9JwbeM9oBOlJr1F4WjSHL7xKiFBFrybUlB+xXdKUO+OA4Fy0OCVsor4vo7r3RH0+Yvu99aiQg92ppCLL7CLbZpl62vwq2+e5aUfARBbTr8HSv5RJPSNY1ryIRE+j3QKViBtufIyYJHllnQYTmKWiO8AQdeL\\\/oIq4nuYXYeekwEoyfPzajvjOo9FdQ\\\/\\\/jMx+CW5JLIhDa6yVfo700txNPgqvXziYpalGYx2qTI2HvHtyjmXwTn7FEYcOjJ7SDofLZR\\\/Z8iTQ\\\/FJNaF20OVDxhJqT\\\/lHikPOsDxHHw1OM3v2W+kAUeRejT56rq6EMQgp9U09y1S8lCYIGWHnQvznB2gny3lvRUxUf5WEa7cvcmx3\\\/PgVo0N5EE9N1kX05vuvPWqc7ugzEDS9r2i4bWA9ttvt6++SJHRDEBejTVpvUC\\\/KhlDDH\\\/XA8yY9MM0XlMl3phDRBpXbzg0v84F\\\/QWJvQ+LjExEN9naRi\\\/Hu\\\/LKcg\\\/QTWUtu8OJHWDOaGnkZ1NB9pwWysvv6aRCrS6RAlEby5avAT+gvndn0NBed9pE1om7Fi9vruSMiT0ebyM5lgpvs2iDFiHIKca8TjdgmjyNYeAX6kDyl9WksZE7q4sCgvqHpPuDJVbtjMHawupcyHhW83MogZWJhTd0U98Y9Kq+rNMT\\\/hZEp\\\/dA+Dd\\\/DU6KVsbWIcQzzUExf2PugW3lRTVPoCqhv\\\/0TQSut7VCUmRBy61Ake1lpVhIzNy8L+SjnBtgEAqr4VASgPVqvf0zU0cNSRZAHUlIhpLtwIT2lJ6InoneVFTLeRzSbFuWkdcjJnX\\\/wlyfMTDpQCYTPqOZHrs6BBC9Ow4hV4nOXDY0E9k+u0hllk14pq4Xz9VoSEQ2UO2aLaJG8uGbGUcUMi5\\\/TA0tAZ2TPbtOmdID4z2XglplAMXFCy8ZcT7JnVUBb85sII0KJMdv\\\/t\\\/j6\\\/bU5TRtMEJ58kOi\\\/LjRd4GY9KOmdVms5iYDvS9AxGYTMKgBHUlSZpIIoB9YdvJtkaf8uZIj0lzj+KAQd1ZFd+0Gc+M1D+G52o0wWV87XiqYgUAlaU1nYqZk7leefqGBq2LwDGW4PJ3VO8bDWGWlMjBpB06f0VdHEd7bLxSklo3XDmRvI0ny+mFJcy7Mdm1aiDzmQG9nu+HaGXFRPg\\\/1+UJF28Q9v88cyq\\\/\\\/NvwLHqNWnTHGjEr2\\\/5E2rGxFnDe8rBrGrfW6qyyDF5qCMhi1Pxrl+69cZZ9+una+oKYKDOfFUDG9YeSryGV1MBqnGvnNI1HTpw5RvgJhGM\\\/EnXNB3vrP8pMO1VVxD5+6FBcBA1nh3ew+X61yYsmavBrlq0nh8BBs5fJhBKexrNdbmOcZ48DIlU46lTNqpKhP\\\/+e5vdM1DVhyot6w+8A3HYRFqCLUzZtWnfJHuXVD92FEdzA=\",\"iv\":\"gUW5DpcsEyQ+xQpMwH53Eg==\"}"
}
```

```json
{
  "params": {
    "cipher": {
      "components": [
        {
          "componentName": "system_vendor",
          "componentVersion": "00.3.00013103.2024022811453523784981.ac8fcd39"
        },
        {
          "componentName": "my_product",
          "componentVersion": "21609_1310_T.9.0.2024020804070123384722.ec4e3e75"
        },
        {
          "componentName": "my_heytap",
          "componentVersion": "11.9.6.0.2024020801000123375469.d935a743"
        },
        {
          "componentName": "my_stock",
          "componentVersion": "26.4.2.3.2024022816265923786726.8561c7ae"
        },
        {
          "componentName": "my_region",
          "componentVersion": "13.1.02331.2024020700002523354770.15bc037d"
        },
        {
          "componentName": "my_carrier",
          "componentVersion": "13.1.02331.2024020800002523375164.e151df65"
        },
        {
          "componentName": "my_bigball",
          "componentVersion": "21609_1310_T.0.0.2024020704012523361245.7bd48ce3"
        },
        {
          "componentName": "my_manifest",
          "componentVersion": "RMX3350_11.F.25_3250_202403011232.97.d78aa642"
        }
      ],
      "mode": 0,
      "time": 1716827580598,
      "isRooted": "1",
      "type": "1",
      "registrationId": "realme_CN_b397690c4935061a9b933ecae707cb07",
      "securityPatch": "2024-03-05",
      "securityPatchVendor": "2024-03-05",
      "strategyVersion": "9",
      "cota": {
        "cotaVersion": "",
        "cotaVersionName": "",
        "buildType": "user"
      },
      "opex": {
        "check": true
      },
      "sota": {
        "sotaProtocolVersion": "1",
        "sotaVersion": "new",
        "otaUpdateTime": -1,
        "updateViaReboot": 2
      },
      "deviceId": "eb7c6e6f57a623fd267c????????????????????002c25c0fb5c85eac1de2bbd",
      "duid": "BACB337B80074E64A3398AC????????????????????C9C5130BCF0053B63A428",
      "h5LinkVersion": 3
    },
    "iv": "kcolHbUsM0IdhDvIF499zA=="
  }
}
```

### Oplus ota check rooted

```java
    String verityMode = SystemProperties.get("ro.boot.veritymode", "");
    String deviceState = SystemProperties.get("ro.boot.vbmeta.device_state", "");
    return "enforcing".equals(verityMode) || ("eio".equals(verityMode) && "locked".equals(deviceState));
```

## Response

### Response Cipher
```json
{
    "parent": "ota-template",
    "components": [
        {
            "componentId": "my_manifest_10010111.202201100200494684051",
            "componentName": "my_manifest",
            "componentVersion": "10010111.202201100200494684051",
            "componentPackets": {
                "size": "1224301",
                "manualUrl": "https://gauss-componentotacostmanual-cn.allawnfs.com/remove-3ec5e9087e94887cbaf403af49048652/component-ota/22/01/12/040d0ba3022144f1b0849615f733d9ee.zip",
                "id": "domestic_my_manifest_10010111.202201100200494684051_1_ae758a87d35bbc9bcf522c2737eda6f5",
                "type": "1",
                "url": "https://gauss-compotacostauto-cn.allawnfs.com/remove-3ec5e9087e94887cbaf403af49048652/component-ota/22/01/12/040d0ba3022144f1b0849615f733d9ee.zip",
                "md5": "ae758a87d35bbc9bcf522c2737eda6f5"
            }
        },
        {
            "componentId": "my_stock_20615.6.25.202201100200494684051",
            "componentName": "my_stock",
            "componentVersion": "20615.6.25.202201100200494684051",
            "componentPackets": {
                "size": "977303936",
                "manualUrl": "https://gauss-componentotacostmanual-cn.allawnfs.com/remove-3ec5e9087e94887cbaf403af49048652/component-ota/22/01/12/f333ad742dfd4261a5ec82af1455c54a.zip",
                "id": "domestic_my_stock_20615.6.25.202201100200494684051_1_70d4469db58ffe3ed3527b53e303a331",
                "type": "1",
                "url": "https://gauss-compotacostauto-cn.allawnfs.com/remove-3ec5e9087e94887cbaf403af49048652/component-ota/22/01/12/f333ad742dfd4261a5ec82af1455c54a.zip",
                "md5": "70d4469db58ffe3ed3527b53e303a331"
            }
        },
        {
            "componentId": "my_heytap_20615.7.9.202201100200494684049",
            "componentName": "my_heytap",
            "componentVersion": "20615.7.9.202201100200494684049",
            "componentPackets": {
                "size": "658009885",
                "manualUrl": "https://gauss-componentotacostmanual-cn.allawnfs.com/remove-3ec5e9087e94887cbaf403af49048652/component-ota/22/01/12/3f6e477367f44662aaac9c976e05e901.zip",
                "id": "domestic_my_heytap_20615.7.9.202201100200494684049_1_288329f6325674da2abe12a3ab4c95f9",
                "type": "1",
                "url": "https://gauss-compotacostauto-cn.allawnfs.com/remove-3ec5e9087e94887cbaf403af49048652/component-ota/22/01/12/3f6e477367f44662aaac9c976e05e901.zip",
                "md5": "288329f6325674da2abe12a3ab4c95f9"
            }
        },
        {
            "componentId": "my_carrier_20615.7.9.202201101143504698468",
            "componentName": "my_carrier",
            "componentVersion": "20615.7.9.202201101143504698468",
            "componentPackets": {
                "size": "2929642",
                "manualUrl": "https://gauss-componentotacostmanual-cn.allawnfs.com/remove-3ec5e9087e94887cbaf403af49048652/component-ota/22/01/12/4b77731588704785a033e83127ae70a5.zip",
                "id": "domestic_my_carrier_20615.7.9.202201101143504698468_1_98a09a0717c23c3accec1785dc5373e5",
                "type": "1",
                "url": "https://gauss-compotacostauto-cn.allawnfs.com/remove-3ec5e9087e94887cbaf403af49048652/component-ota/22/01/12/4b77731588704785a033e83127ae70a5.zip",
                "md5": "98a09a0717c23c3accec1785dc5373e5"
            }
        },
        {
            "componentId": "21609_system_vendor_20615.1.60.202201110003084707238",
            "componentName": "system_vendor",
            "componentVersion": "20615.1.60.202201110003084707238",
            "componentPackets": {
                "size": "2066055481",
                "manualUrl": "https://gauss-componentotacostmanual-cn.allawnfs.com/remove-3ec5e9087e94887cbaf403af49048652/component-ota/22/01/12/b3a202f878a5461ca81874aacb37f689.zip",
                "id": "domestic_21609_system_vendor_20615.1.60.202201110003084707238_1_74ac01748d7fd3a653ef28717aaeacb5",
                "type": "1",
                "url": "https://gauss-compotacostauto-cn.allawnfs.com/remove-3ec5e9087e94887cbaf403af49048652/component-ota/22/01/12/b3a202f878a5461ca81874aacb37f689.zip",
                "md5": "74ac01748d7fd3a653ef28717aaeacb5"
            }
        },
        {
            "componentId": "my_region_20615.7.9.202201100200494684050",
            "componentName": "my_region",
            "componentVersion": "20615.7.9.202201100200494684050",
            "componentPackets": {
                "size": "189994921",
                "manualUrl": "https://gauss-componentotacostmanual-cn.allawnfs.com/remove-3ec5e9087e94887cbaf403af49048652/component-ota/22/01/12/69b80591935647f0a977421ef70c0a92.zip",
                "id": "domestic_my_region_20615.7.9.202201100200494684050_1_00f06286c26e10e359f8b28a011c4773",
                "type": "1",
                "url": "https://gauss-compotacostauto-cn.allawnfs.com/remove-3ec5e9087e94887cbaf403af49048652/component-ota/22/01/12/69b80591935647f0a977421ef70c0a92.zip",
                "md5": "00f06286c26e10e359f8b28a011c4773"
            }
        }
    ],
    "securityPatch": "2022-01-05",
    "otaVersion": "RMX3350_11.A.16_0160_202201110217",
    "isNvDescription": false,
    "description": {
        "share": ".",
        "panelUrl": "https://gauss-compotacostauto-cn.allawnfs.com/remove-3ec5e9087e94887cbaf403af49048652/component-ota/22/01/24/ef818de1ef714d8f81cbe6a5349e3082.html",
        "url": "https://h5fs.coloros.com/c1eeea21f88343e48612a8c5851cee23/static/index.html#/about",
        "firstTitle": "本次更新包含安全补丁更新，提升系统稳定性，优化功耗，修复一些已知问题"
    },
    "versionName": "RMX3350_11_A.16",
    "rid": "ea2fee92-2e1a-4ce5-946c-f9506eb35331",
    "isRecruit": false,
    "realAndroidVersion": "Android 11",
    "isSecret": false,
    "realOsVersion": "ColorOS 11.1",
    "osVersion": "ColorOS 11.1",
    "publishedTime": 1643244613390,
    "componentAssembleType": true,
    "descriptionType": 0,
    "googlePatchInfo": "0",
    "id": "61f1ec2d124eec00cf83d1b8",
    "group": "RMX3350_11.A.16_0160_202201110217",
    "colorOSVersion": "ColorOS 11.1",
    "timestampH5": "2022.01.24",
    "paramFlag": 1,
    "reminderType": 0,
    "noticeType": 0,
    "decentralize": {
        "strategyVersion": "9",
        "round": 28800,
        "offset": 9787
    },
    "versionCode": 160,
    "silenceUpdate": 0,
    "otaTemplateId": "61f0c52c66c96100cfa266f5",
    "securityPatchVendor": "2022-01-05",
    "androidVersion": "Android 11",
    "versionTypeH5": "正式版",
    "aid": "RMX3350NV97_11.A",
    "nvId16": "NV97",
    "status": "published"
}
```
