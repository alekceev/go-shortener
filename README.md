# URL Shortener

Проект сокрощателя ссылок на Go

## Backend

# Add new url
curl -sL -H "Content-type: application/json"  localhost:8080/ -XPOST -d '{"url":"https://yandex.ru"}' | json_pp

{
   "short_url" : "/MGRhZjQ4YjQtY2MyMi00MDcwLWE4ZTEtOGJlZGRkN2RkNjZh",
   "stats_url" : "/stats/MGRhZjQ4YjQtY2MyMi00MDcwLWE4ZTEtOGJlZGRkN2RkNjZh"
}

# Get stats
curl -sL localhost:8080/MGRhZjQ4YjQtY2MyMi00MDcwLWE4ZTEtOGJlZGRkN2RkNjZh/stats | json_pp
{
   "id" : "32bec4ed-65bb-4e7f-a56a-298684a714f6",
   "num_redirects" : 0,
   "short_url" : "MGRhZjQ4YjQtY2MyMi00MDcwLWE4ZTEtOGJlZGRkN2RkNjZh",
   "target" : "https://yandex.ru"
}

## Frontend

TODO

## API

TODO
