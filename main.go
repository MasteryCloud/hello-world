package main

import (
	"fmt"
	"net/http"
	"bytes"
	"text/template"
  "os"
  "log"
)

const defaultListenPort = "80"

type HelloWorldConfig struct {
	Hostname string
	Headers  http.Header
	Host     string
}

const (
	reqInfoID = "reqInfo"

	webHead = `<html>
  <head>
    <title>MasteryCloud</title>
    <link rel="icon" href="img/favicon.png type="image/png">
    <style>
      body {
        margin: 2em 0em;
        padding: 0;
        font-size: 28px;
        font-family: -apple-system,BlinkMacSystemFont,"Roboto","Segoe UI","Helvetica Neue","Lucida Grande",Arial,sans-serif;
        line-height: 1.5;
        position: relative;
        display: flex;
        flex-direction: column;
      }
      a {
        /* text-decoration-line: none; */
        /* text-decoration: none; */
        color: #000;
        text-decoration: none;
      }
      .row {
        padding-top: 1em;
        padding-bottom: 1em;
      }
      .center {
        margin: auto;
        width:fit-content;
      }
      .brand {
        font-size: 32px;
        font-family: -apple-system,BlinkMacSystemFont,"Roboto","Segoe UI","Helvetica Neue","Lucida Grande",Arial,sans-serif;
        display: flex;
        flex-direction: row;
        clear: both;
      }
      .title {
        transition: all 0.2s ease-in-out;
        align-self: center;
        font-weight: bold;
        z-index: 20;
        display: block;
        margin: 0 auto;
        margin-left: 1rem;
      }
      .subtitle {
        font-weight: bold;
        transition: all 0.2s ease-in-out;
        display: block;
        font-size: .625em;
      }
      .logo img {
        max-height: 10rem;
      }

      .footer {
        background-color: #000;
        text-align: left;
      }
      .social-icons li {
        color: #fff !important;
        orphans: 3;
        widows: 3;
        list-style-type: none;
        margin-bottom: 0.5em;
        display: inline-block;
        padding-top: 5px;
        padding-bottom: 5px;
        font-size: .75em;
        text-transform: uppercase;
      }
      ul {
        -webkit-box-direction: normal;
        color: #fff !important;
        box-sizing: border-box;
        orphans: 3;
        widows: 3;
        margin: 0;
        padding: 0;
        list-style-type: none;
      }
      .social-icons a {
        orphans: 3;
        widows: 3;
        list-style-type: none;
        font-size: .75em;
        text-transform: uppercase;
        transition: all 0.2s ease-in-out;
        color: inherit;
        text-decoration: none;
        padding-right: 10px;
        font-weight: bold;
        white-space: nowrap;
      }
      .social-icons a:hover {
        text-decoration: underline;
      }
      .social-icons svg {
        orphans: 3;
        widows: 3;
        list-style-type: none;
        text-transform: uppercase;
        font-weight: bold;
        white-space: nowrap;
        display: inline-block;
        font-size: inherit;
        height: 1em;
        vertical-align: -.125em;
        text-align: center;
        width: 1.25em;
        overflow: visible;
        color: inherit;
      }
      .page__footer-copyright {
        line-height: 1.5;
        color: #fff !important;
        font-family: -apple-system,BlinkMacSystemFont,"Roboto","Segoe UI","Helvetica Neue","Lucida Grande",Arial,sans-serif;
        font-size: .6875em;
      }
      .details {
        text-align: center;
        font-size: 12px;
      }
      .details button {
        margin-bottom: 1em;
      }
    </style>
  </head>
  <body>
  <div class="container">
    <div class="row">
	    <div class="center brand">
        <a class="logo" href="https://masterycloud.com">
        <img src="img/masterycloud-logo.png" alt="MasteryCloud">
        </a>
        <a class="title" href="https://masterycloud.com">
        MasteryCloud
        <span class="subtitle">Building Cloud Infrastructure at Scale</span>
        </a>
      </div>
    </div>
    <div class="row">
      <div class="center">
        My hostname is {{.Hostname}}
      </div>
    </div>`

  webDetails = `
    <div class="row details">
      <div class="center">
        <button class='button' onclick='myFunction()'>Show request details</button>
        <div id="` + reqInfoID + `" style='display:none'>
          <b>Host:</b> {{.Host}} <br />
          <b>Hostname:</b> {{.Hostname}} </b><br />
        {{ range $k,$v := .Headers }}
          <b>{{ $k }}:</b> {{ $v }}<br />
        {{ end }}
        </div>
      </div>
    </div>`

  webLinks = `
      <div class="row footer">
        <div class="center social-icons">
          <ul>
            <li><a href="https://twitter.com/masterycloud" rel="nofollow noopener noreferrer"><svg class="svg-inline--fa fa-twitter-square fa-w-14 fa-fw" aria-hidden="true" focusable="false" data-prefix="fab" data-icon="twitter-square" role="img" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 448 512" data-fa-i2svg=""><path fill="currentColor" d="M400 32H48C21.5 32 0 53.5 0 80v352c0 26.5 21.5 48 48 48h352c26.5 0 48-21.5 48-48V80c0-26.5-21.5-48-48-48zm-48.9 158.8c.2 2.8.2 5.7.2 8.5 0 86.7-66 186.6-186.6 186.6-37.2 0-71.7-10.8-100.7-29.4 5.3.6 10.4.8 15.8.8 30.7 0 58.9-10.4 81.4-28-28.8-.6-53-19.5-61.3-45.5 10.1 1.5 19.2 1.5 29.6-1.2-30-6.1-52.5-32.5-52.5-64.4v-.8c8.7 4.9 18.9 7.9 29.6 8.3a65.447 65.447 0 0 1-29.2-54.6c0-12.2 3.2-23.4 8.9-33.1 32.3 39.8 80.8 65.8 135.2 68.6-9.3-44.5 24-80.6 64-80.6 18.9 0 35.9 7.9 47.9 20.7 14.8-2.8 29-8.3 41.6-15.8-4.9 15.2-15.2 28-28.8 36.1 13.2-1.4 26-5.1 37.8-10.2-8.9 13.1-20.1 24.7-32.9 34z"></path></svg><!-- <i class="fab fa-fw fa-twitter-square" aria-hidden="true"></i> --> Twitter</a></li>
            <li><a href="https://github.com/MasteryCloud" rel="nofollow noopener noreferrer"><svg class="svg-inline--fa fa-github fa-w-16 fa-fw" aria-hidden="true" focusable="false" data-prefix="fab" data-icon="github" role="img" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 496 512" data-fa-i2svg=""><path fill="currentColor" d="M165.9 397.4c0 2-2.3 3.6-5.2 3.6-3.3.3-5.6-1.3-5.6-3.6 0-2 2.3-3.6 5.2-3.6 3-.3 5.6 1.3 5.6 3.6zm-31.1-4.5c-.7 2 1.3 4.3 4.3 4.9 2.6 1 5.6 0 6.2-2s-1.3-4.3-4.3-5.2c-2.6-.7-5.5.3-6.2 2.3zm44.2-1.7c-2.9.7-4.9 2.6-4.6 4.9.3 2 2.9 3.3 5.9 2.6 2.9-.7 4.9-2.6 4.6-4.6-.3-1.9-3-3.2-5.9-2.9zM244.8 8C106.1 8 0 113.3 0 252c0 110.9 69.8 205.8 169.5 239.2 12.8 2.3 17.3-5.6 17.3-12.1 0-6.2-.3-40.4-.3-61.4 0 0-70 15-84.7-29.8 0 0-11.4-29.1-27.8-36.6 0 0-22.9-15.7 1.6-15.4 0 0 24.9 2 38.6 25.8 21.9 38.6 58.6 27.5 72.9 20.9 2.3-16 8.8-27.1 16-33.7-55.9-6.2-112.3-14.3-112.3-110.5 0-27.5 7.6-41.3 23.6-58.9-2.6-6.5-11.1-33.3 2.6-67.9 20.9-6.5 69 27 69 27 20-5.6 41.5-8.5 62.8-8.5s42.8 2.9 62.8 8.5c0 0 48.1-33.6 69-27 13.7 34.7 5.2 61.4 2.6 67.9 16 17.7 25.8 31.5 25.8 58.9 0 96.5-58.9 104.2-114.8 110.5 9.2 7.9 17 22.9 17 46.4 0 33.7-.3 75.4-.3 83.6 0 6.5 4.6 14.4 17.3 12.1C428.2 457.8 496 362.9 496 252 496 113.3 383.5 8 244.8 8zM97.2 352.9c-1.3 1-1 3.3.7 5.2 1.6 1.6 3.9 2.3 5.2 1 1.3-1 1-3.3-.7-5.2-1.6-1.6-3.9-2.3-5.2-1zm-10.8-8.1c-.7 1.3.3 2.9 2.3 3.9 1.6 1 3.6.7 4.3-.7.7-1.3-.3-2.9-2.3-3.9-2-.6-3.6-.3-4.3.7zm32.4 35.6c-1.6 1.3-1 4.3 1.3 6.2 2.3 2.3 5.2 2.6 6.5 1 1.3-1.3.7-4.3-1.3-6.2-2.2-2.3-5.2-2.6-6.5-1zm-11.4-14.7c-1.6 1-1.6 3.6 0 5.9 1.6 2.3 4.3 3.3 5.6 2.3 1.6-1.3 1.6-3.9 0-6.2-1.4-2.3-4-3.3-5.6-2z"></path></svg><!-- <i class="fab fa-fw fa-github" aria-hidden="true"></i> --> GitHub</a></li>

            <li><a href="https://masterycloud.com/en/feed.xml"><svg class="svg-inline--fa fa-rss-square fa-w-14 fa-fw" aria-hidden="true" focusable="false" data-prefix="fas" data-icon="rss-square" role="img" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 448 512" data-fa-i2svg=""><path fill="currentColor" d="M400 32H48C21.49 32 0 53.49 0 80v352c0 26.51 21.49 48 48 48h352c26.51 0 48-21.49 48-48V80c0-26.51-21.49-48-48-48zM112 416c-26.51 0-48-21.49-48-48s21.49-48 48-48 48 21.49 48 48-21.49 48-48 48zm157.533 0h-34.335c-6.011 0-11.051-4.636-11.442-10.634-5.214-80.05-69.243-143.92-149.123-149.123-5.997-.39-10.633-5.431-10.633-11.441v-34.335c0-6.535 5.468-11.777 11.994-11.425 110.546 5.974 198.997 94.536 204.964 204.964.352 6.526-4.89 11.994-11.425 11.994zm103.027 0h-34.334c-6.161 0-11.175-4.882-11.427-11.038-5.598-136.535-115.204-246.161-251.76-251.76C68.882 152.949 64 147.935 64 141.774V107.44c0-6.454 5.338-11.664 11.787-11.432 167.83 6.025 302.21 141.191 308.205 308.205.232 6.449-4.978 11.787-11.432 11.787z"></path></svg><!-- <i class="fas fa-fw fa-rss-square" aria-hidden="true"></i> --> Feed</a></li>
          </ul>
          <div class="page__footer-copyright">© 2020 MasteryCloud</div>
        </div>
      </div>
    </div>`

	webTail = `    <script>
      function myFunction() {
          var x = document.getElementById("` + reqInfoID + `");
          if (x.style.display === "none") {
              x.style.display = "block";
          } else {
              x.style.display = "none";
          }
      }
    </script>
  </body>
</html>`

	HelloWorldTemplate = webHead + `
  ` + webDetails + `
  ` + webLinks + `
` + webTail
)

func CompileTemplateFromMap(tmplt string, configMap interface{}) (string, error) {
	out := new(bytes.Buffer)
	t := template.Must(template.New("compiled_template").Parse(tmplt))
	if err := t.Execute(out, configMap); err != nil {
		return "", err
	}
	return out.String(), nil
}

func (config *HelloWorldConfig) GetManifest() (string, error) {
	return CompileTemplateFromMap(HelloWorldTemplate, config)
}

func (config *HelloWorldConfig) Init(r *http.Request) {
	config.Hostname, _ = os.Hostname()
	config.Host = r.Host
	config.Headers = r.Header
}

func handler(w http.ResponseWriter, r *http.Request) {
	config := &HelloWorldConfig{}
	config.Init(r)
	data, err := config.GetManifest()
	if err != nil {
		fmt.Fprintln(w, err)
	}

	fmt.Fprint(w, data)
}


func main() {
	webPort := os.Getenv("HTTP_PORT")
	if webPort == "" {
		webPort = defaultListenPort
	}

	fmt.Println("Running http service at", webPort, "port")
	http.HandleFunc("/", handler)
	http.Handle("/img/", http.StripPrefix("/img/", http.FileServer(http.Dir(os.Getenv("PWD") + `/img`))))
  log.Fatal(http.ListenAndServe(":" + webPort, nil))
}