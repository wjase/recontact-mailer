# recontact-mailer

This web app takes a json post request from a contact form and emails
the contents to a specified address if the recaptacha check succeeds.

# Server Setup

 1. Create a recaptcha account for your domain at [Google](https://developers.google.com/recaptcha/docs/display)

 2. Configure your settings in the `/etc/recontact-mailer/env.txt` file. [NOTE: Protect this file as it contains sensitive info]

  ```# The private key from Google
  RECAPTCHA_PRIVATE_KEY=
  # The email address you want to receivethe contact details
  TO_MAIL=
  # Email credentials and host address usedfor sending the email
  EMAIL_USERNAME=
  EMAIL_PASSWORD=
  EMAIL_HOST=
  EMAIL_PORT=
  # this app exposes / on the portspecified below
  APP_PORT=7500
```

 3. Copy the config/recontact-mailer.service to /etc/systemd/system/

 4. Enable the service to start at system startup 
    
    `systemctl enble recontact-mailer.service` 

 5. Start the service:
 
    `systemctl start recontact-mailer`

