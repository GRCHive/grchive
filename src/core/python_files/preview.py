#!/usr/bin/env python

import subprocess

# Returns whether or not the conversion was successful
def generatePreviewForFile(fname, outputDir):
    from .filetype import getFiletypeFromFilename, Filetype
    filetype = getFiletypeFromFilename(fname)

    if filetype == Filetype.UNKNOWN:
        return False

    ret = subprocess.call([
        'libreoffice',
        '--headless',
        '--convert-to', 'pdf',
        '--outdir', outputDir,
        fname
    ])
    return (ret == 0)
