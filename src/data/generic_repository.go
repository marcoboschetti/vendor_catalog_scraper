package data

func InsertNew[K any](entity K) error {
	return nil
	// _, err := pgConnection.Model(&entity).Insert()
	// return err
}

func InsertMany[K any](entities []K) error {
	return nil

	// _, err := pgConnection.Model(&entities).Insert()
	// return err
}

func Delete[K any](entity K) error {
	return nil
	// _, err := pgConnection.Model(&entity).WherePK().Delete()
	// return err
}

func Update[K any](entity K) error {
	return nil
	// _, err := pgConnection.Model(&entity).WherePK().Update()
	// return err
}

func GetByID[K any](id string) (*K, error) {
	return nil, nil
	// var entity K

	// err := pgConnection.Model(&entity).
	// 	Where("id = ?", id).
	// 	Select()

	// return &entity, err
}

func GetAll[K any]() ([]K, error) {
	return nil, nil
	// var entities []K
	// err := pgConnection.Model(&entities).Select()

	// return entities, err
}
