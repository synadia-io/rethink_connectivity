package main

#Response: {
	@jsonschema(schema="https://json-schema.org/draft/2020-12/schema")
	lang: string
	words: [...{
		word:  string
		score: int
	}]
	sentences?: [...{
		sentence: string
		score:    int
	}]
	score: int
}
