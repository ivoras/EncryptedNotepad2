package net.ivoras.encryptednotepad2

import android.content.Context
import android.view.LayoutInflater
import com.google.android.material.dialog.MaterialAlertDialogBuilder
import com.google.android.material.textfield.TextInputEditText
import com.google.android.material.textfield.TextInputLayout

/**
 * Dialog for password entry with confirmation.
 * Used when saving a file for the first time or using Save As.
 */
object PasswordDialog {

    private const val MIN_PASSWORD_LENGTH = 1

    /**
     * Shows a dialog asking for a new password with confirmation.
     *
     * @param context The context
     * @param onPasswordEntered Callback with the password if valid, or null if cancelled
     */
    fun showNewPasswordDialog(
        context: Context,
        onPasswordEntered: (password: String?) -> Unit
    ) {
        val view = LayoutInflater.from(context).inflate(R.layout.dialog_password, null)
        val passwordLayout = view.findViewById<TextInputLayout>(R.id.passwordLayout)
        val passwordInput = view.findViewById<TextInputEditText>(R.id.passwordInput)
        val confirmPasswordLayout = view.findViewById<TextInputLayout>(R.id.confirmPasswordLayout)
        val confirmPasswordInput = view.findViewById<TextInputEditText>(R.id.confirmPasswordInput)

        val dialog = MaterialAlertDialogBuilder(context)
            .setTitle(R.string.password_dialog_title)
            .setView(view)
            .setPositiveButton(R.string.ok, null) // Set to null initially to handle validation
            .setNegativeButton(R.string.cancel) { dialog, _ ->
                dialog.dismiss()
                onPasswordEntered(null)
            }
            .setOnCancelListener {
                onPasswordEntered(null)
            }
            .create()

        dialog.setOnShowListener {
            val positiveButton = dialog.getButton(android.app.AlertDialog.BUTTON_POSITIVE)
            positiveButton.setOnClickListener {
                val password = passwordInput.text?.toString() ?: ""
                val confirmPassword = confirmPasswordInput.text?.toString() ?: ""

                // Clear previous errors
                passwordLayout.error = null
                confirmPasswordLayout.error = null

                // Validate password length
                if (password.length < MIN_PASSWORD_LENGTH) {
                    passwordLayout.error = context.getString(R.string.password_too_short)
                    return@setOnClickListener
                }

                // Validate passwords match
                if (password != confirmPassword) {
                    confirmPasswordLayout.error = context.getString(R.string.passwords_dont_match)
                    return@setOnClickListener
                }

                dialog.dismiss()
                onPasswordEntered(password)
            }
        }

        dialog.show()
        passwordInput.requestFocus()
    }

    /**
     * Shows a dialog asking for a password to open/decrypt a file.
     *
     * @param context The context
     * @param onPasswordEntered Callback with the password if entered, or null if cancelled
     */
    fun showOpenPasswordDialog(
        context: Context,
        onPasswordEntered: (password: String?) -> Unit
    ) {
        val view = LayoutInflater.from(context).inflate(R.layout.dialog_open_password, null)
        val passwordLayout = view.findViewById<TextInputLayout>(R.id.passwordLayout)
        val passwordInput = view.findViewById<TextInputEditText>(R.id.passwordInput)

        val dialog = MaterialAlertDialogBuilder(context)
            .setTitle(R.string.enter_password_dialog_title)
            .setView(view)
            .setPositiveButton(R.string.ok, null)
            .setNegativeButton(R.string.cancel) { dialog, _ ->
                dialog.dismiss()
                onPasswordEntered(null)
            }
            .setOnCancelListener {
                onPasswordEntered(null)
            }
            .create()

        dialog.setOnShowListener {
            val positiveButton = dialog.getButton(android.app.AlertDialog.BUTTON_POSITIVE)
            positiveButton.setOnClickListener {
                val password = passwordInput.text?.toString() ?: ""

                passwordLayout.error = null

                if (password.isEmpty()) {
                    passwordLayout.error = context.getString(R.string.password_required)
                    return@setOnClickListener
                }

                dialog.dismiss()
                onPasswordEntered(password)
            }
        }

        dialog.show()
        passwordInput.requestFocus()
    }
}
