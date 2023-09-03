#!/bin/bash
docker run -d --name recontact-mailer -p 7500:7500  --env-file ./env.txt recontact-mailer:latest