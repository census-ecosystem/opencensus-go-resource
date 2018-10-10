# OpenCensus resource discovery

This tool runs all known resource discovery mechanisms in the repository and
returns the result encoded in the well-known environment variables.
It can be used to run discovery before launching another process using:

```bash
env $(ocdiscovery) <my_binary>
```

