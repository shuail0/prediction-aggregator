<!-- 源: https://docs.polymarket.com/api-reference/sports/get-valid-sports-market-types -->

GET

/

sports

/

market-types

Try it

Get valid sports market types

cURL

Copy

Ask AI

```python
curl --request GET \
  --url https://gamma-api.polymarket.com/sports/market-types
```

200

Copy

Ask AI

```python
{
  "marketTypes": [
    "<string>"
  ]
}
```

#### Response

200 - application/json

List of valid sports market types

[​](#response-market-types)

marketTypes

string[]

List of all valid sports market types

[Get sports metadata information](/api-reference/sports/get-sports-metadata-information)[List tags](/api-reference/tags/list-tags)

⌘I