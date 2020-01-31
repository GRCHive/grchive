load("//build:variables.bzl", "env")

def generateBashEnvironment():
    bashVariables = []
    for k, v in env.items():
        if type(v) == "string":
            formatV = "\"" + v + "\""
        else:
            formatV = v

        bashVariables.append("export {name}={value}".format(name=k, value=formatV))

    return "\n".join(bashVariables)
