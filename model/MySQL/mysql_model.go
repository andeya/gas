package MySQLModel

import (
	"github.com/go-gas/Config"
	"github.com/go-gas/SQLBuilder/MySQLBuilder"
	"github.com/go-gas/gas/model"
)

type MySQLModel struct {
	model.Model // extend Model

	b *MySQLBuilder.MySQLBuilder
}

func New(cfg *config.Engine) model.ModelInterface {
	b := MySQLBuilder.New(cfg)
	m := &MySQLModel{}
	m.SetBuilder(b)

	return m
}

func (m *MySQLModel) Insert(s interface{}) (int64, error) {
	return m.Builder().Insert(s)
}

func (m *MySQLModel) MultiInsert(s ...interface{}) ([]int64, error) {
	// m.parseStruct(s)
	if err := m.TransactionStart(); err != nil {
		return nil, err
	}

	allLastID := make([]int64, 0)
	for _, insdata := range s {
		lastid, err := m.Builder().Insert(insdata)
		if err != nil || lastid == 0 {
			m.TransactionRollback()
			return nil, err
		} else {
			allLastID = append(allLastID, lastid)
		}

	}

	if err := m.TransactionCommit(); err != nil {
		return nil, err
	}

	return allLastID, nil

}

func (m *MySQLModel) TransactionStart() error {
	return m.Builder().GetDB().Begin()
}

func (m *MySQLModel) TransactionRollback() error {
	return m.Builder().GetDB().Rollback()
}

func (m *MySQLModel) TransactionCommit() error {
	return m.Builder().GetDB().Commit()
}

// func (m *MySQLModel) TestConn()  {
//     se := "SELECT * from cds where id = ?"

// stmt, err := m.db.Conn.Prepare(se)
//     if err != nil {
//         panic(err.Error())
//     }
//     defer stmt.Close()

//     rows, err := stmt.Query(1)
//     if err != nil {
//         panic(err.Error())
//     }

//     clms, err := rows.Columns()
//     if err != nil {
//         panic(err.Error())
//     } else {
//         fmt.Println(len(clms))
//     }

//     values := make([]sql.RawBytes, len(clms))
//     scanArgs := make([]interface{}, len(values))
//     for i := range values {
//         scanArgs[i] = &values[i]
//     }

//     for rows.Next() {
//         // get RawBytes from data
//         err = rows.Scan(scanArgs...)
//         if err != nil {
//             panic(err.Error()) // proper error handling instead of panic in your app
//         }

//         // Now do something with the data.
//         // Here we just print each column as a string.
//         var value string
//         for i, col := range values {
//             // Here we can check if the value is nil (NULL value)
//             if col == nil {
//                 value = "NULL"
//             } else {
//                 value = string(col)
//             }
//             fmt.Println(clms[i], ": ", value)
//         }
//         fmt.Println("-----------------------------------")
//     }

//     if err = rows.Err(); err != nil {
//         panic(err.Error()) // proper error handling instead of panic in your app
//     }
// }
