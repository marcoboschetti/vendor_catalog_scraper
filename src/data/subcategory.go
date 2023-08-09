package data

func GetAllProductURIsInSubcategory(subcategoryID string) (map[string]interface{}, error) {
	return nil, nil
	// var res []struct {
	// 	URI string
	// }
	// err := pgConnection.Model((*entities.Product)(nil)).
	// 	Column("uri").
	// 	Where("subcategory_id = ? ", subcategoryID).
	// 	Select(&res)

	// if err != nil {
	// 	return nil, err
	// }

	// ans := map[string]interface{}{}
	// for _, r := range res {
	// 	ans[r.URI] = true
	// }

	// return ans, nil
}
