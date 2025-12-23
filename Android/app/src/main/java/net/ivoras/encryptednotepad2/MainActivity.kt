package net.ivoras.encryptednotepad2

import android.app.Activity
import android.content.Intent
import android.net.Uri
import android.os.Bundle
import android.text.Editable
import android.text.TextWatcher
import android.view.View
import android.widget.EditText
import android.widget.PopupMenu
import android.widget.TextView
import android.widget.Toast
import androidx.activity.result.contract.ActivityResultContracts
import androidx.appcompat.app.AppCompatActivity
import com.google.android.material.dialog.MaterialAlertDialogBuilder
import com.google.android.material.floatingactionbutton.FloatingActionButton
import java.io.BufferedReader
import java.io.InputStreamReader
import java.io.OutputStreamWriter

class MainActivity : AppCompatActivity() {

    private lateinit var editText: EditText
    private lateinit var modifiedIndicator: TextView
    private lateinit var filenameText: TextView
    private lateinit var lineCountText: TextView
    private lateinit var fab: FloatingActionButton
    private lateinit var settingsManager: SettingsManager

    // State
    private var currentFileUri: Uri? = null
    private var currentPassword: String? = null
    private var isModified: Boolean = false
    private var savedContent: String = ""

    // File picker launchers
    private val openFileLauncher = registerForActivityResult(
        ActivityResultContracts.OpenDocument()
    ) { uri ->
        uri?.let { handleOpenFile(it) }
    }

    // Use StartActivityForResult for dynamic MIME type support
    private val saveAsLauncher = registerForActivityResult(
        ActivityResultContracts.StartActivityForResult()
    ) { result ->
        if (result.resultCode == Activity.RESULT_OK) {
            result.data?.data?.let { handleSaveAs(it) }
        }
    }

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.activity_main)

        settingsManager = SettingsManager(this)
        initViews()
        setupEditor()
        setupFab()

        // Handle intent if opened from file manager
        handleIncomingIntent(intent)
    }

    override fun onNewIntent(intent: Intent) {
        super.onNewIntent(intent)
        handleIncomingIntent(intent)
    }

    private fun initViews() {
        editText = findViewById(R.id.editText)
        modifiedIndicator = findViewById(R.id.modifiedIndicator)
        filenameText = findViewById(R.id.filenameText)
        lineCountText = findViewById(R.id.lineCountText)
        fab = findViewById(R.id.fab)
    }

    private fun setupEditor() {
        editText.addTextChangedListener(object : TextWatcher {
            override fun beforeTextChanged(s: CharSequence?, start: Int, count: Int, after: Int) {}

            override fun onTextChanged(s: CharSequence?, start: Int, before: Int, count: Int) {}

            override fun afterTextChanged(s: Editable?) {
                updateModifiedState()
                updateLineCount()
            }
        })

        // Update line count on cursor position change
        editText.setOnClickListener {
            updateLineCount()
        }

        editText.accessibilityDelegate = object : View.AccessibilityDelegate() {}

        updateLineCount()
        updateModifiedState()
    }

    private fun setupFab() {
        fab.setOnClickListener { view ->
            showFabMenu(view)
        }
    }

    private fun showFabMenu(anchor: View) {
        val popup = PopupMenu(this, anchor)
        popup.menuInflater.inflate(R.menu.fab_menu, popup.menu)

        popup.setOnMenuItemClickListener { item ->
            when (item.itemId) {
                R.id.action_new -> {
                    confirmNewFile()
                    true
                }
                R.id.action_open -> {
                    confirmOpenFile()
                    true
                }
                R.id.action_save -> {
                    saveFile()
                    true
                }
                R.id.action_save_as -> {
                    saveFileAs()
                    true
                }
                R.id.action_settings -> {
                    SettingsDialog.show(this, settingsManager)
                    true
                }
                R.id.action_about -> {
                    AboutDialog.show(this)
                    true
                }
                else -> false
            }
        }

        popup.show()
    }

    private fun confirmNewFile() {
        if (isModified) {
            MaterialAlertDialogBuilder(this)
                .setTitle(R.string.unsaved_changes_title)
                .setMessage(R.string.unsaved_changes_message)
                .setPositiveButton(R.string.discard) { _, _ ->
                    newFile()
                }
                .setNegativeButton(R.string.cancel, null)
                .setNeutralButton(R.string.save) { _, _ ->
                    saveFileThen { newFile() }
                }
                .show()
        } else {
            newFile()
        }
    }

    private fun confirmOpenFile() {
        if (isModified) {
            MaterialAlertDialogBuilder(this)
                .setTitle(R.string.unsaved_changes_title)
                .setMessage(R.string.unsaved_changes_message)
                .setPositiveButton(R.string.discard) { _, _ ->
                    openFile()
                }
                .setNegativeButton(R.string.cancel, null)
                .setNeutralButton(R.string.save) { _, _ ->
                    saveFileThen { openFile() }
                }
                .show()
        } else {
            openFile()
        }
    }

    private fun newFile() {
        editText.setText("")
        currentFileUri = null
        currentPassword = null
        savedContent = ""
        filenameText.text = getString(R.string.new_file_name)
        updateModifiedState()
        updateLineCount()
    }

    private fun openFile() {
        openFileLauncher.launch(arrayOf("*/*"))
    }

    private fun handleOpenFile(uri: Uri) {
        try {
            // Take persistent permission
            contentResolver.takePersistableUriPermission(
                uri,
                Intent.FLAG_GRANT_READ_URI_PERMISSION or Intent.FLAG_GRANT_WRITE_URI_PERMISSION
            )
        } catch (e: SecurityException) {
            // Permission not available, continue anyway
        }

        try {
            val content = readFileContent(uri)

            if (PgpEncryptionHelper.isPgpMessage(content)) {
                // Encrypted file - ask for password
                PasswordDialog.showOpenPasswordDialog(this) { password ->
                    if (password != null) {
                        decryptAndLoad(uri, content, password)
                    }
                }
            } else {
                // Plain text file - load directly
                loadContent(uri, content, null)
            }
        } catch (e: Exception) {
            Toast.makeText(
                this,
                getString(R.string.error_reading_file, e.message),
                Toast.LENGTH_LONG
            ).show()
        }
    }

    private fun decryptAndLoad(uri: Uri, encryptedContent: String, password: String) {
        try {
            val decrypted = PgpEncryptionHelper.decrypt(encryptedContent, password)
            loadContent(uri, decrypted, password)
        } catch (e: Exception) {
            Toast.makeText(
                this,
                getString(R.string.decryption_failed),
                Toast.LENGTH_LONG
            ).show()
        }
    }

    private fun loadContent(uri: Uri, content: String, password: String?) {
        editText.setText(content)
        currentFileUri = uri
        currentPassword = password
        savedContent = content
        filenameText.text = getFileName(uri)
        editText.setSelection(0)
        updateModifiedState()
        updateLineCount()
    }

    private fun readFileContent(uri: Uri): String {
        val inputStream = contentResolver.openInputStream(uri)
            ?: throw Exception("Cannot open file")

        return BufferedReader(InputStreamReader(inputStream, Charsets.UTF_8)).use { reader ->
            reader.readText()
        }
    }

    private fun saveFile() {
        saveFileThen(null)
    }

    private fun saveFileThen(onComplete: (() -> Unit)?) {
        if (currentFileUri == null) {
            // No file yet - need to use Save As
            saveFileAsThen(onComplete)
            return
        }

        if (currentPassword == null) {
            // Need a password first
            PasswordDialog.showNewPasswordDialog(this) { password ->
                if (password != null) {
                    currentPassword = password
                    performSave(currentFileUri!!, onComplete)
                }
            }
        } else {
            performSave(currentFileUri!!, onComplete)
        }
    }

    private fun saveFileAs() {
        saveFileAsThen(null)
    }

    private fun saveFileAsThen(onComplete: (() -> Unit)?) {
        // For Save As, always ask for a new password
        PasswordDialog.showNewPasswordDialog(this) { password ->
            if (password != null) {
                pendingSavePassword = password
                pendingSaveCallback = onComplete

                val suggestedName = if (currentFileUri != null) {
                    getFileName(currentFileUri!!)
                } else {
                    "document.asc"
                }

                // Create intent with dynamic MIME type based on settings
                val mimeType = if (settingsManager.usePgpEncryptedMime) {
                    "application/pgp-encrypted"
                } else {
                    "*/*"
                }

                val intent = Intent(Intent.ACTION_CREATE_DOCUMENT).apply {
                    addCategory(Intent.CATEGORY_OPENABLE)
                    type = mimeType
                    putExtra(Intent.EXTRA_TITLE, suggestedName)
                }

                saveAsLauncher.launch(intent)
            }
        }
    }

    private var pendingSavePassword: String? = null
    private var pendingSaveCallback: (() -> Unit)? = null

    private fun handleSaveAs(uri: Uri) {
        val password = pendingSavePassword
        val callback = pendingSaveCallback
        pendingSavePassword = null
        pendingSaveCallback = null

        if (password != null) {
            currentPassword = password
            currentFileUri = uri
            filenameText.text = getFileName(uri)
            performSave(uri, callback)
        }
    }

    private fun performSave(uri: Uri, onComplete: (() -> Unit)?) {
        val content = editText.text.toString()
        val password = currentPassword

        if (password == null) {
            Toast.makeText(this, R.string.no_password_set, Toast.LENGTH_SHORT).show()
            return
        }

        try {
            val encrypted = PgpEncryptionHelper.encrypt(content, password)
            writeFileContent(uri, encrypted)
            savedContent = content
            updateModifiedState()
            Toast.makeText(this, R.string.file_saved, Toast.LENGTH_SHORT).show()
            onComplete?.invoke()
        } catch (e: Exception) {
            Toast.makeText(
                this,
                getString(R.string.error_saving_file, e.message),
                Toast.LENGTH_LONG
            ).show()
        }
    }

    private fun writeFileContent(uri: Uri, content: String) {
        val outputStream = contentResolver.openOutputStream(uri, "wt")
            ?: throw Exception("Cannot open file for writing")

        OutputStreamWriter(outputStream, Charsets.UTF_8).use { writer ->
            writer.write(content)
        }
    }

    private fun getFileName(uri: Uri): String {
        var name = uri.lastPathSegment ?: "document.asc"

        // Try to get display name from content resolver
        try {
            contentResolver.query(uri, null, null, null, null)?.use { cursor ->
                if (cursor.moveToFirst()) {
                    val displayNameIndex = cursor.getColumnIndex(android.provider.OpenableColumns.DISPLAY_NAME)
                    if (displayNameIndex >= 0) {
                        name = cursor.getString(displayNameIndex) ?: name
                    }
                }
            }
        } catch (e: Exception) {
            // Use fallback name
        }

        return name
    }

    private fun updateModifiedState() {
        val currentContent = editText.text.toString()
        isModified = currentContent != savedContent
        modifiedIndicator.visibility = if (isModified) View.VISIBLE else View.GONE
    }

    private fun updateLineCount() {
        val text = editText.text.toString()
        val totalLines = if (text.isEmpty()) 1 else text.count { it == '\n' } + 1

        // Get current line
        val selectionStart = editText.selectionStart
        val currentLine = if (selectionStart >= 0) {
            text.substring(0, selectionStart.coerceAtMost(text.length)).count { it == '\n' } + 1
        } else {
            1
        }

        lineCountText.text = getString(R.string.line_count_format, currentLine, totalLines)
    }

    private fun handleIncomingIntent(intent: Intent?) {
        val uri = intent?.data
        if (uri != null && intent.action == Intent.ACTION_VIEW) {
            if (isModified) {
                MaterialAlertDialogBuilder(this)
                    .setTitle(R.string.unsaved_changes_title)
                    .setMessage(R.string.unsaved_changes_message)
                    .setPositiveButton(R.string.discard) { _, _ ->
                        handleOpenFile(uri)
                    }
                    .setNegativeButton(R.string.cancel, null)
                    .show()
            } else {
                handleOpenFile(uri)
            }
        }
    }

    @Deprecated("Use OnBackPressedDispatcher instead")
    override fun onBackPressed() {
        if (isModified) {
            MaterialAlertDialogBuilder(this)
                .setTitle(R.string.unsaved_changes_title)
                .setMessage(R.string.exit_without_saving_message)
                .setPositiveButton(R.string.discard) { _, _ ->
                    @Suppress("DEPRECATION")
                    super.onBackPressed()
                }
                .setNegativeButton(R.string.cancel, null)
                .setNeutralButton(R.string.save) { _, _ ->
                    saveFileThen {
                        @Suppress("DEPRECATION")
                        super.onBackPressed()
                    }
                }
                .show()
        } else {
            @Suppress("DEPRECATION")
            super.onBackPressed()
        }
    }
}
