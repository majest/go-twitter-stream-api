package twitter_api

import (
  "github.com/NOX73/go-oauth"
  "net/http"
  "bufio"
)

const (
  NewRequestMethod = "POST"
  NewRequestURL = "https://stream.twitter.com/1.1/statuses/filter.json"
)

type Credentials struct {
  oauth_consumer_key string
  oauth_token string
  oauth_consumer_secret string
  oauth_token_secret string
}

type Message struct {
  Error error
  Response *http.Response
  Tweet *Tweet
}

type Tweet struct {
  Body string
}

func TwitterStream (ch chan Message, credentials *Credentials, params map[string]string){
  var message Message

  c := oauth.NewCredentials(credentials.oauth_consumer_key, credentials.oauth_token, credentials.oauth_consumer_secret, credentials.oauth_token_secret)

  r, _ := oauth.NewRequest(NewRequestMethod, NewRequestURL, params, c)

  client := http.Client{}
  resp, err := client.Do(r.HttpRequest())

  if err != nil {
    message = Message{
      Error:err,
      Response: resp,
    }

    ch <- message
    return
  }

  body_reader := bufio.NewReader(resp.Body)

  for {
    var part []byte //Part of line
    var prefix bool //Flag. Readln readed only part of line.

    part, prefix, err := body_reader.ReadLine()
    if err != nil { break }

    buffer := append([]byte(nil), part...)

    for prefix && err == nil {
      part, prefix, err = body_reader.ReadLine()
      buffer = append(buffer, part...)
    }
    if err != nil { break }

    tweet := &Tweet{
      Body: string(buffer),
    }

    message = Message{
      Response: resp,
      Tweet: tweet,
    }

    ch <- message
  }
}
