package database

const schema = `
PRAGMA foreign_keys = ON;

CREATE TABLE IF NOT EXISTS users (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	username TEXT NOT NULL UNIQUE,
	password_hash TEXT NOT NULL,
	display_name TEXT NOT NULL,
	created_at TEXT NOT NULL,
	updated_at TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS devices (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	code TEXT NOT NULL UNIQUE,
	name TEXT NOT NULL,
	category TEXT NOT NULL,
	model TEXT NOT NULL,
	serial_number TEXT NOT NULL,
	manufacturer TEXT NOT NULL,
	location TEXT NOT NULL,
	install_date TEXT NOT NULL,
	warranty_end TEXT NOT NULL,
	owner TEXT NOT NULL,
	status TEXT NOT NULL,
	remark TEXT NOT NULL,
	created_at TEXT NOT NULL,
	updated_at TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS inspection_plans (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT NOT NULL,
	scope_type TEXT NOT NULL,
	device_id INTEGER,
	category TEXT NOT NULL DEFAULT '',
	period TEXT NOT NULL,
	start_date TEXT NOT NULL,
	end_date TEXT NOT NULL,
	owner TEXT NOT NULL,
	status TEXT NOT NULL,
	created_at TEXT NOT NULL,
	updated_at TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS plan_check_items (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	plan_id INTEGER NOT NULL,
	name TEXT NOT NULL,
	item_type TEXT NOT NULL,
	standard_range TEXT NOT NULL,
	required INTEGER NOT NULL DEFAULT 0,
	description TEXT NOT NULL,
	sort_order INTEGER NOT NULL DEFAULT 0,
	FOREIGN KEY (plan_id) REFERENCES inspection_plans(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS inspection_tasks (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	plan_id INTEGER NOT NULL,
	device_id INTEGER NOT NULL,
	planned_at TEXT NOT NULL,
	status TEXT NOT NULL,
	actual_executor_id INTEGER,
	executed_at TEXT,
	results TEXT NOT NULL DEFAULT '[]',
	remark TEXT NOT NULL DEFAULT '',
	created_at TEXT NOT NULL,
	updated_at TEXT NOT NULL,
	FOREIGN KEY (plan_id) REFERENCES inspection_plans(id) ON DELETE CASCADE,
	FOREIGN KEY (device_id) REFERENCES devices(id) ON DELETE CASCADE
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_tasks_plan_device_date
	ON inspection_tasks(plan_id, device_id, planned_at);

CREATE TABLE IF NOT EXISTS inspection_records (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	task_id INTEGER NOT NULL,
	device_id INTEGER NOT NULL,
	plan_id INTEGER NOT NULL,
	executor_id INTEGER NOT NULL,
	executed_at TEXT NOT NULL,
	status TEXT NOT NULL,
	result_summary TEXT NOT NULL,
	created_at TEXT NOT NULL,
	FOREIGN KEY (task_id) REFERENCES inspection_tasks(id) ON DELETE CASCADE,
	FOREIGN KEY (device_id) REFERENCES devices(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS anomalies (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	device_id INTEGER NOT NULL,
	task_id INTEGER,
	discoverer_id INTEGER NOT NULL,
	discovered_at TEXT NOT NULL,
	severity TEXT NOT NULL,
	description TEXT NOT NULL,
	status TEXT NOT NULL,
	assignee_id INTEGER,
	progress TEXT NOT NULL DEFAULT '',
	cause_analysis TEXT NOT NULL DEFAULT '',
	close_note TEXT NOT NULL DEFAULT '',
	verification_result TEXT NOT NULL DEFAULT '',
	closed_at TEXT,
	created_at TEXT NOT NULL,
	updated_at TEXT NOT NULL,
	FOREIGN KEY (device_id) REFERENCES devices(id) ON DELETE CASCADE,
	FOREIGN KEY (task_id) REFERENCES inspection_tasks(id) ON DELETE SET NULL
);

CREATE TABLE IF NOT EXISTS repairs (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	device_id INTEGER NOT NULL,
	anomaly_id INTEGER,
	repair_type TEXT NOT NULL,
	content TEXT NOT NULL,
	vendor TEXT NOT NULL,
	started_at TEXT NOT NULL,
	ended_at TEXT NOT NULL,
	cost REAL NOT NULL DEFAULT 0,
	result TEXT NOT NULL,
	remark TEXT NOT NULL,
	created_by INTEGER NOT NULL,
	created_at TEXT NOT NULL,
	updated_at TEXT NOT NULL,
	FOREIGN KEY (device_id) REFERENCES devices(id) ON DELETE CASCADE,
	FOREIGN KEY (anomaly_id) REFERENCES anomalies(id) ON DELETE SET NULL
);
`
