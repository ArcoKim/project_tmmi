# Project Tmmi
## 설명
국내의 음악 차트를 모아 인사이트를 얻고, AI나 그래프 등을 통해 음악 차트에 편리하게 접근할 수 있습니다.

## 리소스 준비
1. 인프라 배포
```bash
terraform init
terraform apply
```

2. 테이블 & 인덱스 생성
```bash
psql -h $PG_HOST -U postgres -d tmmi -a -f database/init.sql
```

3. 크롤링을 이용한 데이터 주입
```bash
export PG_HOST=(Host)
export PG_PORT=(Port)
export PG_USER=(User Name)
export PG_DATABASE=(Database Name)
export PG_PASSWORD=(Password)

cd database/crawling
go run scraper.go
```

4. Aurora -> S3 내보내기 (Bedrock Knowledge Base Sync가 필요함)
```sql
SELECT aws_commons.create_s3_uri ('ap-tmmi-postgres-XXXX', 'song.csv', 'ap-northeast-2') AS s3_uri \gset
SELECT * FROM aws_s3.query_export_to_s3('SELECT m.name, a.artist, a.name album, m.lyrics FROM music m JOIN album a ON m.album_id = a.id', :'s3_uri', options :='format csv, header true');
```

## 리소스 삭제
```bash
terraform destroy
```