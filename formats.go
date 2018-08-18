package main

var format = `[
  {
    "title":"Osx system log",
    "multiline" : false,
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
    "multiline" : false,
    "regex" : [
      {
        "name": "basic",
        "regex": "(?m)(?P<timestamp>^\\d\\d\\d\\d/\\d\\d/\\d\\d \\d\\d:\\d\\d:\\d\\d) (?P<body>.*)",
        "timestamp":"2006/01/02 15:04:05"
      }
    ]
  }
]`
