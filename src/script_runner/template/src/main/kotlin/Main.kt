package com.grchive.script.runner

import grchive.core.runner.invokeWithMetadata
import grchive.core.runner.loadMetadataFromResource

import com.github.ajalt.clikt.core.CliktCommand
import com.github.ajalt.clikt.parameters.arguments.argument
import com.github.ajalt.clikt.parameters.types.long

class Cli : CliktCommand(help = "Wrapper program to run client scripts.") {
    val cls : String by argument(help = "Client script Kotlin class full path.")
    val fn : String by argument(help = "Client script function name.")
    val meta : String by argument(help = "Metadata resource name.")
    val runId : Long by argument(help = "Run ID.").long()

    override fun run() {
        invokeWithMetadata(runId, cls, fn, loadMetadataFromResource(meta))
    }
}

fun main(args : Array<String>)  = Cli().main(args)
