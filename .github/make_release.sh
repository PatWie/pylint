rm -rf pylint

cd ../
go build pylint-server.go
go build pylint-worker.go
cd .github
mkdir -p pylint
cp ../pylint-server pylint/pylint-server
cp ../pylint-worker pylint/pylint-worker

mkdir -p pylint/scripts
cp ../scripts/run_job.sh pylint/scripts/run_job.sh
cp ../README.md pylint/README.md
cp ../LICENSE pylint/LICENSE
tar -zcvf pylint.tar.gz pylint