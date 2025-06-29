package repository

import (
	"database/sql"
	"fmt"
	"payment/internal/errors"
	"payment/internal/model"

	"github.com/lib/pq"
)

type PackageRepository struct {
	db *sql.DB
}

func NewPackageRepository(db *sql.DB) *PackageRepository {
	return &PackageRepository{db: db}
}

func (r *PackageRepository) FindAllPackage() ([]model.Packages, error) {

	var packages []model.Packages

	rows, err := r.db.Query(`
        SELECT id, name, price, description, benefits 
        FROM packages
        ORDER BY id
    `)
	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var p model.Packages
		err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Description, pq.Array(&p.Benefits))
		if err != nil {
			return nil, fmt.Errorf("scan error: %w", err)
		}
		packages = append(packages, p)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return packages, nil
}

func (r *PackageRepository) FindPackageById(id int) (*model.Packages, error) {
	var p model.Packages

	err := r.db.QueryRow(`
		SELECT id, name, price, description, benefits
		FROM packages
		WHERE id = $1
	`, id).Scan(&p.ID, &p.Name, &p.Price, &p.Description, pq.Array(&p.Benefits))

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.ErrPackageNotFound
		}
		return nil, fmt.Errorf("failed to find package by id: %w", err)
	}

	return &p, nil
}
