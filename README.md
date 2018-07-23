# httpdoc

Auto genrate http interface document.

JSON 文档格式

```json
{
  "method": "GET",
  "url": "http://host:port/api/getXXXList",
  "param": {
    "size": {
      "type": "int",
      "require": false,
      "comment": "分页大小",
      "default": 10
    },
    "offset": {
      "type": "int",
      "require": false,
      "comment": "分页偏移"
    }
  },
  "return": {
    "data": [
      {
        "id": {
          "type": "string",
          "comment": "数据id"
        },
        "name": {
          "type": "string",
          "comment": "数据名"
        }
      }
    ],
    "total": "int"
  }
}
```

TOML 文档格式

```toml
# TODO
```
