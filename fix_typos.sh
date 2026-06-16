# Fix some specific issues in Go and JSON files instead of ignoring them globally
sed -i 's/var datas/var data/' agent/utils/firewall/client/ufw.go
sed -i 's/datas = append/data = append/' agent/utils/firewall/client/ufw.go
sed -i 's/return datas/return data/' agent/utils/firewall/client/ufw.go

sed -i 's/var datas/var data/g' agent/utils/firewall/client/firewalld.go
sed -i 's/datas = append/data = append/g' agent/utils/firewall/client/firewalld.go
sed -i 's/return datas/return data/g' agent/utils/firewall/client/firewalld.go

sed -i 's/var datas/var data/g' agent/utils/cloud_storage/client/kodo.go
sed -i 's/datas = append/data = append/g' agent/utils/cloud_storage/client/kodo.go
sed -i 's/return datas/return data/g' agent/utils/cloud_storage/client/kodo.go

sed -i 's/var datas/var data/g' agent/utils/cloud_storage/client/cos.go
sed -i 's/datas = append/data = append/g' agent/utils/cloud_storage/client/cos.go
sed -i 's/return datas/return data/g' agent/utils/cloud_storage/client/cos.go
sed -i 's/datas, _, err/data, _, err/g' agent/utils/cloud_storage/client/cos.go
sed -i 's/datas\.Contents/data.Contents/g' agent/utils/cloud_storage/client/cos.go

sed -i 's/fo := files.NewFileOp()/fOp := files.NewFileOp()/' agent/utils/cloud_storage/client/kodo.go
sed -i 's/fo\.DownloadFile/fOp.DownloadFile/' agent/utils/cloud_storage/client/kodo.go

sed -i 's/var entrys/var entries/g' agent/utils/toolbox/pure-ftpd.go
sed -i 's/entrys = append/entries = append/g' agent/utils/toolbox/pure-ftpd.go
sed -i 's/range entrys/range entries/g' agent/utils/toolbox/pure-ftpd.go

sed -i 's/Broswer/Browser/g' agent/utils/nginx/components/location.go

sed -i 's/var datas/var data/g' agent/utils/postgresql/client/remote.go
sed -i 's/datas = append/data = append/g' agent/utils/postgresql/client/remote.go
sed -i 's/return datas/return data/g' agent/utils/postgresql/client/remote.go

sed -i 's/var datas/var data/g' agent/utils/postgresql/client/local.go
sed -i 's/datas = append/data = append/g' agent/utils/postgresql/client/local.go
sed -i 's/return datas/return data/g' agent/utils/postgresql/client/local.go

sed -i 's/"formatEN": "clean monitor datas"/"formatEN": "clean monitor data"/g' core/cmd/server/docs/x-log.json
sed -i 's/"formatEN": "clean monitor datas"/"formatEN": "clean monitor data"/g' core/cmd/server/docs/docs.go
sed -i 's/"formatEN": "clean monitor datas"/"formatEN": "clean monitor data"/g' core/cmd/server/docs/swagger.json

sed -i 's/"summary": "Sycn host SSH secret"/"summary": "Sync host SSH secret"/g' core/cmd/server/docs/docs.go
sed -i 's/"summary": "Sycn host SSH secret"/"summary": "Sync host SSH secret"/g' core/cmd/server/docs/swagger.json

sed -i 's/Lable         string    `json:"lable"`/Label         string    `json:"label"`/g' core/app/dto/script_library.go

sed -i 's/datas := writer/data := writer/g' core/middleware/operation.go
sed -i 's/datas, _ = io/data, _ = io/g' core/middleware/operation.go
sed -i 's/Unmarshal(datas/Unmarshal(data/g' core/middleware/operation.go

sed -i 's/isRuning/isRunning/g' core/utils/firewall/firewall.go
