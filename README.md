[![Build Status](https://travis-ci.org/indix/abelwatch.svg?branch=master)](https://travis-ci.org/indix/abelwatch)
# abelwatch

Abel Watch is an alerting tool for Abel stats aggregation system.

## Usage
```
Usage of ./abelwatch:
  -abel-url string
    	Abel URL (Eg. http://abel.domain.tld:3330) without the trailing slash
  -pid string
    	File to write PID file (default "PID")
  -slack-webhook string
    	Slack webhook to post the alert
  -wasp-namespace string
    	Namespace in WASP to get the AbelWatch rules (default "dev.abel.watchers.rules")
  -wasp-url string
    	WASP URL (Eg. http://wasp.domain.tld:9000) without the trailing slash
```

## License
https://www.apache.org/licenses/LICENSE-2.0
