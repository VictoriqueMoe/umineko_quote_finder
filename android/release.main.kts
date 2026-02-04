#!/usr/bin/env kotlin

import java.io.File

val RESET = "\u001B[0m"
val RED = "\u001B[31m"
val GREEN = "\u001B[32m"
val YELLOW = "\u001B[33m"
val BLUE = "\u001B[34m"

fun printColoured(message: String, colour: String) {
    println("$colour$message$RESET")
}

fun validateSemanticVersion(version: String): Boolean {
    val regex = Regex("""^\d+\.\d+\.\d+$""")
    return regex.matches(version)
}

fun getCurrentVersionInfo(buildFile: File): Pair<Int, String> {
    val content = buildFile.readText()
    val versionCodeRegex = Regex("""versionCode\s*=\s*(\d+)""")
    val versionNameRegex = Regex("""versionName\s*=\s*"([^"]+)"""")

    val versionCode = versionCodeRegex.find(content)?.groupValues?.get(1)?.toInt()
        ?: throw IllegalStateException("Could not find versionCode in build.gradle.kts")
    val versionName = versionNameRegex.find(content)?.groupValues?.get(1)
        ?: throw IllegalStateException("Could not find versionName in build.gradle.kts")

    return Pair(versionCode, versionName)
}

fun updateBuildFile(buildFile: File, newVersionCode: Int, newVersionName: String) {
    var content = buildFile.readText()

    content = content.replace(
        Regex("""versionCode\s*=\s*\d+"""),
        "versionCode = $newVersionCode"
    )

    content = content.replace(
        Regex("""versionName\s*=\s*"[^"]+""""),
        """versionName = "$newVersionName""""
    )

    buildFile.writeText(content)
}

fun runCommand(command: String): Boolean {
    printColoured("\n→ Running: $command", BLUE)
    val process = ProcessBuilder(*command.split(" ").toTypedArray())
        .redirectOutput(ProcessBuilder.Redirect.INHERIT)
        .redirectError(ProcessBuilder.Redirect.INHERIT)
        .start()

    val exitCode = process.waitFor()
    return exitCode == 0
}

try {
    printColoured("╔════════════════════════════════════════╗", GREEN)
    printColoured("║   Umineko Quotes Release Build         ║", GREEN)
    printColoured("╚════════════════════════════════════════╝", GREEN)

    val buildFile = File("app/build.gradle.kts")
    if (!buildFile.exists()) {
        printColoured("Error: app/build.gradle.kts not found!", RED)
        printColoured("Make sure you're running this script from the android/ directory.", YELLOW)
        kotlin.system.exitProcess(1)
    }

    val (currentVersionCode, currentVersionName) = getCurrentVersionInfo(buildFile)
    printColoured("\nCurrent version: $currentVersionName (code: $currentVersionCode)", BLUE)

    print("\n${YELLOW}Enter new version (x.x.x format): $RESET")
    val newVersionName = readLine()?.trim() ?: ""

    if (newVersionName.isEmpty()) {
        printColoured("Error: Version cannot be empty!", RED)
        kotlin.system.exitProcess(1)
    }

    if (!validateSemanticVersion(newVersionName)) {
        printColoured("Error: Version must be in semantic versioning format (x.x.x)", RED)
        printColoured("Example: 1.0.0, 1.1.0, 2.0.1", YELLOW)
        kotlin.system.exitProcess(1)
    }

    val newVersionCode = currentVersionCode + 1

    printColoured("\n┌─────────────────────────────────────┐", GREEN)
    printColoured("│ Version Changes:                    │", GREEN)
    printColoured("│ Version Name: $currentVersionName → $newVersionName${" ".repeat(maxOf(0, 17 - currentVersionName.length - newVersionName.length))}│", GREEN)
    printColoured("│ Version Code: $currentVersionCode → $newVersionCode${" ".repeat(19)}│", GREEN)
    printColoured("└─────────────────────────────────────┘", GREEN)

    print("\n${YELLOW}Proceed with release build? (y/n): $RESET")
    val confirm = readLine()?.trim()?.lowercase()

    if (confirm != "y" && confirm != "yes") {
        printColoured("Build cancelled.", YELLOW)
        kotlin.system.exitProcess(0)
    }

    printColoured("\n[1/3] Updating build.gradle.kts...", BLUE)
    updateBuildFile(buildFile, newVersionCode, newVersionName)
    printColoured("✓ Version updated successfully", GREEN)

    printColoured("\n[2/3] Running Gradle clean...", BLUE)
    if (!runCommand("./gradlew clean")) {
        printColoured("✗ Gradle clean failed!", RED)
        kotlin.system.exitProcess(1)
    }
    printColoured("✓ Clean completed", GREEN)

    printColoured("\n[3/3] Running Gradle assembleRelease...", BLUE)
    if (!runCommand("./gradlew assembleRelease")) {
        printColoured("✗ Gradle assembleRelease failed!", RED)
        kotlin.system.exitProcess(1)
    }
    printColoured("✓ APK build completed", GREEN)

    printColoured("\n╔════════════════════════════════════════╗", GREEN)
    printColoured("║        BUILD SUCCESSFUL! ✓             ║", GREEN)
    printColoured("╚════════════════════════════════════════╝", GREEN)
    printColoured("\nAPK location: app/build/outputs/apk/release/umineko-quotes-v$newVersionName.apk", BLUE)

} catch (e: Exception) {
    printColoured("\n✗ Error: ${e.message}", RED)
    kotlin.system.exitProcess(1)
}
