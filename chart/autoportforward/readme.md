helm install yourspotify . -n yourspotify --create-namespace --values values.yaml

 helm upgrade yourspotify . -n yourspotify -f values.yaml