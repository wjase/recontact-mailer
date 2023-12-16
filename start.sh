#!/bin/bash
docker run -d --name recontact-mailer -p 7500:7500  --env-file ./test-env.txt recontact-mailer:latest