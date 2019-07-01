# url-shortener
Simple URL shortener backed by sqlite.

Using the API

        $ curl -X POST http://mydomain.com/api/save -d '{"url": "http://google.com"}'
        {"error":"","id":"M","url":"http://mydomain.com/123", "short_url": "123"}

        $ curl -X POST http://mydomain.com/api/save -d '{"url": "http://google.com", "shorturl": "abcd"}'
        {"error":"","id":"M","url":"http://mydomain.com/abcd", "short_url": "abcd"}

There's also a simple web ui available

#### Run in docker:

    docker run -dv /local/data/path:/data \
        -p 1337:1337 \
        -e BASE_URL=http://mydomain.com \
        -e DB_PATH=/data \
        -e API_KEY=abcdefg \
        -e PORT=1337 \
        omropfryslan/shorten
