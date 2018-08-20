package main

var format = `[
  {
    "title":"Osx system log",
    "regex" : [
      {
        "name": "basic",
        "regex": "(?m)(?P<timestamp>^\\w\\w\\w \\d\\d \\d\\d:\\d\\d:\\d\\d) (?P<hostname>[\\w-]*) (?P<body>.*)",
        "timestamp":"Jan 02 15:04:05"
      }
    ]
  },
  {
    "title":"Golang log",
    "regex" : [
      {
        "name": "basic",
        "regex": "(?m)(?P<timestamp>^\\d\\d\\d\\d/\\d\\d/\\d\\d \\d\\d:\\d\\d:\\d\\d) (?P<body>.*)",
        "timestamp":"2006/01/02 15:04:05"
      }
    ]
  },
 {
    "title":"Stern log",
    "regex" : [
      {
        "name": "basic",
        "regex": "(?m)(?P<pod>[\\S-]*) (?P<system>[\\S-]*) (?P<timestamp>\\d\\d\\d\\d-\\d\\d-\\d\\dT\\d\\d:\\d\\d:\\d\\d,\\d\\d\\dZ) (?P<thread>[\\S-]*) (?P<level>[\\S-]*) (?P<body>.*)",
        "timestamp":"2006-01-02T15:04:05.999Z"
      },
      {
        "name": "ml",
        "regex": "(?m)(?P<pod>[\\S-]*) (?P<system>[\\S-]*) (?P<body>.*)"
      }
    ]
  }
]`
