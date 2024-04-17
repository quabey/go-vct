# Go VCT

Sends automated messages with Valorant Champions Tournament Updates to a services webhook.

Updates Inculde

- Upcoming matches
- Starting matches
- Ending machtes

Also theres a selection of offical streams and watch parties.

## How to use

- Install docker if you havent already
- Create a `.env` based on the `.env.example` and add the webhook url
- Run `docker build . -t bot:latest`
- Run `docker run --env-file ./.env -it bot`
