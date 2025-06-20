package main

import (
	"database/sql"
	"errors"
)

// ----------  LUGARES ----------

func CreateLugar(l Lugar) (int64, error) {
	var id int64
	err := DB.QueryRow(
		`INSERT INTO lugar (nombre, descripcion, ubicacion, imagen, horario)
         VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		l.Nombre, l.Descripcion, l.Ubicacion, l.Imagen, l.Horario,
	).Scan(&id)
	return id, err
}

func GetLugar(id int64) (Lugar, error) {
	var l Lugar
	err := DB.QueryRow(
		`SELECT nombre, descripcion, ubicacion, imagen, horario
		 FROM lugar WHERE id = $1`, id,
	).Scan(&l.Nombre, &l.Descripcion, &l.Ubicacion, &l.Imagen, &l.Horario)

	if err == sql.ErrNoRows {
		return l, errors.New("lugar no encontrado")
	}
	return l, err
}

func UpdateLugar(id int64, l Lugar) error {
	_, err := DB.Exec(
		`UPDATE lugar
		 SET nombre=$1, descripcion=$2, ubicacion=$3, imagen=$4, horario=$5
		 WHERE id=$16`,
		l.Nombre, l.Descripcion, l.Ubicacion, l.Imagen, l.Horario, id,
	)
	return err
}

func DeleteLugar(id int64) error {
	_, err := DB.Exec(`DELETE FROM lugar WHERE id = $1`, id)
	return err
}

func ListLugares() ([]Lugar, error) {
	rows, err := DB.Query(`SELECT nombre, descripcion, ubicacion, imagen, horario FROM lugar`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []Lugar
	for rows.Next() {
		var l Lugar
		if err := rows.Scan(&l.Nombre, &l.Descripcion, &l.Ubicacion, &l.Imagen, &l.Horario); err != nil {
			return nil, err
		}
		result = append(result, l)
	}
	return result, rows.Err()
}

// ----------  USUARIOS ----------

func CreateUsuario(u Usuario) (int64, error) {
	var id int64
	err := DB.QueryRow(
		`INSERT INTO usuario (nombre, usuario, password, imagen, ubicacion)
		 VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		u.Nombre, u.Usuario, u.Password, u.Imagen, u.Ubicacion,
	).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, err
}

func GetUsuarioByID(id int64) (Usuario, error) {
	var u Usuario
	err := DB.QueryRow(
		`SELECT nombre, usuario, password, imagen, ubicacion
		 FROM usuario WHERE id = $1`, id,
	).Scan(&u.Nombre, &u.Usuario, &u.Password, &u.Imagen, &u.Ubicacion)
	if err == sql.ErrNoRows {
		return u, errors.New("usuario no encontrado")
	}
	return u, err
}

func GetUsuarioByUsername(username string) (Usuario, error) {
	var u Usuario
	err := DB.QueryRow(
		`SELECT nombre, usuario, password, imagen, ubicacion
		 FROM usuario WHERE usuario = $1`, username,
	).Scan(&u.Nombre, &u.Usuario, &u.Password, &u.Imagen, &u.Ubicacion)
	if err == sql.ErrNoRows {
		return u, errors.New("usuario no encontrado")
	}
	return u, err
}

func UpdateUsuario(id int64, u Usuario) error {
	_, err := DB.Exec(
		`UPDATE usuario
		 SET nombre=$1, usuario=$2, password=$3, imagen=$4, ubicacion=$5
		 WHERE id=$6`,
		u.Nombre, u.Usuario, u.Password, u.Imagen, u.Ubicacion, id,
	)
	return err
}

func DeleteUsuario(id int64) error {
	_, err := DB.Exec(`DELETE FROM usuario WHERE id = $1`, id)
	return err
}
