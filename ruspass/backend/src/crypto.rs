use argon2::{Argon2, PasswordHasher, password_hash::SaltString};
use chacha20poly1305::{ChaCha20Poly1305, Key, Nonce, AeadInPlace, KeyInit};
use rand::RngCore;
use base64::{Engine as _, engine::general_purpose::STANDARD as BASE64};
use thiserror::Error;

#[derive(Error, Debug)]
pub enum CryptoError {
    #[error("Encryption failed: {0}")]
    Encryption(String),
    #[error("Decryption failed: {0}")]
    Decryption(String),
    #[error("Key derivation failed: {0}")]
    KeyDerivation(String),
}

pub struct CryptoService;

impl CryptoService {
    /// Derive a key from the master password using Argon2id
    pub fn derive_key(password: &str, salt: &[u8]) -> Result<Vec<u8>, CryptoError> {
        let argon2 = Argon2::default();
        let salt_str = SaltString::from_b64(&BASE64.encode(salt))
            .map_err(|e| CryptoError::KeyDerivation(e.to_string()))?;
        
        let hash = argon2.hash_password(password.as_bytes(), &salt_str)
            .map_err(|e| CryptoError::KeyDerivation(e.to_string()))?;
        
        Ok(hash.hash.unwrap().as_bytes().to_vec())
    }

    /// Generate a random salt for key derivation
    pub fn generate_salt() -> Vec<u8> {
        let mut salt = vec![0u8; 16];
        rand::thread_rng().fill_bytes(&mut salt);
        salt
    }

    /// Generate a random nonce for encryption
    pub fn generate_nonce() -> Vec<u8> {
        let mut nonce = vec![0u8; 12];
        rand::thread_rng().fill_bytes(&mut nonce);
        nonce
    }

    /// Encrypt data using ChaCha20-Poly1305
    pub fn encrypt(data: &[u8], key: &[u8], nonce: &[u8]) -> Result<Vec<u8>, CryptoError> {
        if key.len() != 32 {
            return Err(CryptoError::Encryption("Key must be 32 bytes".to_string()));
        }
        if nonce.len() != 12 {
            return Err(CryptoError::Encryption("Nonce must be 12 bytes".to_string()));
        }

        let cipher = ChaCha20Poly1305::new(Key::from_slice(key));
        let nonce_obj = Nonce::from_slice(nonce);
        
        let mut buffer = data.to_vec();
        cipher.encrypt_in_place(nonce_obj, b"", &mut buffer)
            .map_err(|e: chacha20poly1305::Error| CryptoError::Encryption(e.to_string()))?;
        
        Ok(buffer)
    }

    /// Decrypt data using ChaCha20-Poly1305
    pub fn decrypt(encrypted_data: &[u8], key: &[u8], nonce: &[u8]) -> Result<Vec<u8>, CryptoError> {
        if key.len() != 32 {
            return Err(CryptoError::Decryption("Key must be 32 bytes".to_string()));
        }
        if nonce.len() != 12 {
            return Err(CryptoError::Decryption("Nonce must be 12 bytes".to_string()));
        }

        let cipher = ChaCha20Poly1305::new(Key::from_slice(key));
        let nonce_obj = Nonce::from_slice(nonce);
        
        let mut buffer = encrypted_data.to_vec();
        cipher.decrypt_in_place(nonce_obj, b"", &mut buffer)
            .map_err(|e: chacha20poly1305::Error| CryptoError::Decryption(e.to_string()))?;
        
        Ok(buffer)
    }

    /// Encode bytes to base64 string
    pub fn encode_base64(data: &[u8]) -> String {
        BASE64.encode(data)
    }

    /// Decode base64 string to bytes
    pub fn decode_base64(data: &str) -> Result<Vec<u8>, CryptoError> {
        BASE64.decode(data)
            .map_err(|e| CryptoError::Decryption(e.to_string()))
    }
}
