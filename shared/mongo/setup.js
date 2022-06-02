db.trip.createIndex({"trip.userid":1,"trip.status":1},{unique:true,partialFilterExpression:{"trip.status":1}})
db.auth.createIndex({"open_id":1},{unique:true})