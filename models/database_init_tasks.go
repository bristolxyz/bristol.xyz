package models

var databaseInitTasks []func()

// RunDatabaseInitTasks is used to run all database initialisation tasks.
func RunDatabaseInitTasks() {
	for _, v := range databaseInitTasks {
		v()
	}
	databaseInitTasks = []func(){}
}
