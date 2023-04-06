import base64

__version__ = '1.0.0'


class ArgLibrary(object):
    ROBOT_LIBRARY_VERSION = __version__
    ROBOT_LIBRARY_SCOPE = 'GLOBAL'

    def string_to_arg(self, string):
        bytes = len(string).to_bytes(4, 'little') + string.encode('utf-8')
        return base64.b64encode(bytes).decode('utf-8')

