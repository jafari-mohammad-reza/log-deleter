

```
docker run --rm \
  -e ADDR="address to elastic search" \
  -e INDEX="index name" \
  -e QUERY="kibana copied query" \
  delete-query
```

or


```
ADDR="address to elastic search" INDEX="index name"   QUERY="kibana copied query" SCROLL_SIZE="5000"  SLICES="5" ./delete-by-query
```