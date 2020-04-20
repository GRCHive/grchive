load("//build:variables.bzl", "env")
def core_lib_version_str():
    return "{0}.{1}".format(
        env["KOTLIN_CORE_LIB_MAJOR_VERSION"],
        env["KOTLIN_CORE_LIB_MINOR_VERSION"],
    )

def core_lib_maven_dep():
    return "{0}:{1}:{2}".format(
        env["KOTLIN_CORE_LIB_GROUP_ID"],
        env["KOTLIN_CORE_LIB_ARTIFACT_ID"],
        core_lib_version_str(),
    )

def core_lib_maven_dep_to_dir():
    return env["KOTLIN_CORE_LIB_GROUP_ID"].replace(".", "/") + "/" + \
        env["KOTLIN_CORE_LIB_ARTIFACT_ID"].replace(".", "/") + "/" + \
        core_lib_version_str()

def core_lib_fname(ext = 'jar'):
    return '{0}-{1}.{2}'.format(env["KOTLIN_CORE_LIB_ARTIFACT_ID"], core_lib_version_str(), ext)
