import 'dart:ffi' as ffi;
import 'dart:io';
import 'dart:typed_data';
import 'package:ffi/ffi.dart';
import 'package:path/path.dart';

// Dart type definition for calling the C foreign function
typedef FetchWebsite = ffi.Pointer<FetchWebsiteReturn> Function(
    ffi.Pointer<Utf8>, ffi.Pointer<Utf8>);

class FetchWebsiteReturn extends ffi.Struct {
  external ffi.Pointer<ffi.Void> content;

  @ffi.Int32()
  external int size;

  external ffi.Pointer<Utf8> error;
}

class Binding {
  static final String _libraryName = 'libthyra';
  static final Binding _singleton = Binding._internal();

  late ffi.DynamicLibrary _library;

  factory Binding() {
    return _singleton;
  }

  Binding._internal() {
    _library = openLib();
  }

  ffi.DynamicLibrary openLib() {
    if (Platform.isMacOS || Platform.isWindows) {
      throw ("Lib doesn't exist for this plateform.");
    }
    if (Platform.isIOS) {
      return ffi.DynamicLibrary.process();
    }
    if (Platform.isLinux) {
      return ffi.DynamicLibrary.open(
          absolute("../../build/libraries/linux_amd64/libthyra.so"));
    }
    if (Platform.isAndroid) {
      return ffi.DynamicLibrary.open('$_libraryName.so');
    }
    try {
      return ffi.DynamicLibrary.open("$_libraryName.so");
    } catch (e) {
      print("fallback to open DynamicLibrary on older devices");
      //fallback for devices that cannot load dynamic libraries by name: load the library with an absolute path
      //read the app id
      var appid = File("/proc/self/cmdline").readAsStringSync();
      // the file /proc/self/cmdline returns a string with many trailing \0 characters, which makes the string pretty useless for dart, many
      // operations will not work correctly. remove these trailing zero bytes.
      appid = String.fromCharCodes(
          appid.codeUnits.where((element) => element != 0));
      final loadPath = "/data/data/$appid/lib/$_libraryName.so";
      return ffi.DynamicLibrary.open(loadPath);
    }
  }

  Uint8List fetchWebsite(String address, String file) {
    final FetchWebsite nativeFunc = _library
        .lookup<ffi.NativeFunction<FetchWebsite>>('fetchWebsite')
        .asFunction();

    var resp = nativeFunc(address.toNativeUtf8(), file.toNativeUtf8());

    if (resp.ref.error.address != ffi.nullptr.address) {
      var message = resp.ref.error.toDartString();
      if (!Platform.isWindows) {
        malloc.free(resp);
      }
      throw Exception(message);
    }

    //print('Content size is: ${resp.ref.size}');
    final content =
        resp.ref.content.cast<ffi.Uint8>().asTypedList(resp.ref.size);

    if (!Platform.isWindows) {
      malloc.free(resp);
    }

    return content;
  }
}
