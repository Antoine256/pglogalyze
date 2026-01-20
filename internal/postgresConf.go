package internal

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
)

func GetPathByDatabaseConn(params []CmdParam) string {
	//----------------------- GETTING DATABASE INFO ------------------
	//définis dans le fichier de conf then => (ici env pour exemple)
	// conf complète à vérifier avec fonction vérif conf qui prend tout ou rien ?

	var conn *pgx.Conn
	var err error
	var connStr string

	if HasParams("db", params) {

		connStr = GetParams("db", params, 0)
		conn, err = pgx.Connect(context.Background(), connStr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
			return ""
		}
		defer conn.Close(context.Background())

	} else if HasParams("c", params) /* OR by default ??*/ {

		parameterFilePath := GetParams("c", params, 0)
		connStr = fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", "", "", "", "", "")
		conn, err = pgx.Connect(context.Background(), connStr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to connect to database %v\n", err)
			fmt.Fprintf(os.Stderr, "Connection with parameter didn't work \n")
			fmt.Fprintln(os.Stderr, "Parameter file "+parameterFilePath+" is incorrect \n")
			return ""
		}
		defer conn.Close(context.Background())

	} else {

		conn, err = pgx.Connect(context.Background(), "postgresql://postgres:postgres@localhost:5432/postgres")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
			os.Exit(1)
		}
		defer conn.Close(context.Background())

	}

	filePath, err := GetCurrentLogAbsoluteFilePath(conn, context.Background())
	if err != nil || !PathExists(filePath) {
		fmt.Fprintf(os.Stderr, "Unable to find log file: %v\n", err)
		os.Exit(1)
	}

	return filePath
}

func verifConf(databaseConn *pgx.Conn, ctx context.Context) bool {
	var loggingcol string
	err := databaseConn.QueryRow(ctx, "show logging_collector;").Scan(&loggingcol)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}
	return loggingcol == "on"
}

func GetCurrentLogAbsoluteFilePath(databaseConn *pgx.Conn, ctx context.Context) (string, error) {

	isVerif := verifConf(databaseConn, ctx)

	if !isVerif {
		return "", fmt.Errorf("logging_collector isn't set on")
	}

	filePath := getVariable(databaseConn, ctx, "SELECT  pg_current_logfile();")
	dataDir := getVariable(databaseConn, ctx, "show  data_directory;")

	res := dataDir + "/" + filePath

	return res, nil
}

func getVariable(databaseConn *pgx.Conn, ctx context.Context, request string) string {
	var value string
	err := databaseConn.QueryRow(ctx, request).Scan(&value)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}
	return value
}
