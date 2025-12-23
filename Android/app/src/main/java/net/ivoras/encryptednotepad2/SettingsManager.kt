package net.ivoras.encryptednotepad2

import android.content.Context
import android.content.SharedPreferences

/**
 * Manages app settings using SharedPreferences.
 */
class SettingsManager(context: Context) {

    private val prefs: SharedPreferences = context.getSharedPreferences(
        PREFS_NAME,
        Context.MODE_PRIVATE
    )

    /**
     * Whether to use application/pgp-encrypted MIME type when saving files.
     * If false, uses wildcard for broader compatibility.
     */
    var usePgpEncryptedMime: Boolean
        get() = prefs.getBoolean(KEY_USE_PGP_MIME, DEFAULT_USE_PGP_MIME)
        set(value) = prefs.edit().putBoolean(KEY_USE_PGP_MIME, value).apply()

    companion object {
        private const val PREFS_NAME = "encrypted_notepad_settings"
        private const val KEY_USE_PGP_MIME = "use_pgp_encrypted_mime"
        private const val DEFAULT_USE_PGP_MIME = false
    }
}
