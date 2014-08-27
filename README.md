docker-registry-debug
=====================

    Usage: ./registry-debug [options] command
    
    options:
      -i="https://index.docker.io": override index endpoint
      -q=false: disable debug logs
      -r="": override registry endpoint
    
    commands:
      info <repos_name>: lookup repos meta-data
      layerinfo <repos_name> <layer_id>: lookup layer meta-data
      curlme <repos_name> <layer_id>: print a curl command for fetching the layer


Examples:

    ./registry-debug info ubuntu
    ./registry-debug layerinfo ubuntu 3db9c44f45209632d6050b35958829c3a2aa256d81b9a7be45b362ff85c54710
    ./registry-debug curlme ubuntu 3db9c44f45209632d6050b35958829c3a2aa256d81b9a7be45b362ff85c54710


Using via Docker
----------------
Project provides an easy way to build and run the app via Docker. In this way, you don't have to install go and its dependencies to local environment. Usage example:

```bash
# Build the container
docker build -t docker-registry-debug .
# Run the command inside the container
docker run -i -t docker-registry-debug info ubuntu
```
