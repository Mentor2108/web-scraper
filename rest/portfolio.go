package rest

// func RetrievePortfolio(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
// 	ctx := r.Context()
// 	log := util.GetGlobalLogger(ctx)

// 	rows, err := data.DB.Query(ctx, "SELECT id, title, description, url FROM portfolio_items;")
// 	if err != nil {
// 		http.Error(w, "Failed to fetch items", http.StatusInternalServerError)
// 		return
// 	}
// 	defer rows.Close()

// 	// var items []PortfolioItem
// 	var items []map[string]interface{}
// 	for rows.Next() {
// 		item := make(map[string]interface{})
// 		var id, title, description, url string
// 		if err := rows.Scan(&id, &title, &description, &url); err != nil {
// 			log.Fatal(err)
// 		}
// 		item["id"] = id
// 		item["title"] = title
// 		item["description"] = description
// 		item["url"] = url
// 		items = append(items, item)
// 	}

// 	json.NewEncoder(w).Encode(items)
// }
