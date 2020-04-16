package com.grchive.script.runner

import grchive.core.utility.hexEncode
import com.grchive.web.client.s2

fun main() {
    s2()

    val test = "Testing".toByteArray()
    println(hexEncode(test))
}
