import 'package:thyra_bridge_dart/thyra_bridge_dart.dart';
import 'package:test/test.dart';

void main() {
  group('A group of tests', () {
    final binding = Binding();

    setUp(() {
      // Additional setup goes here.
    });

    test('First Test', () {
      var webContent = String.fromCharCodes(binding.fetchWebsite(
          "A1aMywGBgBywiL6WcbKR4ugxoBtdP9P3waBVi5e713uvj7F1DJL", "index.html"));

      expect(webContent.length, 19065);
      expect(webContent, startsWith("<!DOCTYPE html>"));
    });
  });
}
