## How to use

- Install docker if you havent already
- Create a `.env` based on the `.env.example` and add the webhook url
- Run `docker build . -t bot:latest`
- Run `docker run --env-file ./.env -it bot`
