package net.ivoras.encryptednotepad2

import org.bouncycastle.openpgp.PGPException
import org.pgpainless.PGPainless
import org.pgpainless.algorithm.SymmetricKeyAlgorithm
import org.pgpainless.decryption_verification.ConsumerOptions
import org.pgpainless.encryption_signing.EncryptionOptions
import org.pgpainless.encryption_signing.ProducerOptions
import org.pgpainless.util.Passphrase
import java.io.ByteArrayInputStream
import java.io.ByteArrayOutputStream
import java.io.IOException

/**
 * Helper class for OpenPGP symmetric encryption/decryption using PGPainless.
 * Produces ASCII-armored output compatible with other OpenPGP implementations.
 */
object PgpEncryptionHelper {

    /**
     * Encrypts plaintext using symmetric (password-based) encryption.
     *
     * @param plaintext The text to encrypt
     * @param password The password for encryption
     * @return ASCII-armored encrypted message
     * @throws PGPException if encryption fails
     * @throws IOException if I/O error occurs
     */
    @Throws(PGPException::class, IOException::class)
    fun encrypt(plaintext: String, password: String): String {
        val passphrase = Passphrase.fromPassword(password)

        val plaintextBytes = plaintext.toByteArray(Charsets.UTF_8)
        val inputStream = ByteArrayInputStream(plaintextBytes)
        val outputStream = ByteArrayOutputStream()

        val encryptionOptions = EncryptionOptions()
            .addPassphrase(passphrase)
            .overrideEncryptionAlgorithm(SymmetricKeyAlgorithm.AES_256)

        val producerOptions = ProducerOptions
            .encrypt(encryptionOptions)
            .setAsciiArmor(true)

        val encryptionStream = PGPainless.encryptAndOrSign()
            .onOutputStream(outputStream)
            .withOptions(producerOptions)

        inputStream.copyTo(encryptionStream)
        encryptionStream.close()

        return outputStream.toString(Charsets.UTF_8.name())
    }

    /**
     * Decrypts an ASCII-armored message using the provided password.
     *
     * @param armoredCiphertext The ASCII-armored encrypted message
     * @param password The password for decryption
     * @return The decrypted plaintext
     * @throws PGPException if decryption fails (wrong password or corrupted data)
     * @throws IOException if I/O error occurs
     */
    @Throws(PGPException::class, IOException::class)
    fun decrypt(armoredCiphertext: String, password: String): String {
        val passphrase = Passphrase.fromPassword(password)

        val ciphertextBytes = armoredCiphertext.toByteArray(Charsets.UTF_8)
        val inputStream = ByteArrayInputStream(ciphertextBytes)
        val outputStream = ByteArrayOutputStream()

        val consumerOptions = ConsumerOptions()
            .addDecryptionPassphrase(passphrase)

        val decryptionStream = PGPainless.decryptAndOrVerify()
            .onInputStream(inputStream)
            .withOptions(consumerOptions)

        decryptionStream.copyTo(outputStream)
        decryptionStream.close()

        return outputStream.toString(Charsets.UTF_8.name())
    }

    /**
     * Checks if the given text appears to be an ASCII-armored PGP message.
     *
     * @param text The text to check
     * @return true if it looks like a PGP message
     */
    fun isPgpMessage(text: String): Boolean {
        val trimmed = text.trim()
        return trimmed.startsWith("-----BEGIN PGP MESSAGE-----") &&
               trimmed.contains("-----END PGP MESSAGE-----")
    }
}
