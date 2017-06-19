# Slack Visitor

Slack app that lets your Slack team chat with people who are not part of your Slack team via a web client.

![Slack Visitor Demo](http://g.recordit.co/09ZfILeOZh.gif)


## Deploy to heroku

### Slack App Configuration

1. Create an app in your Slack team: https://api.slack.com/apps?new_app=1

2. Add a Bot user to the Slack app.

3. Enable interactive messages.

4. Add the following permissions to your app:

  - channels:write
  - groups:write
  - users:read

5. Install the app to the team that you want.

### Deploy

1. Create a Heroku app.

2. Add the nodejs and the go buildpacks:

    ```
    $ heroku buildpacks:add --index 1 heroku/nodejs
    $ heroku buildpacks:add --index 2 heroku/go
    ```
3. Add a Postgres database to the heroku app
4. Setup the following env vars:

  - VISITOR_CHANNEL_ID=XXXXXXXXXX (The ID of the channel where the chat requests will be displayed)
  - SLACK_BOT_API_KEY=xoxb-xxxxxxxxxx-xxxxxxxxx-xxxxxxxx
  - SLACK_APP_API_KEY=xoxp-yyyyyyyyy-yyyyyyyyyy-yyyyyyyy
  - SLACK_VERIFICATION_TOKEN=xyxyxyxyxyxyxyxyxy
5. Deploy

    ```
    $ git push heroku master
    ```
6. Run migrations

    ```
    $ heroku run bin/slack-visitor db migrate
    ```
7. Go back to the Slack App configuration page and setup the Interactive Messages Request URL that should look like this `https://<name-of-your-heroku-app>.herokuapp.com/api/slack/action`

8. Enjoy :)

## Development

After checking out the repo, run `glide install && npm install` to install dependencies. Then, run `go run main.go` and `cd client && npm run dev` to run the app.

## Contributing

Bug reports and pull requests are welcome on GitHub at https://github.com/simon0191/slack-visitor.

## License

MIT License

Copyright (c) 2017 Simon Soriano

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
