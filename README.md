# Mine klager - Min side microfrontend

Static HTML microfrontend served by NGINX.

## Routes

- `/nb` - Norwegian Bokmål
- `/nn` - Norwegian Nynorsk
- `/en` - English

## Local Development

Build and run with Docker:

```sh
./scripts/build.sh
docker build -t mine-klager-microfrontend .
docker run -p 8080:8080 mine-klager-microfrontend
```

Then open:
- [localhost:8080/nb](http://localhost:8080/nb)
- [localhost:8080/nn](http://localhost:8080/nn)
- [localhost:8080/en](http://localhost:8080/en)
