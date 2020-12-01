package main

// (c) Leonid Titov, 2019, All rights reserved

import (
	"fmt"

	"database/sql"
	_ "github.com/lib/pq"
)

type MyDB struct {
	*sql.DB
}

func (db *MyDB) GoParcelPath(
	absroot uint64,
	parcel_path []string,
) (
	ParcelPathButName string,
	ParcelName string,
	ParcelID uint64,
	err_res error,
) {
	defer func() {
		if err := recover(); err != nil {
			ParcelName = "ERROR: NOT FOUND"
			err_res = fmt.Errorf("GoParcelPath() can't resolve path, PATH NOT FOUND: %v", err)
		}
	}()

	if len(parcel_path) > 1 {
	} else if len(parcel_path) == 1 {
		ParcelName = parcel_path[0]
	} else {
		ParcelID = absroot
		return
	}

	pps := parcel_path[0:len(parcel_path)]

	id := absroot

	var err error

	for i := 0; i < len(pps); i++ {
		id, err = db.OfParentGetChildByName(id, pps[i])
		if err != nil {
			panic(fmt.Errorf("#1 : %v", err))
		}
		if i < len(pps)-1 {
			ParcelPathButName += " /" + pps[i]
		} else {
			ParcelName = pps[i]
		}
	}
	ParcelID = id
	return
}

func (db *MyDB) OfParentGetChildByName(
	id1 uint64,
	name string,
) (
	id2 uint64,
	err error,
) {

	names, ids := db.ListChildrenOfId(id1)

	for n, p := range names {
		if p == name {
			id2 = ids[n]
			return
		}
	}
	err = fmt.Errorf("OfParentGetChildByName(), name '%v' not found", name)
	return
}

func (db *MyDB) ListChildrenOfId(
	id uint64,
) (
	names []string,
	ids []uint64,
) {

	names = make([]string, 0, 20)
	ids = make([]uint64, 0, 20)

	reslist, err := db.Query(`
		SELECT	parcels.name	AS name,
			parcels.barcode_id	AS id
		FROM parcels, parcels_data
		WHERE		parcels.barcode_id = parcels_data.points_to_barcode_id
			AND	parcels_data.parent_parcel_barcode_id = $1
		;
	`, id)
	if err != nil {
		panic(fmt.Errorf("#2 ListChildrenOfId() failed (1): %v", err))
	}

	for {
		if !reslist.Next() {
			break
		}
		var name string
		err = reslist.Scan(&name, &id)
		if err != nil {
			panic(fmt.Errorf("#3 ListChildrenOfId() failed (2): %v", err))
		}
		names = append(names, name)
		ids = append(ids, id)
	}
	return
}

func (db *MyDB) TryToIdentifyObjectType(
	id uint64,
) (
	obj_type string,
	err_res error,
) {
	defer func() {
		if err := recover(); err != nil {
			err_res = fmt.Errorf("TryToIdentifyObjectType() failed: %v", err)
		}
	}()

	err := db.QueryRow(`
		SELECT points_to_table AS obj_type
		FROM barcode_primary
		WHERE barcode_id=$1
		;
	`, id).Scan(&obj_type)
	if err != nil {
		panic(fmt.Errorf("#4 : %v", err))
	}

	switch obj_type {
	case "parcels":
		obj_type = "<b>PARCEL</b>"
		p := struct {
			p1_desc sql.NullString
			p2_desc sql.NullString
			t_name  sql.NullString
		}{}
		err := db.QueryRow(`
			SELECT description
			FROM parcels
			WHERE barcode_id=$1
			;
		`, id).Scan(&p.p1_desc)
		if err != nil {
			err = nil
		}
		err = db.QueryRow(`
			SELECT	p2.description,
				p2.name
			FROM parcels AS p1, parcels AS p2
			WHERE		p1.barcode_id=$1
				AND	p2.barcode_id=p1.type
				AND	p1.type IS NOT NULL
			;
		`, id).Scan(&p.p2_desc, &p.t_name)
		if err != nil {
			err = nil
		}

		if p.p1_desc.Valid && p.p1_desc.String != "" {
			obj_type += ", " + p.p1_desc.String
		}
		if p.t_name.Valid && p.t_name.String != "" {
			obj_type += " <b>OF TYPE:</b> " + p.t_name.String
		}
		if p.p2_desc.Valid && p.p2_desc.String != "" {
			obj_type += ", " + p.p2_desc.String
		}

	case "documents_a4_cas":
		obj_type = "Document A4 CAS"
	case "employees":
		obj_type = "Employee"
	case "locations":
		obj_type = "Location"
	case "bulk_store_1":
		obj_type = "CAD Part, table " + obj_type
	default:
		obj_type = "(database object type unidentified)"
	}

	return
}

func (db *MyDB) TryToIdentifyObjectTypeV2(
	id uint64,
) (
	obj_type string,
	obj_descr string,
	has_children bool,
	err_res error,
) {
	defer func() {
		if err := recover(); err != nil {
			err_res = fmt.Errorf("TryToIdentifyObjectTypeV2() failed: %v", err)
		}
	}()

	err := db.QueryRow(`
		SELECT points_to_table AS obj_type
		FROM barcode_primary
		WHERE barcode_id=$1
		;
	`, id).Scan(&obj_type)
	if err != nil {
		panic(fmt.Errorf("#5 : %v", err))
	}

	switch obj_type {
	case "parcels":
		obj_type = ""

		var ss1 sql.NullString

		err := db.QueryRow(`
			SELECT description
			FROM parcels
			WHERE barcode_id=$1
			;
		`, id).Scan(&ss1)
		if err != nil {
			err = nil
		}
		obj_descr = ss1.String

		err = db.QueryRow(`
			SELECT	p2.name
			FROM parcels AS p1, parcels AS p2
			WHERE		p1.barcode_id=$1
				AND	p2.barcode_id=p1.type
				AND	p1.type IS NOT NULL
			;
		`, id).Scan(&ss1)
		if err == nil {
			obj_type = ss1.String
		}

		var h int
		err = db.QueryRow(`
			SELECT COUNT(parcels.barcode_id) AS h
			FROM parcels, parcels_data
			WHERE		parcels.barcode_id = parcels_data.points_to_barcode_id
				AND	parcels_data.parent_parcel_barcode_id = $1
			;
		`, id).Scan(&h)
		if err != nil {
			panic(fmt.Errorf("#6 : %v", err))
		}
		if h > 0 {
			has_children = true
		} else {
			has_children = false
		}

	case "documents_a4_cas":
		obj_type = "Document A4 CAS"
	case "employees":
		obj_type = "Employee"
	case "locations":
		obj_type = "Location"
	case "bulk_store_1":
		obj_type = "CAD Part, table " + obj_type
	default:
		obj_type = "(database object type unidentified)"
	}

	return
}

func (db *MyDB) TraceParentalPath(
	id uint64,
) (
	parcel_path_reversed []string,
	err_res error,
) {
	defer func() {
		if err := recover(); err != nil {
			err_res = fmt.Errorf("TraceParentalPath() failed: %v", err)
		}
	}()

	parcel_path_reversed = make([]string, 0, 10)

L1:
	name, id, found, err := db.GetSingleParent(id)
	if err != nil {
		panic(fmt.Errorf("#7 : %v", err))
	}
	if !found {
		goto L1_e
	}
	parcel_path_reversed = append(parcel_path_reversed, name)
	goto L1
L1_e:

	return
}

func (db *MyDB) GetSingleParent(
	id uint64,
) (
	name string,
	id_par uint64,
	found bool,
	err_res error,
) {
	defer func() {
		if err := recover(); err != nil {
			err_res = fmt.Errorf("GetSingleParent(): %v", err)
		}
	}()

	err := db.QueryRow(`
		SELECT	parcels.name	AS name,
			parcels.barcode_id	AS id
		FROM parcels, parcels_data
		WHERE		parcels.barcode_id = parcels_data.parent_parcel_barcode_id
			AND	parcels_data.points_to_barcode_id = $1
		;
	`, id).Scan(&name, &id_par)
	if err == nil {
		found = true
	} else if err == sql.ErrNoRows {
		found = false
	} else {
		panic(err)
	}
	return
}

func Transact(db *sql.DB, txFunc func(*sql.Tx)) (err_res error) {
	defer func() {
		if err := recover(); err != nil {
			err_res = fmt.Errorf("Transact failed: %v", err)
		}
	}()
	tx, err := db.Begin()
	if err != nil {
		panic(fmt.Errorf("#12 : %v", err))
	}
	err = func() (err_res error) {
		defer func() {
			if err := recover(); err != nil {
				err_res = fmt.Errorf("Rollback because of: %v", err)
				tx.Rollback()
			}
		}()
		txFunc(tx)
		err = tx.Commit()
		if err != nil {
			panic(fmt.Errorf("#13 : %v", err))
		}
		return nil
	}()
	if err != nil {
		panic(fmt.Errorf("#14 : %v", err))
	}
	return nil
}

func (db *MyDB) Transact2(txFunc func(*sql.Tx)) (err_res error) {
	defer func() {
		if err := recover(); err != nil {
			err_res = fmt.Errorf("Transact failed: %v", err)
		}
	}()
	tx, err := db.Begin()
	if err != nil {
		panic(fmt.Errorf("#15 : %v", err))
	}
	err = func() (err_res error) {
		defer func() {
			if err := recover(); err != nil {
				err_res = fmt.Errorf("Rollback because of: %v", err)
				tx.Rollback()
			}
		}()
		txFunc(tx)
		err = tx.Commit()
		if err != nil {
			panic(fmt.Errorf("#16 : %v", err))
		}
		return nil
	}()
	if err != nil {
		panic(fmt.Errorf("#17 : %v", err))
	}
	return nil
}
