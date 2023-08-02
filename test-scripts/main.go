package main

import (
  "bytes"
  "crypto/tls"
  "encoding/json"
  "fmt"
  "io/ioutil"
  "math/rand"
  "net/http"
  "net/url"
  "os"
  "strconv"
  "flag"
)

var (
  clientId string
  clientSecret string
  ssoEndpoint string
  apiEndpoint string
)

const (
  AUTH = iota
  REST_WRITE
  REST_READ
)

type authResponse struct {
  AccessToken string `json:"access_token"`
  NotBeforePolicy int `json:"not-before-policy"`
  Scope string `json:"scope"`
}

type queryResponse struct {
  UserId int `json:"user_id"`
  Title string `json:"title"`
  Body string `json:"body"`
}

const (
  tokenPath = "/auth/realms/3Scale-demo/protocol/openid-connect/token"
)

func performRequest(endPoint string, reqType int, insecure bool, token string) []byte {
  var tr *http.Transport;
  requestURL := endPoint;

  if insecure {
    tr = &http.Transport{
      TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
    }
  } else {
    tr = &http.Transport{}
  }

  client := &http.Client{Transport: tr}

  if (reqType == AUTH) {
    // prepare auth payload
    payload := url.Values{};
    payload.Add("client_id", clientId);
    payload.Add("client_secret", clientSecret);
    payload.Add("scope", "openid");
    payload.Add("grant_type", "client_credentials");

    req, err := http.NewRequest(http.MethodPost, requestURL, bytes.NewReader([]byte(payload.Encode())));
    if err != nil {
      fmt.Printf("error building request object: %s\n", err);
      os.Exit(1)
    }
    req.Header.Set("Content-Type", "application/x-www-form-urlencoded");
    req.Header.Add("Content-Length", strconv.Itoa(len(payload.Encode())))

    res, err := client.Do(req);
    if err != nil {
      fmt.Printf("error making http request: %s\n", err);
      os.Exit(1)
    }
    defer res.Body.Close();

    resBody, err := ioutil.ReadAll(res.Body)
    if err != nil {
      fmt.Printf("client: could not read response body: %s\n", err)
      os.Exit(1)
    }

    // check error code
    if (res.StatusCode != 200) {
      fmt.Printf("client: error - status code %d\n", res.StatusCode);
      os.Exit(1);
    }

    // return response
    return resBody;
  } else if (reqType == REST_WRITE) {
    var payloadTemplate string = `{"userId": "%s", "title": "%s", "body": "%s"}`;

    payloadBytes := []byte(fmt.Sprintf(payloadTemplate, strconv.Itoa(rand.Intn(10)), "Test Call number " + strconv.Itoa(rand.Int()), "Test Body Number " + strconv.Itoa(rand.Int())))

    req, err := http.NewRequest(http.MethodPost, requestURL, bytes.NewReader(payloadBytes));
    if err != nil {
      fmt.Printf("error building request object: %s\n", err);
      os.Exit(1)
    }
    req.Header.Set("Content-Type", "application/json");
    req.Header.Set("Authorization", "Bearer " + token);

    res, err := client.Do(req);
    if err != nil {
      fmt.Printf("error making http request: %s\n", err);
      os.Exit(1)
    }
    defer res.Body.Close();

    resBody, err := ioutil.ReadAll(res.Body)
    if err != nil {
      fmt.Printf("client: could not read response body: %s\n", err)
      os.Exit(1)
    }
    fmt.Printf("%s\n", resBody);

    // return response
    return resBody;

  } else if (reqType == REST_READ) {
    req, err := http.NewRequest(http.MethodGet, requestURL, nil);
    if err != nil {
      fmt.Printf("error building request object: %s\n", err);
      os.Exit(1)
    }
    req.Header.Set("Content-Type", "application/json");
    req.Header.Set("Authorization", "Bearer " + token);

    res, err := client.Do(req);
    if err != nil {
      fmt.Printf("error making http request: %s\n", err);
      os.Exit(1)
    }
    defer res.Body.Close();

    resBody, err := ioutil.ReadAll(res.Body)
    if err != nil {
      fmt.Printf("client: could not read response body: %s\n", err)
      os.Exit(1)
    }

    if (res.StatusCode != 200) {
      fmt.Printf("client: error - status code %d\n", res.StatusCode);
      os.Exit(1);
    }

    // return response
    return resBody;
  }

  return []byte{};
}

func validateStringParm(s string) bool {
  if len(s) > 0 {
    return true;
  }

  return false;
}

func main() {
  var r []byte;
  var authDict authResponse;
  var queryDict []queryResponse;

  flag.StringVar(&clientId, "key", "", "The API Key");
  flag.StringVar(&clientSecret, "secret", "", "The API secret");
  flag.StringVar(&ssoEndpoint, "sso", "https://keycloak-fuse-jdbc-demo.apps.demo3scale.sandbox644.opentlc.com", "The SSO Authentication Endpoint");
  flag.StringVar(&apiEndpoint, "api_endpoint", "https://fuse-product-fuse-spring-apicast-staging.apps.demo3scale.sandbox644.opentlc.com", "The API REST Endpoint");

  // parse CMD Line flags
  flag.Parse()

  if ! (validateStringParm(clientId) && validateStringParm(clientSecret)) {
    fmt.Println("Missing Parameters:");
    flag.PrintDefaults();
    os.Exit(1);
  }

  tokenURL := fmt.Sprintf("%s%s", ssoEndpoint, tokenPath);
  fmt.Println("Contacting the SSO server to get a valid access token...");
  r = performRequest(tokenURL, AUTH, true, "");

  err := json.Unmarshal(r, &authDict);
  if (err != nil) {
    fmt.Printf("client error %s\n", err);
    os.Exit(1);
  }

  fmt.Printf("-> Got access Token: [%s]\n", authDict.AccessToken);

  // make a POST call
  postURL := fmt.Sprintf("%s%s", apiEndpoint, "/api/post");
  fmt.Println("Contacting the API Gateway to perform a REST WRITE Operation....");
  r = performRequest(postURL, REST_WRITE, true, authDict.AccessToken);

  fmt.Printf("-> Write Complete\n");

  // make a GET call
  getURL := fmt.Sprintf("%s%s", apiEndpoint, "/api/get");
  fmt.Println("Contacting the API Gateway to perform a REST READ Operation....");
  r = performRequest(getURL, REST_READ, true, authDict.AccessToken);

  err = json.Unmarshal(r, &queryDict);
  if (err != nil) {
    fmt.Printf("client error %s\n", err);
    os.Exit(1);
  }

  fmt.Printf("-> Read Results\n");
  for i := range queryDict {
    fmt.Printf("User ID: %d -- Title: [%s] -- Body: [%s]\n", queryDict[i].UserId, queryDict[i].Title, queryDict[i].Body);
  }
}
