package oxilib

import (
	"time"
)

const sqlDeleteItem = `UPDATE items SET deleted=1, updated=? WHERE item_id=? `
const sqlDeleteItemField = `UPDATE fields SET deleted=1, updated=? WHERE item_id=? `

func (sdb *storageDB) dbDeleteItem(itemID int64) (err error) {
	if sdb.sTX == nil {
		return formError(BSERR00003DbTransactionFailed, "dbDeleteItem")
	}

	updateTime := prepareTimeForDb(time.Now())

	stmtItem, err := sdb.sTX.Prepare(sqlDeleteItem)
	if err != nil {
		return err
	}

	_, err = stmtItem.Exec(updateTime, itemID)
	if err != nil {
		return err
	}

	errClose := stmtItem.Close()
	if errClose != nil {
		return formError(BSERR00016DbDeleteFailed, errClose.Error())
	}

	stmtItemFields, err := sdb.sTX.Prepare(sqlDeleteItemField)
	if err != nil {
		return err
	}

	_, err = stmtItemFields.Exec(updateTime, itemID)
	if err != nil {
		return err
	}

	errClose = stmtItemFields.Close()
	if errClose != nil {
		return formError(BSERR00016DbDeleteFailed, errClose.Error())
	}
	return nil
}

const sqlInsertItem = `
	INSERT 
		INTO items (name,icon,created,updated,deleted) 
		VALUES (?,?,?,?, 0)
`

func (sdb *storageDB) dbInsertItem(itemName string, itemIcon string) (itemID int64, err error) {
	if sdb.sTX == nil {
		return 0, formError(BSERR00003DbTransactionFailed, "dbInsertItem")
	}

	creationTime := prepareTimeForDb(time.Now())

	stmt, err := sdb.sTX.Prepare(sqlInsertItem)
	if err != nil {
		return 0, formError(BSERR00006DbInsertFailed, err.Error(), "dbInsertItem")
	}
	res, errStmt := stmt.Exec(itemName,
		itemIcon,
		creationTime,
		creationTime)

	if errStmt != nil {
		return 0, formError(BSERR00006DbInsertFailed, errStmt.Error(), "dbInsertItem")
	}
	itemID, err = res.LastInsertId()
	if err != nil {
		return 0, formError(BSERR00006DbInsertFailed, err.Error(), "dbInsertItem")
	}
	errClose := stmt.Close()
	if errClose != nil {
		return 0, formError(BSERR00006DbInsertFailed, errClose.Error(), "dbInsertItem")
	}

	return itemID, nil
}

// List all non-deleted items
const sqlListItems = `
	SELECT item_id, name, icon, created, updated, deleted
		FROM items 
		WHERE deleted='0'
`

// List all non-deleted items
const sqlListItemsWithDeleted = `
	SELECT item_id, name, icon, created, updated, deleted
		FROM items
`

func (sdb *storageDB) dbSelectAllItems(returnDeleted bool) (items []OxiItem, err error) {
	var sqlList string
	if returnDeleted {
		sqlList = sqlListItemsWithDeleted
	} else {
		sqlList = sqlListItems
	}
	rows, err := sdb.sDB.Query(sqlList)
	if err != nil {
		return items, formError(BSERR00014ItemsReadFailed, err.Error())
	}
	defer func() {
		errClose := rows.Close()
		if errClose != nil {
			err = formError(BSERR00014ItemsReadFailed, err.Error(), errClose.Error())
		}

	}()

	var item OxiItem

	for rows.Next() {
		err = rows.Scan(&item.ID,
			&item.Name,
			&item.Icon,
			&item.Created,
			&item.Updated,
			&item.Deleted)
		if err != nil {
			return items, err
		}

		if item.Deleted && !returnDeleted {
			continue
		}

		items = append(items, item)
	}
	return items, nil
}

const sqlUpdateItemName = `UPDATE items SET name=?, updated=? WHERE item_id=? `

func (sdb *storageDB) dbUpdateItemName(itemID int64, newName string) (err error) {
	if sdb.sTX == nil {
		return formError(BSERR00003DbTransactionFailed, "dbUpdateItemName")
	}

	stmt, err := sdb.sTX.Prepare(sqlUpdateItemName)
	if err != nil {
		return err
	}

	updateTime := prepareTimeForDb(time.Now())

	_, err = stmt.Exec(newName, updateTime, itemID)
	if err != nil {
		return err
	}

	errClose := stmt.Close()
	if errClose != nil {
		return formError(BSERR00018DbItemNameUpdateFailed, errClose.Error())
	}
	return nil
}

// Get item by id
const sqlGetItemById = `
	SELECT item_id, name, icon, created, updated, deleted
		FROM items 
		WHERE item_id=? and deleted=0
`

// Get item by id including deleted
const sqlGetItemByIdWithDeleted = `
	SELECT item_id, name, icon, created, updated, deleted
		FROM items 
		WHERE item_id=?
`

func (sdb *storageDB) dbGetItemById(itemId int64, withDeleted bool) (item OxiItem, err error) {
	var sqlRequest string
	if withDeleted {
		sqlRequest = sqlGetItemByIdWithDeleted
	} else {
		sqlRequest = sqlGetItemById
	}
	stmt, err := sdb.sDB.Prepare(sqlRequest)
	if err != nil {
		return item, formError(BSERR00014ItemsReadFailed, err.Error())
	}
	rows, err := stmt.Query(itemId)
	if err != nil {
		return item, formError(BSERR00014ItemsReadFailed, err.Error())
	}

	defer func() {
		errClose := rows.Close()
		if errClose != nil {
			err = formError(BSERR00014ItemsReadFailed, err.Error(), errClose.Error())
		}
		errClose = stmt.Close()
		if errClose != nil {
			err = formError(BSERR00014ItemsReadFailed, err.Error(), errClose.Error())
		}
	}()

	if rows.Next() {

		err = rows.Scan(&item.ID,
			&item.Name,
			&item.Icon,
			&item.Created,
			&item.Updated,
			&item.Deleted)
		if err != nil {
			return item, err
		}
		return item, nil
	}
	return item, formError(BSERR00019ItemNotFound)
}

const sqlUpdateItemIcon = `UPDATE items SET icon=?, updated=? WHERE item_id=? `

func (sdb *storageDB) dbUpdateItemIcon(itemID int64, newIcon string) (err error) {
	if sdb.sTX == nil {
		return formError(BSERR00003DbTransactionFailed, "dbUpdateItemIcon")
	}

	stmt, err := sdb.sTX.Prepare(sqlUpdateItemIcon)
	if err != nil {
		return err
	}

	updateTime := prepareTimeForDb(time.Now())

	_, err = stmt.Exec(newIcon, updateTime, itemID)
	if err != nil {
		return err
	}

	errClose := stmt.Close()
	if errClose != nil {
		return formError(BSERR00026DbItemIconUpdateFailed, errClose.Error())
	}
	return nil
}
