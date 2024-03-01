## Fam Yt
Caches yt results for predefined search queries for faster access.
Dashboard: [Dashboard](https://fam-yt.sudonick.me)
Backend: [Backend](https://aws.sudonick.me)

## More Info
This project has three components.
1. Server: A http server that serves cached results for certain pre-defined 
search queries.
2. Fetcher: A background process that fetches the results for the pre-defined
search queries and stores them in the database.
3. Web: A dashboard to view the cached results and the search queries. 
You can also add new search queries from the dashboard.

## Stuff Implemented
- [x] Poller and Fetcher (Used worker pool for fetching results)
- [x] Server for serving cached results (With Pagination & Sorting)
- [x] Api key rotation when quota exceeds
- [x] Dashboard UI

## Tech Stack 
- Golang (Echo Framework)
- MongoDB 
- React
- Tailwind CSS 

## Quickstart
Have the env variables ready, you can use the `env.example` file as a template.
```
make server
make fetcher
make web
```

## API Documentation
GET /api/videos
- Query Params:
    - query: search query
    - limit (optional): default 10
    - pageno (optional): default 1
    - order (optional): default "desc" //order by published date
- Returns:
    - Array of videos

GET /api/queries
- Returns:
    - Array of queries

POST /api/queries
- Body:
    - query: search query
- Returns:
    - Success/Error message

GET /api/info
- Returns:
    - pollInterval - interval for fetching results
    - maxResults - fetched by fetcher

## Screenshots
![image](https://github.com/sudo-nick16/fam-yt/assets/73229823/5a2e0246-0358-4bb8-ac8c-2d849ba831c1)

