POST http://localhost:8080/api/users
{
  "_key": "johndoe",
  "name": "John Doe",
  "type": "person",
  "information": {
		"instagram": "johndoe",
		"admin": "true",
		"moderator": "true"
	},
  "descriptions": {
		"de-DE": {
			"title": "John Doe",
			"text": "John Doe ist ein sehr schöner Name."
		}
	},
  "photos": [
    {
      "src": "image.jpg",
      "w": 640,
      "h": 480
    }
  ],
  "comments": [
    {
      "title": "Hey there",
      "text": "I have something to say here...",
      "author": "johndoe",
      "created": "2006-01-02T15:04:05Z"
    }
  ]
}

HTTP 200
[Captures]
userkey: jsonpath "$._key"

GET http://localhost:8080/api/users/{{userkey}}

HTTP 200
[Asserts]
jsonpath "$._key" == {{userkey}}
jsonpath "$.name" == "John Doy"