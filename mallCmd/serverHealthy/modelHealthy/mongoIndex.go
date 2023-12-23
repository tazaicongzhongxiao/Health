package modelHealthy

import mongo "MyTestMall/mallBase/basics/pkg/mongo"

func HealthyIndex() map[string]interface{} {
	var req = make(map[string]interface{})
	req["body_param"] = []mongo.IndexData{
		{Keys: []mongo.IndexKey{{"age", "1"}}},
	}
	req["food"] = []mongo.IndexData{
		{Keys: []mongo.IndexKey{{"heat", "1"}}},
	}
	req["sports"] = []mongo.IndexData{
		{Keys: []mongo.IndexKey{{"name", "1"}}},
	}
	req["pic"] = []mongo.IndexData{
		{Keys: []mongo.IndexKey{{"name", "1"}}},
	}
	return req
}
