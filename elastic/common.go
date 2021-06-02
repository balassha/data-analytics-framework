package elastic

type keyword struct {
	Type string `json:"type"`
}

type properties struct {
	Postcode keyword `json:"postcode"`
	Recipe   keyword `json:"recipe"`
	Delivery keyword `json:"delivery"`
	Start    keyword `json:"start"`
	End      keyword `json:"end"`
}

type mapping struct {
	Properties properties `json:"properties"`
}

type IndexName struct {
	Index string `json:"_index"`
}

type IndexWrapper struct {
	Index IndexName `json:"index"`
}

type Data struct {
	Postcode string `json:"postcode"`
	Recipe   string `json:"recipe"`
	Delivery string `json:"delivery"`
	Start    int    `json:"start"`
	End      int    `json:"end"`
}

// Aggregator query to find unique recipes
// {
//     "aggs": {
//         "recipes": {
//             "terms": {
//                 "field": "recipe"
//             }
//         }
//     }
// }
type uniqueRecipes struct {
	Aggregator aggregator `json:"aggs"`
}

type aggregator struct {
	Recipes terms `json:"recipes"`
}

type terms struct {
	Terms field `json:"terms"`
}

type field struct {
	Field string `json:"field"`
	Size  *int   `json:"size,omitempty"`
}

// Response fields
type aggregatorResponse struct {
	Aggregations recipes `json:"aggregations"`
	Hits         *hits   `json:"hits,omitempty"`
}

type hits struct {
	Total total       `json:"total"`
	Hits  []hitsArray `json:"hits"`
}

type hitsArray struct {
	Array fields `json:"fields"`
}

type fields struct {
	Recipe []string `json:"recipe"`
}

type total struct {
	Value int `json:"value"`
}

type recipes struct {
	Recipe buckets `json:"recipes"`
}

type buckets struct {
	Buckets []RecipeCount `json:"buckets"`
}

type RecipeCount struct {
	Key   string `json:"key"`
	Count int    `json:"doc_count"`
}

// List the recipe names (alphabetically ordered) that contain in their name
// {
//     "query": {
//         "query_string": {
//             "query": "(*Veggie*) OR (*Chops*) OR (*Mushroom*)",
//             "default_field": "recipe"
//         }
//     },
//     "collapse": {
//         "field": "recipe"
//     }
// }
type regexQuery struct {
	Query    queryString `json:"query"`
	Collapse collapse    `json:"collapse"`
}

type collapse struct {
	Field string `json:"field"`
}

type queryString struct {
	QueryString queryDefaultField `json:"query_string"`
}

type queryDefaultField struct {
	Query        string `json:"query"`
	DefaultField string `json:"default_field"`
}

// Aggregator query to find the postcode with most delivered recipes.
// Reuse some of the types from above
// {
//     "aggs": {
//         "recipes": {
//             "terms": {
//                 "field": "postcode",
//                 "size":"1"
//             }
//         }
//     }
// }

// Query to Get deliveries on a postcode within the given time
// {
//     "query": {
//         "bool": {
//             "must": [
//                 {
//                     "match": {
//                         "postcode": "10120"
//                     }
//                 },
//                 {
//                     "range": {
//                         "start": {
//                             "gte": 10
//                         }
//                     }
//                 },
//                 {
//                     "range": {
//                         "end": {
//                             "lte": 15
//                         }
//                     }
//                 }
//             ]
//         }
//     }
// }

type query struct {
	Query boolType `json:"query"`
}

type boolType struct {
	Bool must `json:"bool"`
}

type must struct {
	Must []matchRange `json:"must"`
}

type matchRange struct {
	Match *match     `json:"match,omitempty"`
	Range *rangeType `json:"range,omitempty"`
}

type match struct {
	Postcode string `json:"postcode"`
}

type rangeType struct {
	Start *start `json:"start,omitempty"`
	End   *end   `json:"end,omitempty"`
}

type start struct {
	Gte *int `json:"gte,omitempty"`
	Lte *int `json:"lte,omitempty"`
}

type end struct {
	Gte *int `json:"gte,omitempty"`
	Lte *int `json:"lte,omitempty"`
}

//Final Response Structs
type FinalResponse struct {
	Count                   int                     `json:"unique_recipe_count"`
	CountPerRecipeList      []RecipeCountItem       `json:"count_per_recipe"`
	BusiestPostcode         BusiestPostcode         `json:"busiest_postcode"`
	CountPerPostcodeAndTime CountPerPostcodeAndTime `json:"count_per_postcode_and_time"`
	MatchByName             []string                `json:"match_by_name"`
}

type RecipeCountItem struct {
	Recipe string `json:"recipe"`
	Count  int    `json:"count"`
}

type BusiestPostcode struct {
	Postcode      string `json:"postcode"`
	DeliveryCount int    `json:"delivery_count"`
}

type CountPerPostcodeAndTime struct {
	Postcode      string `json:"postcode"`
	From          string `json:"from"`
	To            string `json:"to"`
	DeliveryCount int    `json:"delivery_count"`
}
