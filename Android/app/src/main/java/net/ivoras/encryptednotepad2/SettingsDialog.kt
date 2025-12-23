package net.ivoras.encryptednotepad2

import android.content.Context
import android.view.LayoutInflater
import android.widget.CheckBox
import com.google.android.material.dialog.MaterialAlertDialogBuilder

/**
 * Dialog for app settings.
 */
object SettingsDialog {

    /**
     * Shows the Settings dialog.
     *
     * @param context The context
     * @param settingsManager The settings manager instance
     */
    fun show(context: Context, settingsManager: SettingsManager) {
        val view = LayoutInflater.from(context).inflate(R.layout.dialog_settings, null)
        val pgpMimeCheckbox = view.findViewById<CheckBox>(R.id.pgpMimeCheckbox)

        // Set current value
        pgpMimeCheckbox.isChecked = settingsManager.usePgpEncryptedMime

        // Save immediately on change
        pgpMimeCheckbox.setOnCheckedChangeListener { _, isChecked ->
            settingsManager.usePgpEncryptedMime = isChecked
        }

        MaterialAlertDialogBuilder(context)
            .setTitle(R.string.settings_title)
            .setView(view)
            .setPositiveButton(R.string.close, null)
            .show()
    }
}
