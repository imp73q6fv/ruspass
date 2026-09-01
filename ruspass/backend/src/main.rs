mod crypto;
mod database;
mod models;

use database::DatabaseService;
use crypto::CryptoService;
use models::PasswordEntry;

fn main() {
    println!("=== RusPass Password Manager ===");
    println!("Backend initialized with ChaCha20-Poly1305 encryption and Argon2id key derivation");
    
    // For demonstration, create a test database
    let db_path = "ruspass.db";
    let db = DatabaseService::new(db_path).expect("Failed to create database");
    
    println!("\nDatabase created at: {}", db_path);
    println!("Ready for frontend integration via IPC or API");
    
    // Demo: Initialize database with master password
    println!("\n--- Demo Mode ---");
    demo_operations(&db);
}

fn demo_operations(db: &DatabaseService) {
    let master_password = "demo_master_password_123";
    
    // Initialize database
    if let Err(e) = db.init_database(master_password) {
        println!("Database already initialized or error: {:?}", e);
    } else {
        println!("Database initialized with master password");
    }
    
    // Get config (salt and nonce)
    match db.get_config() {
        Ok(Some((salt, nonce))) => {
            println!("Retrieved database config (salt: {} bytes, nonce: {} bytes)", 
                     salt.len(), nonce.len());
            
            // Derive key from master password
            match CryptoService::derive_key(master_password, &salt) {
                Ok(key) => {
                    println!("Key derived successfully ({} bytes)", key.len());
                    
                    // Create a sample entry
                    let entry = PasswordEntry {
                        id: None,
                        title: "GitHub".to_string(),
                        username: "user@example.com".to_string(),
                        password: "secure_password_123".to_string(),
                        url: Some("https://github.com".to_string()),
                        notes: Some("My GitHub account".to_string()),
                        category: "Work".to_string(),
                        created_at: 0,
                        updated_at: 0,
                    };
                    
                    // Add entry
                    match db.add_entry(&entry, &key, &nonce) {
                        Ok(id) => println!("Added entry with ID: {}", id),
                        Err(e) => println!("Failed to add entry: {:?}", e),
                    }
                    
                    // Get all entries
                    match db.get_all_entries(&key, &nonce) {
                        Ok(entries) => {
                            println!("\nRetrieved {} entries:", entries.len());
                            for e in &entries {
                                println!("  - {} ({})", e.title, e.username);
                            }
                        },
                        Err(e) => println!("Failed to get entries: {:?}", e),
                    }
                    
                    // Search entries
                    match db.search_entries("git", &key, &nonce) {
                        Ok(results) => {
                            println!("\nSearch results for 'git': {} entries", results.len());
                        },
                        Err(e) => println!("Search failed: {:?}", e),
                    }
                },
                Err(e) => println!("Key derivation failed: {:?}", e),
            }
        },
        Ok(None) => println!("Database not initialized"),
        Err(e) => println!("Failed to get config: {:?}", e),
    }
    
    println!("\n--- Demo Complete ---");
    println!("The backend is ready for integration with the Go frontend.");
}
