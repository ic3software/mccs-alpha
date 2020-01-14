package es

import (
	"context"
	"log"

	"github.com/olivere/elastic/v7"
)

func checkIndexes(client *elastic.Client) {
	for _, indexName := range indexes {
		checkIndex(client, indexName)
	}
}

func checkIndex(client *elastic.Client, index string) {
	ctx := context.Background()

	exists, err := client.IndexExists(index).Do(ctx)
	if err != nil {
		panic(err)
	}

	if exists {
		return
	}

	createIndex, err := client.CreateIndex(index).BodyString(indexMappings[index]).Do(ctx)
	if err != nil {
		panic(err)
	}
	if !createIndex.Acknowledged {
		panic("CreateIndex " + index + " was not acknowledged.")
	} else {
		log.Println("Successfully created " + index + " index")
	}
}

var indexes = []string{"businesses", "users", "tags"}

// Notes:
// 1. Using nested fields for arrays of objects.
var indexMappings = map[string]string{
	"businesses": `
	{
		"settings": {
			"analysis": {
				"analyzer": {
					"tag_analyzer": {
						"type": "custom",
						"tokenizer": "whitespace",
						"filter": [
							"lowercase",
							"asciifolding"
						]
					}
				}
			}
		},
		"mappings": {
			"properties": {
				"adminTags": {
					"type": "text",
					"analyzer": "tag_analyzer",
					"fields": {
						"keyword": {
							"type": "keyword",
							"ignore_above": 256
						}
					}
				},
				"businessID": {
					"type": "text",
					"fields": {
						"keyword": {
							"type": "keyword",
							"ignore_above": 256
						}
					}
				},
				"businessName": {
					"type": "text",
					"fields": {
						"keyword": {
							"type": "keyword",
							"ignore_above": 256
						}
					}
				},
				"locationCity": {
					"type": "text",
					"fields": {
						"keyword": {
							"type": "keyword",
							"ignore_above": 256
						}
					}
				},
				"locationCountry": {
					"type": "text",
					"fields": {
						"keyword": {
							"type": "keyword",
							"ignore_above": 256
						}
					}
				},
				"status": {
					"type": "text",
					"fields": {
						"keyword": {
							"type": "keyword",
							"ignore_above": 256
						}
					}
				},
				"offers": {
					"type" : "nested",
					"properties": {
						"createdAt": {
							"type": "date"
						},
						"name": {
							"type": "text",
							"analyzer": "tag_analyzer",
							"fields": {
								"keyword": {
									"type": "keyword",
									"ignore_above": 256
								}
							}
						}
					}
				},
				"wants": {
					"type" : "nested",
					"properties": {
						"createdAt": {
							"type": "date"
						},
						"name": {
							"type": "text",
							"analyzer": "tag_analyzer",
							"fields": {
								"keyword": {
									"type": "keyword",
									"ignore_above": 256
								}
							}
						}
					}
				}
			}
		}
	}`,
	"users": `
	{
		"mappings": {
			"properties": {
				"email": {
					"type": "text",
					"fields": {
						"keyword": {
							"type": "keyword",
							"ignore_above": 256
						}
					}
				},
				"firstName": {
					"type": "text",
					"fields": {
						"keyword": {
							"type": "keyword",
							"ignore_above": 256
						}
					}
				},
				"lastName": {
					"type": "text",
					"fields": {
						"keyword": {
							"type": "keyword",
							"ignore_above": 256
						}
					}
				},
				"userID": {
					"type": "text",
					"fields": {
						"keyword": {
							"type": "keyword",
							"ignore_above": 256
						}
					}
				}
			}
		}
	}`,
	"tags": `
	{
		"settings": {
			"analysis": {
				"analyzer": {
					"tag_analyzer": {
						"type": "custom",
						"tokenizer": "whitespace",
						"filter": [
							"lowercase",
							"asciifolding"
						]
					}
				}
			}
		},
		"mappings": {
			"properties": {
				"name": {
					"type": "text",
					"analyzer": "tag_analyzer",
					"fields": {
						"keyword": {
							"type": "keyword",
							"ignore_above": 256
						}
					}
				},
				"offerAddedAt": {
					"type": "date"
				},
				"tagID": {
					"type": "text",
					"fields": {
						"keyword": {
							"type": "keyword",
							"ignore_above": 256
						}
					}
				},
				"wantAddedAt": {
					"type": "date"
				}
			}
		}
	}`,
}
