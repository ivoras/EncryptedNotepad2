package net.ivoras.encryptednotepad2

import android.content.Context
import android.content.pm.PackageManager
import android.os.Build
import com.google.android.material.dialog.MaterialAlertDialogBuilder

/**
 * Dialog showing information about the app.
 */
object AboutDialog {

    /**
     * Shows the About dialog.
     *
     * @param context The context
     */
    fun show(context: Context) {
        val versionName = getVersionName(context)

        val message = context.getString(R.string.about_message, versionName)

        MaterialAlertDialogBuilder(context)
            .setTitle(R.string.app_name)
            .setMessage(message)
            .setPositiveButton(R.string.ok) { dialog, _ ->
                dialog.dismiss()
            }
            .setIcon(R.mipmap.ic_launcher)
            .show()
    }

    private fun getVersionName(context: Context): String {
        return try {
            val packageInfo = if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.TIRAMISU) {
                context.packageManager.getPackageInfo(
                    context.packageName,
                    PackageManager.PackageInfoFlags.of(0)
                )
            } else {
                @Suppress("DEPRECATION")
                context.packageManager.getPackageInfo(context.packageName, 0)
            }
            packageInfo.versionName ?: "1.0"
        } catch (e: PackageManager.NameNotFoundException) {
            "1.0"
        }
    }
}
