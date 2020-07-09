package database

import "fmt"

func DeleteById(id string) error {
	stmt, err := Conn().Prepare("DELETE FROM customers WHERE id=$1")
	if err != nil {
		return fmt.Errorf("can't prepare delete statement: %w", err)
	}

	if _, err := stmt.Exec(id); err != nil {
		return fmt.Errorf("can't execute delete: %w", err)
	}

	return nil
}
