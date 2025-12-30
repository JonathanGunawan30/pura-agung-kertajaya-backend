package repository

import "gorm.io/gorm"

type Repository[T any] struct {
	DB *gorm.DB
}

func (r *Repository[T]) Create(db *gorm.DB, entity *T) error {
	return db.Create(entity).Error
}

func (r *Repository[T]) Update(db *gorm.DB, entity *T) error {
	return db.Save(entity).Error
}

func (r *Repository[T]) Delete(db *gorm.DB, entity *T) error {
	return db.Delete(entity).Error
}

func (r *Repository[T]) CountById(db *gorm.DB, id any) (int64, error) {
	var total int64
	err := db.Model(new(T)).Where("id = ?", id).Count(&total).Error
	return total, err
}

func (r *Repository[T]) FindById(db *gorm.DB, entity *T, id any) error {
	return db.Where("id = ?", id).Take(entity).Error
}

func (r *Repository[T]) FindAll(db *gorm.DB, dest *[]T) error {
	return db.Find(dest).Error
}

func (r *Repository[T]) FindBySlug(db *gorm.DB, entity *T, slug string) error {
	return db.Where("slug = ?", slug).Take(entity).Error
}

func (r *Repository[T]) CountBySlug(db *gorm.DB, slug string) (int64, error) {
	var total int64
	err := db.Model(new(T)).Where("slug = ?", slug).Count(&total).Error
	return total, err
}

func (r *Repository[T]) CountBySlugIgnoringID(db *gorm.DB, slug string, id any) (int64, error) {
	var total int64
	err := db.Model(new(T)).
		Where("slug = ? AND id != ?", slug, id).
		Count(&total).Error
	return total, err
}
