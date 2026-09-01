use rusqlite::{Connection, params};
use std::time::{SystemTime, UNIX_EPOCH};
use super::models::PasswordEntry;
use super::crypto::{CryptoService, CryptoError};

// Helper function to convert CryptoError to rusqlite::Error
fn crypto_to_sqlite_error(e: CryptoError) -> rusqlite::Error {
    e.into()
}

pub struct DatabaseService {
    conn: Connection,
}

impl DatabaseService {
    pub fn new(path: &str) -> Result<Self, rusqlite::Error> {
        let conn = Connection::open(path)?;
        
        // Create tables
        conn.execute(
            "CREATE TABLE IF NOT EXISTS config (
                id INTEGER PRIMARY KEY CHECK (id = 1),
                key_salt BLOB NOT NULL,
                nonce BLOB NOT NULL,
                created_at INTEGER NOT NULL
            )",
            [],
        )?;
        
        conn.execute(
            "CREATE TABLE IF NOT EXISTS entries (
                id INTEGER PRIMARY KEY AUTOINCREMENT,
                title TEXT NOT NULL,
                username TEXT NOT NULL,
                password BLOB NOT NULL,
                url TEXT,
                notes TEXT,
                category TEXT NOT NULL DEFAULT 'General',
                created_at INTEGER NOT NULL,
                updated_at INTEGER NOT NULL
            )",
            [],
        )?;
        
        conn.execute(
            "CREATE INDEX IF NOT EXISTS idx_entries_category ON entries(category)",
            [],
        )?;
        
        Ok(DatabaseService { conn })
    }

    pub fn init_database(&self, _master_password: &str) -> Result<(), CryptoError> {
        let salt = CryptoService::generate_salt();
        let nonce = CryptoService::generate_nonce();
        let now = SystemTime::now()
            .duration_since(UNIX_EPOCH)
            .unwrap()
            .as_secs() as i64;

        self.conn.execute(
            "INSERT OR REPLACE INTO config (id, key_salt, nonce, created_at) VALUES (1, ?1, ?2, ?3)",
            params![&salt, &nonce, now],
        ).map_err(|e| CryptoError::Encryption(e.to_string()))?;

        Ok(())
    }

    pub fn get_config(&self) -> Result<Option<(Vec<u8>, Vec<u8>)>, rusqlite::Error> {
        let mut stmt = self.conn.prepare("SELECT key_salt, nonce FROM config WHERE id = 1")?;
        let result = stmt.query_row([], |row| {
            let salt: Vec<u8> = row.get(0)?;
            let nonce: Vec<u8> = row.get(1)?;
            Ok((salt, nonce))
        });

        match result {
            Ok(config) => Ok(Some(config)),
            Err(rusqlite::Error::QueryReturnedNoRows) => Ok(None),
            Err(e) => Err(e),
        }
    }

    pub fn add_entry(&self, entry: &PasswordEntry, key: &[u8], nonce: &[u8]) -> Result<i64, CryptoError> {
        let now = SystemTime::now()
            .duration_since(UNIX_EPOCH)
            .unwrap()
            .as_secs() as i64;

        // Encrypt the password
        let password_bytes = entry.password.as_bytes();
        let encrypted = CryptoService::encrypt(password_bytes, key, nonce)?;
        let encrypted_b64 = CryptoService::encode_base64(&encrypted);

        let mut stmt = self.conn.prepare(
            "INSERT INTO entries (title, username, password, url, notes, category, created_at, updated_at) 
             VALUES (?1, ?2, ?3, ?4, ?5, ?6, ?7, ?8)"
        ).map_err(|e| CryptoError::Encryption(e.to_string()))?;

        stmt.execute(params![
            entry.title,
            entry.username,
            encrypted_b64,
            entry.url,
            entry.notes,
            entry.category,
            now,
            now
        ]).map_err(|e| CryptoError::Encryption(e.to_string()))?;

        Ok(self.conn.last_insert_rowid())
    }

    pub fn get_all_entries(&self, key: &[u8], nonce: &[u8]) -> Result<Vec<PasswordEntry>, CryptoError> {
        let mut stmt = self.conn.prepare(
            "SELECT id, title, username, password, url, notes, category, created_at, updated_at 
             FROM entries ORDER BY title ASC"
        ).map_err(|e| CryptoError::Decryption(e.to_string()))?;

        let entries = stmt.query_map([], |row| {
            let id: i64 = row.get(0)?;
            let title: String = row.get(1)?;
            let username: String = row.get(2)?;
            let encrypted_password: String = row.get(3)?;
            let url: Option<String> = row.get(4)?;
            let notes: Option<String> = row.get(5)?;
            let category: String = row.get(6)?;
            let created_at: i64 = row.get(7)?;
            let updated_at: i64 = row.get(8)?;

            // Decrypt the password
            let encrypted_bytes = CryptoService::decode_base64(&encrypted_password)
                .map_err(crypto_to_sqlite_error)?;
            let decrypted_bytes = CryptoService::decrypt(&encrypted_bytes, key, nonce)
                .map_err(crypto_to_sqlite_error)?;
            let password = String::from_utf8(decrypted_bytes)
                .map_err(|e| crypto_to_sqlite_error(CryptoError::Decryption(e.to_string())))?;

            Ok(PasswordEntry {
                id: Some(id),
                title,
                username,
                password,
                url,
                notes,
                category,
                created_at,
                updated_at,
            })
        }).map_err(|e| CryptoError::Decryption(e.to_string()))?;

        let mut result = Vec::new();
        for entry in entries {
            match entry {
                Ok(e) => result.push(e),
                Err(e) => return Err(CryptoError::Decryption(e.to_string())),
            }
        }

        Ok(result)
    }

    pub fn get_entry_by_id(&self, id: i64, key: &[u8], nonce: &[u8]) -> Result<PasswordEntry, CryptoError> {
        let mut stmt = self.conn.prepare(
            "SELECT id, title, username, password, url, notes, category, created_at, updated_at 
             FROM entries WHERE id = ?1"
        ).map_err(|e| CryptoError::Decryption(e.to_string()))?;

        let entry = stmt.query_row(params![id], |row| {
            let entry_id: i64 = row.get(0)?;
            let title: String = row.get(1)?;
            let username: String = row.get(2)?;
            let encrypted_password: String = row.get(3)?;
            let url: Option<String> = row.get(4)?;
            let notes: Option<String> = row.get(5)?;
            let category: String = row.get(6)?;
            let created_at: i64 = row.get(7)?;
            let updated_at: i64 = row.get(8)?;

            // Decrypt the password
            let encrypted_bytes = CryptoService::decode_base64(&encrypted_password)
                .map_err(crypto_to_sqlite_error)?;
            let decrypted_bytes = CryptoService::decrypt(&encrypted_bytes, key, nonce)
                .map_err(crypto_to_sqlite_error)?;
            let password = String::from_utf8(decrypted_bytes)
                .map_err(|e| crypto_to_sqlite_error(CryptoError::Decryption(e.to_string())))?;

            Ok(PasswordEntry {
                id: Some(entry_id),
                title,
                username,
                password,
                url,
                notes,
                category,
                created_at,
                updated_at,
            })
        }).map_err(|e| CryptoError::Decryption(e.to_string()))?;

        Ok(entry)
    }

    pub fn update_entry(&self, entry: &PasswordEntry, key: &[u8], nonce: &[u8]) -> Result<(), CryptoError> {
        let now = SystemTime::now()
            .duration_since(UNIX_EPOCH)
            .unwrap()
            .as_secs() as i64;

        // Encrypt the password
        let password_bytes = entry.password.as_bytes();
        let encrypted = CryptoService::encrypt(password_bytes, key, nonce)?;
        let encrypted_b64 = CryptoService::encode_base64(&encrypted);

        self.conn.execute(
            "UPDATE entries SET title = ?1, username = ?2, password = ?3, url = ?4, 
             notes = ?5, category = ?6, updated_at = ?7 WHERE id = ?8",
            params![
                entry.title,
                entry.username,
                encrypted_b64,
                entry.url,
                entry.notes,
                entry.category,
                now,
                entry.id
            ],
        ).map_err(|e| CryptoError::Encryption(e.to_string()))?;

        Ok(())
    }

    pub fn delete_entry(&self, id: i64) -> Result<(), rusqlite::Error> {
        self.conn.execute("DELETE FROM entries WHERE id = ?1", params![id])?;
        Ok(())
    }

    pub fn search_entries(&self, query: &str, key: &[u8], nonce: &[u8]) -> Result<Vec<PasswordEntry>, CryptoError> {
        let search_pattern = format!("%{}%", query);
        let mut stmt = self.conn.prepare(
            "SELECT id, title, username, password, url, notes, category, created_at, updated_at 
             FROM entries WHERE title LIKE ?1 OR username LIKE ?1 OR url LIKE ?1 ORDER BY title ASC"
        ).map_err(|e| CryptoError::Decryption(e.to_string()))?;

        let entries = stmt.query_map(params![search_pattern], |row| {
            let entry_id: i64 = row.get(0)?;
            let title: String = row.get(1)?;
            let username: String = row.get(2)?;
            let encrypted_password: String = row.get(3)?;
            let url: Option<String> = row.get(4)?;
            let notes: Option<String> = row.get(5)?;
            let category: String = row.get(6)?;
            let created_at: i64 = row.get(7)?;
            let updated_at: i64 = row.get(8)?;

            // Decrypt the password
            let encrypted_bytes = CryptoService::decode_base64(&encrypted_password)
                .map_err(crypto_to_sqlite_error)?;
            let decrypted_bytes = CryptoService::decrypt(&encrypted_bytes, key, nonce)
                .map_err(crypto_to_sqlite_error)?;
            let password = String::from_utf8(decrypted_bytes)
                .map_err(|e| crypto_to_sqlite_error(CryptoError::Decryption(e.to_string())))?;

            Ok(PasswordEntry {
                id: Some(entry_id),
                title,
                username,
                password,
                url,
                notes,
                category,
                created_at,
                updated_at,
            })
        }).map_err(|e| CryptoError::Decryption(e.to_string()))?;

        let mut result = Vec::new();
        for entry in entries {
            match entry {
                Ok(e) => result.push(e),
                Err(e) => return Err(CryptoError::Decryption(e.to_string())),
            }
        }

        Ok(result)
    }

    pub fn get_categories(&self) -> Result<Vec<String>, rusqlite::Error> {
        let mut stmt = self.conn.prepare("SELECT DISTINCT category FROM entries ORDER BY category")?;
        let categories = stmt.query_map([], |row| row.get::<_, String>(0))?;
        
        let mut result = Vec::new();
        for category in categories {
            result.push(category?);
        }
        
        Ok(result)
    }
}
