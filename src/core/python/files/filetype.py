from enum import Enum, auto

import magic

class Filetype(Enum):
    PDF = auto()
    IMAGE = auto()
    TEXT = auto()
    WORD = auto()
    EXCEL = auto()
    UNKNOWN = auto()

def getFiletypeFromFilename(fname):
    magicType = magic.detect_from_filename(fname)

    if magicType.encoding == 'binary':
        if magicType.mime_type == 'application/pdf':
            return Filetype.PDF
        elif magicType.mime_type.startswith('image/'):
            return Filetype.IMAGE
        elif magicType.mime_type.startswith('application/vnd.openxmlformats-officedocument'):
            if 'spreadsheetml' in magicType.mime_type:
                return Filetype.EXCEL
            elif 'wordprocessingml' in magicType.mime_type:
                return Filetype.WORD
        return Filetype.UNKNOWN

    return Filetype.TEXT
