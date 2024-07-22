###Install ctop on coreos
#make bashrc editable on coreos
cd $HOME
rm .bashrc
cp /usr/share/skel/.bashrc .
#add alias to bashrc on coreos
alias monitor="docker run --rm -ti --name=ctop -v /var/run/docker.sock:/var/run/docker.sock quay.io/vektorlab/ctop:latest"
#RUN and relax
monitor

Copy ./config-old to ./config as graylog uses it through docker and locks the directory

Run main.go to create a message

uses admin:admin still

curl -u admin:admin -H 'Accept: application/json' -X GET 'http://0.0.0.0:9000/api/cluster?pretty=true'


curl -u admin:admin -i -H 'Accept: application/json' 'http://localhost:9000/api/search/universal/absolute?query=*&from=2017-12-19T14:16:56.329Z&to=2017-12-20T15:00:00.000Z&fields=timstamp,source,message&limit=5&pretty=true'

Creates a couple of messages on the std index of graylog_0
go run main.go

Queries the messages from elastic exposed via graylog rest api
 go run elastic/main.go

http://docs.graylog.org/en/2.3/pages/archiving.html -> message ttl


# CERTS

cd ./graylog/certs

openssl req -x509 -days 365 -nodes -newkey rsa:2048 -config openssl-graylog.cnf -keyout pkcs5-plain.pem -out cert.pem

openssl pkcs8 -in pkcs5-plain.pem -topk8 -nocrypt -out pkcs8-plain.pem

openssl pkcs8 -in pkcs5-plain.pem -topk8 -out graylog-key.pem -passout pass:secret

to accept those certs in chrome go to

chrome://settings/certificates

click on server/import and choose:
graylog/certs/cert.pem
this will show up under others tab.
Chrome will still give a warning that the connection is not trusted but it works